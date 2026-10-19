package frontend

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	. "Assembler/common"
)

type Instr struct {
	Op  TokenType // TtKwAdd, TtKwMov, TtKwJz, ...
	Dst TokenType // TtKwA..TtKwF or TtKwOut (mov only)
	Src TokenType // TtKwA..TtKwF (mov only)
	Imm int
	Ref string // non-empty for "imm <ident>", resolved after parsing
	Pos Position
}

type Parser struct {
	lexer  LexerInterface
	token  Token
	out    *os.File
	instrs []Instr
	labels map[string]int
	consts map[string]int
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

func isReg(t TokenType) bool { return t >= TtKwA && t <= TtKwF }

func (p *Parser) isInstrStart() bool {
	t := p.token.Type
	return t >= TtKwNand && t <= TtKwJg && t != TtKwConst
}

func (p *Parser) Label() {
	name := strings.TrimSuffix(p.token.Lexeme, ":") // drop if the lexer already strips it
	pos := p.token.Pos
	p.match(TtLabel)

	if _, dup := p.labels[name]; dup {
		p.matchError(fmt.Sprintf("duplicate label '%s' at %s", name, pos))
	}
	p.labels[name] = len(p.instrs)

	for p.isInstrStart() {
		p.Instr()
	}
}

func (p *Parser) reg() {
	if !isReg(p.token.Type) {
		p.matchError(fmt.Sprintf("expected register at %s, found '%s'", p.token.Pos, p.token.Type))
	}
	p.match(p.token.Type)
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

func (parser *Parser) Instr() {
	in := Instr{Op: parser.token.Type, Pos: parser.token.Pos}
	parser.match(in.Op)

	switch in.Op {
	case TtKwMov:
		in.Dst = parser.token.Type
		if !parser.matchIf(TtKwOut) {
			parser.reg()
		}
		parser.match(TtComma)
		in.Src = parser.token.Type
		parser.reg()
	case TtKwImm:

		parser.match(TtKwImm)

		if parser.token.Type != TtInt {
			in.Imm = Const[parser.token.Lexeme]
			parser.match(TtIdent)
		} else {
			var err error
			in.Imm, err = strconv.Atoi(parser.token.Lexeme)
			if err != nil {
				panic(err)
			}
			parser.match(TtInt)
		}
	}
	parser.match(TtSemi)
	parser.instrs = append(parser.instrs, in)
}
