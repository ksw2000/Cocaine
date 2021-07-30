package main

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Variable struct {
	token      int // map to const I8, I16, ...
	value      interface{}
	size       uint // unit: byte
	indirect   uint
	initialize bool
	readOnly   bool
}

type Program struct {
	importList       []string // TODO: use []Program in the future
	cLibDependencies map[string]interface{}
	main             *Scope
	functions        map[string]*Scope
}

type Scope struct {
	varibaleMap        map[string]*Variable
	annonymousVariable []*Variable
	subScope           []*Scope
}

var buildinLabelSet = map[string]int{
	"import": ParsingImport,
	"main":   ParsingMain,
	"end":    ParsingEnd,
}

var mnemonics = map[string]int{
	"mov": MneMov,
}

var instruction = map[string]int{
	"if":   InsIf,
	"else": InsElse,
	"for":  InsFor,
}

func NewScope() *Scope {
	return &Scope{
		varibaleMap:        map[string]*Variable{},
		annonymousVariable: []*Variable{},
		subScope:           []*Scope{},
	}
}
func NewProgram() *Program {
	return &Program{
		importList:       []string{},
		cLibDependencies: map[string]interface{}{},
		main:             NewScope(),
	}
}

func Analysis(r *bufio.Reader) *Program {
	buf, isPrefix, err := r.ReadLine()
	var lineNumber int = 0
	lineBuf := []byte{}
	state := ParsingInit
	program := NewProgram()
	for err == nil {
		if !isPrefix {
			lineNumber++
			line := string(lineBuf) + string(buf)
			lineBuf = []byte{}
			color.Blue("%3d | %s", lineNumber, line)
			state = program.processSingleLine(line, state)
			if state == ParsingEnd {
				break
			}
		} else {
			lineBuf = append(lineBuf, buf...)
		}
		buf, isPrefix, err = r.ReadLine()
	}
	color.Green("parsing done")
	fmt.Printf("%#v", program)
	color.Green("prepare to generate assembly code")
	return program
}

func (p *Program) processSingleLine(line string, previousState int) (nextState int) {
	if previousState == ParsingInit || previousState == ParsingLabel {
		if line != "" && line[0] != ' ' {
			return parseLabel(line, previousState)
		}
	} else if previousState == ParsingImport {
		if line != "" {
			if line[0] != ' ' {
				return parseLabel(line, previousState)
			}
			line = strings.TrimSpace(line)
			p.importList = append(p.importList, line)
			fmt.Println("import package:", line)
		}
	} else if previousState == ParsingMain {
		if line != "" {
			if line[0] != ' ' {
				return parseLabel(line, previousState)
			}
			return p.ParseCodeInFunction(line, previousState)
		}
	} else if previousState == ParsingEnd {
		fmt.Println("END!")
	}
	return previousState
}

func parseLabel(line string, previousState int) (nextState int) {
	if line != "" {
		nextState, ok := buildinLabelSet[line]
		if ok {
			fmt.Println("built-in label", line)
		} else {
			fmt.Println("custom label", line)
		}
		return nextState
	}
	return previousState
}

func (p *Program) ParseCodeInFunction(line string, previousState int) (nextState int) {
	line = strings.TrimSpace(line)
	if line == "" {
		return previousState
	}

	line = strings.TrimSpace(line)
	words := strings.Split(line, " ")
	if len(words) > 0 {
		if _, ok := mnemonics[words[0]]; ok {
			fmt.Printf("mnemonics: %s\n", words[0])
		} else if ins, ok := instruction[words[0]]; ok {
			fmt.Printf("instructions: %s %d\n", words[0], ins)
		} else {
			// variable declration or call funcion
			pointPart := strings.Split(words[0], ".")
			if len(pointPart) > 1 {
				if pointPart[0] == "c" {
					p.cLibDependencies[pointPart[1]] = struct{}{}
					fmt.Printf("call c library function: %s\n", pointPart[1])
				} else {
					fmt.Printf("call function in another package: %s\n", pointPart[1])
				}
			} else if f, ok := p.functions[words[0]]; ok {
				fmt.Printf("call function in the same package: %s %#v\n", words[0], f)
			} else if len(words) >= 3 {
				// name type val

				// TODO:
				// Should choose the specify scope
				p.main.ParseVariable(words[0], words[1], strings.Join(words[2:], " "))
				fmt.Printf("declare variable name: %s type: %s val: %s\n", words[0], words[1], words[2:])
			} else {
				fmt.Println("error")
			}
		}
	}
	/*
		1. parse mnemonics
		2. parse scope
			1. if
			2. else
			3. for
		3. parse normal function
		4. parse c library function
		5. parse variable declration
	*/
	return previousState
}

type VariableType struct {
	unit  uint
	token int
}

var buildInType = map[string]VariableType{
	"i64": {
		unit:  8,
		token: I64,
	},
	"rune": {
		unit:  32,
		token: Rune,
	},
}

func (s *Scope) ParseVariable(name string, varType string, value string) {
	if _, ok := s.varibaleMap[name]; !ok {
		s.varibaleMap[name] = &Variable{
			token:    buildInType[varType].token,    // TODO: consider custom type
			size:     buildInType[varType].unit * 1, // TODO: calcuate size
			indirect: 0,                             // TODO: the amount of *
			readOnly: false,                         // TODO: feature in the futre
		}
		// fill value
		if value != "?" {
			s.varibaleMap[name].initialize = true
			s.varibaleMap[name].value = value
		}
	} else {
		color.Red("ERROR! Declare again!")
	}
}

// only support naive string now
func (s *Scope) AnnoymousStringVariable(value string) {
	s.annonymousVariable = append(s.annonymousVariable, &Variable{
		token: buildInType["rune"].token,
		value: value,
		//size:       1,
		indirect:   1,
		initialize: true,
		readOnly:   true,
	})
}
