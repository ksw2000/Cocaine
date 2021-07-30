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

	/*
		Lexical analysis
		Syntax analysis
		Semantic analysis
		Intermediate code generation
		Optimization
		Code generation
	*/
	program := Analysis(bufio.NewReader(file))
	interCode := IntermidiateCodeGenerator(program)
	nasm := NASMGenerator(interCode)
	fmt.Println(nasm)
}
