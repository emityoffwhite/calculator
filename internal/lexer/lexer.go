package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

// Lexer превращает строку исходного выражения в последовательность токенов.
// Работает по принципу "один символ за раз" (one-pass scanner), что типично
// для простых языковых процессоров: O(n) по времени, без бэктрекинга.
type Lexer struct {
	input        string
	position     int  // индекс текущего символа
	readPosition int  // индекс следующего символа для чтения
	ch           byte // текущий рассматриваемый символ
}

// New создаёт новый Lexer для переданной строки.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // 0 означает "конец входной строки" (NUL-байт как сторожевое значение)
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken считывает и возвращает следующий токен входной строки.
func (l *Lexer) NextToken() (Token, error) {
	l.skipWhitespace()

	pos := l.position

	switch {
	case l.ch == 0:
		return Token{Type: EOF, Literal: "", Pos: pos}, nil
	case l.ch == '+':
		l.readChar()
		return Token{Type: PLUS, Literal: "+", Pos: pos}, nil
	case l.ch == '-':
		l.readChar()
		return Token{Type: MINUS, Literal: "-", Pos: pos}, nil
	case l.ch == '*':
		l.readChar()
		return Token{Type: ASTERISK, Literal: "*", Pos: pos}, nil
	case l.ch == '/':
		l.readChar()
		return Token{Type: SLASH, Literal: "/", Pos: pos}, nil
	case l.ch == '^':
		l.readChar()
		return Token{Type: CARET, Literal: "^", Pos: pos}, nil
	case l.ch == '%':
		l.readChar()
		return Token{Type: PERCENT, Literal: "%", Pos: pos}, nil
	case l.ch == '(':
		l.readChar()
		return Token{Type: LPAREN, Literal: "(", Pos: pos}, nil
	case l.ch == ')':
		l.readChar()
		return Token{Type: RPAREN, Literal: ")", Pos: pos}, nil
	case l.ch == ',':
		l.readChar()
		return Token{Type: COMMA, Literal: ",", Pos: pos}, nil
	case l.ch == '=':
		l.readChar()
		return Token{Type: ASSIGN, Literal: "=", Pos: pos}, nil
	case isDigit(l.ch):
		return l.readNumber(pos)
	case isLetter(l.ch):
		return l.readIdentifier(pos)
	default:
		ch := l.ch
		l.readChar()
		return Token{}, fmt.Errorf("unexpected character %q at position %d", ch, pos)
	}
}

func (l *Lexer) readNumber(startPos int) (Token, error) {
	var sb strings.Builder
	hasDot := false

	for isDigit(l.ch) || (l.ch == '.' && !hasDot) {
		if l.ch == '.' {
			hasDot = true
		}
		sb.WriteByte(l.ch)
		l.readChar()
	}

	literal := sb.String()
	if strings.HasSuffix(literal, ".") {
		return Token{}, fmt.Errorf("malformed number %q at position %d", literal, startPos)
	}

	return Token{Type: NUMBER, Literal: literal, Pos: startPos}, nil
}

func (l *Lexer) readIdentifier(startPos int) (Token, error) {
	var sb strings.Builder

	for isLetter(l.ch) || isDigit(l.ch) {
		sb.WriteByte(l.ch)
		l.readChar()
	}

	return Token{Type: IDENT, Literal: sb.String(), Pos: startPos}, nil
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
