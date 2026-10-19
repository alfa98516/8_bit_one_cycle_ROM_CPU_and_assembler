package common

type LexerInterface interface {
	Init(src []byte, handler ErrorHandler)
	NextToken() Token
}

type ParserInterface interface {
	ParseSrc(src []byte, handler ErrorHandler) string
}
