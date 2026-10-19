package frontend

import (
	"fmt"
	"os"
	"slices"
	"strconv"

	. "Assembler/common"
)

type Parser struct {
	lexer LexerInterface
	token Token
	out   *os.File
}

func (parser *Parser) oneOf(ttcheck TokenType, tts ...TokenType) bool {
	return slices.Contains(tts, ttcheck)
}

func (parser *Parser) matchError(msg string) {
	fmt.Println("Parsing error:", msg)
	panic("Parser.match: invalid token")
}

func (parser *Parser) match(tt TokenType) {
	if parser.token.Type != tt {
		msg := fmt.Sprintf("found '%s' but expected '%s' at position %s.", parser.token.Type, tt, parser.token.Pos)
		parser.matchError(msg)
	}
	parser.token = parser.lexer.NextToken()
}

func (parser *Parser) matchIf(tt TokenType) bool {
	if parser.token.Type == tt {
		parser.match(tt)
		return true
	}
	return false
}

func (parser *Parser) Init(lex LexerInterface) {
	parser.lexer = lex
	parser.token = parser.lexer.NextToken()
}

func (parser *Parser) ParserSrc(src []byte, handler ErrorHandler, oPath string) {
	var lex Lexer
	lex.Init(src, handler)
	parser.Init(&lex)
	file, err := os.Create(oPath)
	if err != nil {
		panic(err)
	}
	parser.out = file

	parser.Program()
}

var (
	Labels map[string]int
	Const  map[string]int
)

func (parser *Parser) Program() {
	for parser.token.Type == TtKwConst {
		parser.Const()
	}
}

func (parser *Parser) Const() {
	ident := parser.token.Lexeme
	parser.match(TtIdent)
	integer := parser.token.Lexeme
	parser.match(TtInt)
	i, err := strconv.Atoi(integer)
	if err != nil {
		panic(err)
	}

	if i > 63 {
		msg := fmt.Sprintf("The integer %d at position %s is too large (this instruction set only supports immediation below 64)", i, parser.token.Pos)
		parser.matchError(msg)
	}
	parser.match(TtSemi)
	Const[ident] = i
}
