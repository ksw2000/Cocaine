package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// parse flag
	output := flag.String("o", "", "`output` name")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "coco src.coco")
		fmt.Fprintln(os.Stderr, "coco -o src.exe src.coco")
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
		return
	}

	// generate output filename
	input := flag.Arg(0)
	if *output == "" {
		parts := strings.Split(input, ".")
		extName := parts[len(parts)-1]
		if extName == "coco" {
			*output = strings.Join(parts[0:len(parts)-1], ".") + ".exe"
		} else {
			*output = input + ".exe"
		}
	}

	// read file
	file, err := os.Open(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s does not exist", input)
		return
	}

	// TODO:
	// 0. preprocess (delete comment, process \r \n and spaces)
	// 1. compile
	// 2. generate assembly code
	// 3. assemble by NASM

	compile(bufio.NewReader(file))
}

const (
	parsingInit = iota
	parsingImport
	parsingMain
	parsingLabel
	parsingEnd
)

func compile(r *bufio.Reader) {
	buf, isPrefix, err := r.ReadLine()
	var lineNumber int = 0
	lineBuf := []byte{}
	state := parsingInit
	for err == nil {
		if !isPrefix {
			lineNumber++
			line := string(lineBuf) + string(buf)
			lineBuf = []byte{}
			fmt.Printf("%3d | %s\n", lineNumber, line)
			state = process(line, state)
			if state == parsingEnd {
				break
			}
		} else {
			lineBuf = append(lineBuf, buf...)
		}
		buf, isPrefix, err = r.ReadLine()
	}
}

var buildinLabelSet = map[string]int{
	"import": parsingImport,
	"main":   parsingMain,
	"end":    parsingEnd,
}

var importList []string

func process(line string, previousState int) (nextState int) {
	if previousState == parsingInit || previousState == parsingLabel {
		if line != "" && line[0] != ' ' {
			return parseLabel(line, previousState)
		}
	} else if previousState == parsingImport {
		if line != "" {
			if line[0] != ' ' {
				return parseLabel(line, previousState)
			}
			line = strings.TrimSpace(line)
			importList = append(importList, line)
			fmt.Println("import package:", line)
		}
	} else if previousState == parsingMain {
		if line != "" {
			if line[0] != ' ' {
				return parseLabel(line, previousState)
			}
			return parseCode(line, previousState)
		}
	} else if previousState == parsingEnd {
		fmt.Println("END!")
	}
	return previousState
}

func parseLabel(line string, previousState int) (nextState int) {
	if line != "" {
		nextState, ok := buildinLabelSet[line]
		if ok {
			fmt.Println("built in label!", line)
		} else {
			fmt.Println("custom label", line)
		}
		return nextState
	}
	return previousState
}

func parseCode(line string, previousState int) (nextState int) {
	line = strings.TrimSpace(line)
	if line == "" {
		return previousState
	}

	line = strings.TrimSpace(line)
	fmt.Println("code:", line)
	/*
		1. parse variable
		2. parse function
		3. parse some instruction
			1. if
			2. else
			3. for
			4. foo = bar
	*/
	return previousState
}
