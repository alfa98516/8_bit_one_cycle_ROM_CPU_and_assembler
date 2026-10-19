package frontend

import (
	"strings"
	"unicode"

	. "Assembler/common"
)

type Lexer struct {
	src          []rune
	n            int
	i            int
	ch           rune
	eoi          bool
	token        Token
	pos          Position
	errorHandler ErrorHandler
}

func (lexer *Lexer) peekRune() (rune, bool) {
	if lexer.eoi || lexer.i+1 >= lexer.n {
		return rune(0), true
	}
	return lexer.src[lexer.i+1], false
}

func (lexer *Lexer) peekRuneIs(check rune) bool {
	r, err := lexer.peekRune()
	return !err && r == check
}

func (lexer *Lexer) nextRune() {
	if lexer.eoi {
		return
	}
	lexer.i++
	if lexer.i >= lexer.n {
		lexer.ch, lexer.eoi = rune(0), true
		return
	}
	ch := lexer.src[lexer.i]
	if ch == '\n' { // Note: '\r'
		lexer.pos.Line += 1
		lexer.pos.Col = -1
	}
	lexer.pos.Col += 1
	lexer.ch, lexer.eoi = ch, false
}

func (lexer *Lexer) setToken(Tt TokenType, r rune) {
	lexer.token = Token{Type: Tt, Lexeme: string(lexer.ch), Pos: lexer.pos}
	lexer.nextRune()
	lexer.token.PosEnd = lexer.pos
}

func (lexer *Lexer) identifier() string {
	var builder strings.Builder
	for {
		builder.WriteRune(lexer.ch)
		nextR, _ := lexer.peekRune()
		if unicode.IsLetter(nextR) || nextR == '_' || unicode.IsDigit(nextR) {
			lexer.nextRune()
		} else {
			break
		}
	}
	return builder.String()
}

func (lexer *Lexer) integer(pos Position) Token {
	builder := strings.Builder{}
	for !lexer.eoi && unicode.IsDigit(lexer.ch) {
		builder.WriteRune(lexer.ch)
		lexer.nextRune()
	}
	return Token{Type: TtInt, Lexeme: builder.String(), Pos: pos}
}

func (lexer *Lexer) resetToken(Tt TokenType) {
	lexer.token = Token{Type: Tt, Lexeme: lexer.token.Lexeme + string(lexer.ch), Pos: lexer.token.Pos}
	lexer.nextRune()
	lexer.token.PosEnd = lexer.pos
}

func (lexer *Lexer) processWhiteSpaces() bool {
	n := 0
	for !lexer.eoi && unicode.IsSpace(lexer.ch) {
		n += 1
		lexer.nextRune()
	}
	return n > 0
}

func (lexer *Lexer) Init(src []byte, handler ErrorHandler) {
	str := []rune(string(src))                               // Fix: cast to string
	lexer.src, lexer.n, lexer.i = []rune(str), len(str), -1. // Fix: len of str
	lexer.eoi = false
	lexer.pos = Position{Line: 1, Col: 0}
	lexer.errorHandler = handler
	lexer.nextRune()
}

func (lexer *Lexer) NextToken() Token {
	for lexer.processWhiteSpaces() {
	}
	if lexer.eoi {
		lexer.token = Token{Type: TtEOI, Lexeme: TtEOI.String(), Pos: lexer.pos}
		return lexer.token
	}

	for lexer.ch == '#' {
		for lexer.ch != '\n' {
			lexer.nextRune()
		}
		lexer.nextRune()
	}
	for lexer.processWhiteSpaces() {
	}
	switch lexer.ch {
	case ',':
		lexer.setToken(TtComma, ',')
	case ';':
		lexer.setToken(TtSemi, ';')
	default:
		if unicode.IsLetter(lexer.ch) || lexer.ch == '_' {
			ps := lexer.pos
			lexeme := lexer.identifier()
			lexer.token = Token{Type: Lookup(lexeme), Lexeme: lexeme, Pos: ps}
			lexer.nextRune()
		} else if unicode.IsDigit(lexer.ch) {
			lexer.token = lexer.integer(lexer.pos)
		} else {
			lexer.token = Token{Type: TtUnknown, Lexeme: "UNKNOWN", Pos: lexer.pos}
		}
	}
	return lexer.token
}
