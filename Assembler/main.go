package main

import (
	. "Assembler/common"
	. "Assembler/frontend"
)

func main() {
	var lex LexerInterface
	lex.Init()
	var parser *Parser
	parser.Init(lex)
}
