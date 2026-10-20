package frontend

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	. "Assembler/common"
)

const MaxImm = 63 // 6-bit immediate

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
	instrs []Instr
	labels map[string]int
	consts map[string]int
}

// Instrs exposes the parsed instruction list to other packages.
func (parser *Parser) Instrs() []Instr { return parser.instrs }

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
	parser.labels = make(map[string]int)
	parser.consts = make(map[string]int)
	parser.instrs = nil
	parser.token = parser.lexer.NextToken()
}

func (parser *Parser) ParserSrc(src []byte, handler ErrorHandler) {
	var lex Lexer
	lex.Init(src, handler)
	parser.Init(&lex)
	parser.Program()
}

func (parser *Parser) Program() {
	for parser.token.Type == TtKwConst {
		parser.Const()
	}

	// The first label must be main (PC resets to address 0).
	if strings.TrimSuffix(parser.token.Lexeme, ":") != "main" {
		parser.matchError(fmt.Sprintf("program must start with 'main:' at %s", parser.token.Pos))
	}
	parser.Label()
	for parser.token.Type != TtEOI {
		parser.Label()
	}
	parser.resolve()
}

func isReg(t TokenType) bool { return t >= TtKwA && t <= TtKwF }

// Relies on the order of the keyword constants in token.go:
// all mnemonics lie between TtKwNand and TtKwJg, except const.
func (parser *Parser) isInstrStart() bool {
	t := parser.token.Type
	return t >= TtKwNand && t <= TtKwJg && t != TtKwConst
}

func (parser *Parser) Label() {
	name := strings.TrimSuffix(parser.token.Lexeme, ":") // harmless if the lexer already strips it
	pos := parser.token.Pos
	parser.match(TtLabel)

	if _, dup := parser.labels[name]; dup {
		parser.matchError(fmt.Sprintf("duplicate label '%s' at %s", name, pos))
	}
	if _, dup := parser.consts[name]; dup {
		parser.matchError(fmt.Sprintf("label '%s' at %s clashes with a constant", name, pos))
	}
	parser.labels[name] = len(parser.instrs) // address of the next instruction

	for parser.isInstrStart() {
		parser.Instr()
	}
}

// reg consumes a register token and returns its type.
func (parser *Parser) reg() TokenType {
	t := parser.token.Type
	if !isReg(t) {
		parser.matchError(fmt.Sprintf("expected register at %s, found '%s'", parser.token.Pos, t))
	}
	parser.match(t)
	return t
}

func (parser *Parser) Const() {
	parser.match(TtKwConst)
	pos := parser.token.Pos
	ident := parser.token.Lexeme
	parser.match(TtIdent)
	integer := parser.token.Lexeme
	parser.match(TtInt)
	i, err := strconv.Atoi(integer)
	if err != nil {
		parser.matchError(fmt.Sprintf("invalid integer '%s' at %s", integer, pos))
	}
	if i < 0 || i > MaxImm {
		parser.matchError(fmt.Sprintf("the integer %d at %s is out of range (immediates are 0..%d)", i, pos, MaxImm))
	}
	if _, dup := parser.consts[ident]; dup {
		parser.matchError(fmt.Sprintf("duplicate constant '%s' at %s", ident, pos))
	}
	parser.match(TtSemi)
	parser.consts[ident] = i
}

func (parser *Parser) Instr() {
	in := Instr{Op: parser.token.Type, Pos: parser.token.Pos}
	parser.match(in.Op)

	switch in.Op {
	case TtKwMov:
		in.Dst = parser.token.Type
		if !parser.matchIf(TtKwOut) {
			in.Dst = parser.reg()
		}
		parser.match(TtComma)
		in.Src = parser.reg()
	case TtKwImm:
		if parser.token.Type == TtInt {
			v, err := strconv.Atoi(parser.token.Lexeme)
			if err != nil {
				parser.matchError(fmt.Sprintf("invalid integer '%s' at %s", parser.token.Lexeme, parser.token.Pos))
			}
			in.Imm = v
			parser.match(TtInt)
		} else {
			in.Ref = parser.token.Lexeme // const or label, resolved later
			parser.match(TtIdent)
		}
	}
	parser.match(TtSemi)
	parser.instrs = append(parser.instrs, in)
}

// resolve fills in every "imm <ident>" and range-checks all immediates.
func (parser *Parser) resolve() {
	for i := range parser.instrs {
		in := &parser.instrs[i]
		if in.Op != TtKwImm {
			continue
		}
		if in.Ref != "" {
			if v, ok := parser.consts[in.Ref]; ok {
				in.Imm = v
			} else if v, ok := parser.labels[in.Ref]; ok {
				in.Imm = v
			} else {
				parser.matchError(fmt.Sprintf("undefined identifier '%s' at %s", in.Ref, in.Pos))
			}
		}
		if in.Imm < 0 || in.Imm > MaxImm {
			parser.matchError(fmt.Sprintf("value %d at %s doesn't fit in a 6-bit immediate", in.Imm, in.Pos))
		}
	}
}
