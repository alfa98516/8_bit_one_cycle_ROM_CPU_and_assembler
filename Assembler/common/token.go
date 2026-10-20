package common

import (
	"fmt"
	"strconv"
)

type (
	TokenType int
	Position  struct {
		Line int
		Col  int
	}
)

func (p Position) String() string {
	return fmt.Sprintf("(%d %d)", p.Line, p.Col)
}

type Token struct {
	Type   TokenType
	Lexeme string
	Pos    Position
	PosEnd Position
}

type Error struct {
	Pos Position
	Msg string
}

func (e Error) string() string {
	return fmt.Sprintf("Error: '%s' '%s'.", e.Pos, e.Msg)
}

type ErrorHandler func(pos Position, msg string)

const (
	TtEOI TokenType = iota
	TtUnknown
	TtInt
	TtIdent
	TtLabel
	TtComma
	TtSemi
	beginKeyword
	TtKwNand
	TtKwOr
	TtKwAnd
	TtKwNor
	TtKwAdd
	TtKwSub
	TtKwXor
	TtKwNot
	TtKwMov
	TtKwImm
	TtKwNop
	TtKwConst
	TtKwJmp
	TtKwJz
	TtKwJnz
	TtKwJl
	TtKwJge
	TtKwJle
	TtKwJg
	TtKwOut
	TtKwA
	TtKwB
	TtKwC
	TtKwD
	TtKwE
	TtKwF
	endKeyword
)

var tokenTypes = [...]string{
	TtEOI:     "EOI",
	TtUnknown: "UNKNOWN",
	TtInt:     "INT",
	TtIdent:   "IDENT",
	TtLabel:   "LABEL",
	TtComma:   ",",
	TtSemi:    ";",
	TtKwNand:  "nand",
	TtKwOr:    "or",
	TtKwAnd:   "and",
	TtKwNor:   "nor",
	TtKwAdd:   "add",
	TtKwSub:   "sub",
	TtKwXor:   "xor",
	TtKwNot:   "not",
	TtKwMov:   "mov",
	TtKwImm:   "imm",
	TtKwNop:   "nop",
	TtKwConst: "const",
	TtKwJmp:   "jmp",
	TtKwJz:    "jz",
	TtKwJnz:   "jnz",
	TtKwJl:    "jl",
	TtKwJge:   "jge",
	TtKwJle:   "jle",
	TtKwJg:    "jg",
	TtKwA:     "A",
	TtKwB:     "B",
	TtKwC:     "C",
	TtKwD:     "D",
	TtKwE:     "E",
	TtKwF:     "F",
	TtKwOut:   "out",
}

var keywords map[string]TokenType

func (Tt TokenType) String() string {
	s := ""
	if 0 <= Tt && Tt < TokenType(len(tokenTypes)) {
		s = tokenTypes[Tt]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(Tt)) + ")"
	}
	return s
}

func Lookup(ident string) TokenType {
	if tok, isKeyword := keywords[ident]; isKeyword {
		return tok
	}
	return TtIdent
}

func init() {
	keywords = make(map[string]TokenType, endKeyword-(beginKeyword+1))
	for i := beginKeyword + 1; i < endKeyword; i++ {
		keywords[tokenTypes[i]] = i
	}
}
