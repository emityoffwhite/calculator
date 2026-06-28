package lexer

import "testing"

func TestNextToken(t *testing.T) {
	input := `3 + 4.5 * (2 - 1) / x ^ 2 % 3, y = sqrt(16)`

	expected := []Token{
		{Type: NUMBER, Literal: "3"},
		{Type: PLUS, Literal: "+"},
		{Type: NUMBER, Literal: "4.5"},
		{Type: ASTERISK, Literal: "*"},
		{Type: LPAREN, Literal: "("},
		{Type: NUMBER, Literal: "2"},
		{Type: MINUS, Literal: "-"},
		{Type: NUMBER, Literal: "1"},
		{Type: RPAREN, Literal: ")"},
		{Type: SLASH, Literal: "/"},
		{Type: IDENT, Literal: "x"},
		{Type: CARET, Literal: "^"},
		{Type: NUMBER, Literal: "2"},
		{Type: PERCENT, Literal: "%"},
		{Type: NUMBER, Literal: "3"},
		{Type: COMMA, Literal: ","},
		{Type: IDENT, Literal: "y"},
		{Type: ASSIGN, Literal: "="},
		{Type: IDENT, Literal: "sqrt"},
		{Type: LPAREN, Literal: "("},
		{Type: NUMBER, Literal: "16"},
		{Type: RPAREN, Literal: ")"},
		{Type: EOF, Literal: ""},
	}

	l := New(input)

	for i, want := range expected {
		got, err := l.NextToken()
		if err != nil {
			t.Fatalf("token %d: unexpected error: %v", i, err)
		}
		if got.Type != want.Type {
			t.Fatalf("token %d: type = %v, want %v", i, got.Type, want.Type)
		}
		if got.Literal != want.Literal {
			t.Fatalf("token %d: literal = %q, want %q", i, got.Literal, want.Literal)
		}
	}
}

func TestNextToken_InvalidCharacter(t *testing.T) {
	l := New("3 & 4")

	// Считываем '3', пробел игнорируется автоматически.
	if _, err := l.NextToken(); err != nil {
		t.Fatalf("unexpected error reading first token: %v", err)
	}

	_, err := l.NextToken()
	if err == nil {
		t.Fatal("expected error for unsupported character '&', got nil")
	}
}

func TestNextToken_MalformedNumber(t *testing.T) {
	l := New("3.")

	_, err := l.NextToken()
	if err == nil {
		t.Fatal("expected error for malformed number '3.', got nil")
	}
}

func TestNextToken_EmptyInput(t *testing.T) {
	l := New("")

	tok, err := l.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != EOF {
		t.Errorf("Type = %v, want EOF", tok.Type)
	}
}

func TestTokenType_String(t *testing.T) {
	if PLUS.String() != "+" {
		t.Errorf("PLUS.String() = %q, want %q", PLUS.String(), "+")
	}
	unknown := TokenType(999)
	if unknown.String() != "UNKNOWN" {
		t.Errorf("unknown type String() = %q, want %q", unknown.String(), "UNKNOWN")
	}
}
