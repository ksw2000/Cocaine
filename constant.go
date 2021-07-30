package main

const (
	ParsingInit = iota
	ParsingImport
	ParsingMain
	ParsingLabel
	ParsingEnd
	ParsingErr
)

const (
	MneMov = iota
)

const (
	InsIf   = iota
	InsElse = iota
	InsFor  = iota
)

const (
	I8 = iota
	I16
	I32
	I64
	U8
	U16
	U32
	U64
	Bool
	Rune
	F32
	F64
	Any
)
