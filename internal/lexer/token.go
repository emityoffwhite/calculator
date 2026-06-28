package lexer

// TokenType определяет вид токена, полученного из исходной строки.
type TokenType int

const (
	EOF TokenType = iota
	NUMBER
	IDENT // имена переменных и функций: x, sqrt, ans

	PLUS     // +
	MINUS    // -
	ASTERISK // *
	SLASH    // /
	CARET    // ^ (возведение в степень)
	PERCENT  // % (остаток от деления)

	LPAREN // (
	RPAREN // )
	COMMA  // , (разделитель аргументов функций)
	ASSIGN // =
)

// Token представляет собой одну лексическую единицу: её тип, исходный текст
// и позицию в строке (для информативных сообщений об ошибках).
type Token struct {
	Type    TokenType
	Literal string
	Pos     int
}

func (t TokenType) String() string {
	names := map[TokenType]string{
		EOF:      "EOF",
		NUMBER:   "NUMBER",
		IDENT:    "IDENT",
		PLUS:     "+",
		MINUS:    "-",
		ASTERISK: "*",
		SLASH:    "/",
		CARET:    "^",
		PERCENT:  "%",
		LPAREN:   "(",
		RPAREN:   ")",
		COMMA:    ",",
		ASSIGN:   "=",
	}
	if name, ok := names[t]; ok {
		return name
	}
	return "UNKNOWN"
}
