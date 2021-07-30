package main

import "fmt"

type Asm struct {
	extern []string
	data   []*SectionData
	bss    []*SectionBss
	text   []*SectionText
}

func NewAsm() *Asm {
	return &Asm{
		extern: []string{},
		data:   []*SectionData{},
		bss:    []*SectionBss{},
		text:   []*SectionText{},
	}
}

type SectionData struct {
	name string
	d    string
	val  interface{}
	// ex: msg db "hello"
}

type SectionBss struct {
	name string
	res  string
	size uint
}

const (
	Register = iota
	Memory
	Immediate
)

type SectionText struct {
	ins      string
	arg1     string
	arg1Type int
	arg2     string
	arg2Type int
}

type Function struct {
	stack uint // default 32
}

func IntermidiateCodeGenerator(p *Program) *Asm {
	a := NewAsm()

	for cld := range p.cLibDependencies {
		a.extern = append(a.extern, cld)
	}

	for name, v := range p.main.varibaleMap {
		if v.initialize {
			a.data = append(a.data, &SectionData{
				name: fmt.Sprintf("main.%s", name),
				d:    "db",
				val:  v.value,
			})
		} else {
			a.bss = append(a.bss, &SectionBss{
				name: fmt.Sprintf("main.%s", name),
				res:  "resb",
				size: v.size,
			})
		}
	}

	/*for _, v := range p.main.annonymousVariable {

	}*/

	return a
}

func NASMGenerator(a *Asm) string {
	/*
	   global main
	   extern printf

	   section .data
	       var1 dq  2021
	       var2 db  "The primary Cocaine Compiler is created in %d", 0

	   section .text
	   main:
	*/
	ret := ""
	ret += "global main\n"
	if len(a.extern) > 0 {
		ret += "extern"
		for _, e := range a.extern {
			ret += fmt.Sprintf(" %s", e)
		}
		ret += "\n"
	}
	ret += "\n"
	if len(a.bss) > 0 {
		ret += "section .bss\n"
	}
	if len(a.data) > 0 {
		ret += "section .data\n"
		for _, data := range a.data {
			ret += fmt.Sprintf("\t%s %s %s\n", data.name, data.d, data.val)
		}
	}
	ret += "section .text\n"
	ret += "main:\n"

	return ret
}
