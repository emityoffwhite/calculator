package parser

import "testing"

// TestParseExpression_String проверяет правильность приоритетов и ассоциативности
// через строковое представление AST. Например, для "3 + 4 * 2" корректный AST
// должен быть "(3 + (4 * 2))" - умножение связывает операнды теснее, чем сложение.
func TestParseExpression_String(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"3 + 4", "(3 + 4)"},
		{"3 + 4 * 2", "(3 + (4 * 2))"},
		{"(3 + 4) * 2", "((3 + 4) * 2)"},
		{"3 - 4 - 2", "((3 - 4) - 2)"}, // левая ассоциативность
		{"2 ^ 3 ^ 2", "(2 ^ (3 ^ 2))"}, // правая ассоциативность
		{"-5 + 3", "((-5) + 3)"},
		{"2 * -3", "(2 * (-3))"},
		{"10 % 3", "(10 % 3)"},
		{"x = 5", "x = 5"},
		{"sqrt(16)", "sqrt(16)"},
		{"pow(2, 10)", "pow(2, 10)"},
		{"2 + 3 * 4 - 1", "((2 + (3 * 4)) - 1)"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p, err := New(tt.input)
			if err != nil {
				t.Fatalf("New() unexpected error: %v", err)
			}

			node, err := p.ParseExpression()
			if err != nil {
				t.Fatalf("ParseExpression() unexpected error: %v", err)
			}

			if got := node.String(); got != tt.want {
				t.Errorf("AST = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseExpression_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unclosed paren", "(3 + 4"},
		{"trailing garbage", "2 + 3 4"},
		{"missing operand", "3 +"},
		{"empty parens for binary", "()"},
		{"unexpected operator at start", "* 3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New(tt.input)
			if err != nil {
				// Ошибка на этапе создания парсера (например, lexer-ошибка) тоже валидна.
				return
			}

			_, err = p.ParseExpression()
			if err == nil {
				t.Errorf("ParseExpression() expected error for input %q, got nil", tt.input)
			}
		})
	}
}

func TestParseExpression_MultiArgFunction(t *testing.T) {
	p, err := New("max(1, 2)")
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	node, err := p.ParseExpression()
	if err != nil {
		t.Fatalf("ParseExpression() unexpected error: %v", err)
	}

	want := "max(1, 2)"
	if got := node.String(); got != want {
		t.Errorf("AST = %q, want %q", got, want)
	}
}
