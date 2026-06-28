package evaluator

import (
	"math"
	"testing"

	"github.com/emityoffwhite/go-calculator/internal/parser"
)

// evalString - тестовый хелпер, объединяющий парсинг и вычисление в одну строку,
// чтобы не дублировать parser.New + ParseExpression в каждом тест-кейсе.
func evalString(t *testing.T, input string, env *Environment) float64 {
	t.Helper()

	p, err := parser.New(input)
	if err != nil {
		t.Fatalf("parser.New(%q) unexpected error: %v", input, err)
	}

	node, err := p.ParseExpression()
	if err != nil {
		t.Fatalf("ParseExpression(%q) unexpected error: %v", input, err)
	}

	result, err := Eval(node, env)
	if err != nil {
		t.Fatalf("Eval(%q) unexpected error: %v", input, err)
	}

	return result
}

func TestEval_Arithmetic(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"3 + 4", 7},
		{"3 + 4 * 2", 11},
		{"(3 + 4) * 2", 14},
		{"10 - 3 - 2", 5},
		{"2 ^ 10", 1024},
		{"2 ^ 3 ^ 2", 512}, // право-ассоциативность: 2^(3^2) = 2^9
		{"10 / 4", 2.5},
		{"10 % 3", 1},
		{"-5 + 3", -2},
		{"-(2 + 3)", -5},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := evalString(t, tt.input, NewEnvironment())
			if got != tt.want {
				t.Errorf("Eval(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestEval_DivisionByZero(t *testing.T) {
	tests := []string{"5 / 0", "5 % 0"}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			p, err := parser.New(input)
			if err != nil {
				t.Fatalf("unexpected parser error: %v", err)
			}
			node, err := p.ParseExpression()
			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}

			_, err = Eval(node, NewEnvironment())
			if err == nil {
				t.Fatalf("Eval(%q) expected division by zero error, got nil", input)
			}
		})
	}
}

func TestEval_Variables(t *testing.T) {
	env := NewEnvironment()

	got := evalString(t, "x = 5", env)
	if got != 5 {
		t.Fatalf("Eval(x = 5) = %v, want 5", got)
	}

	got = evalString(t, "x * 2 + 1", env)
	if got != 11 {
		t.Fatalf("Eval(x * 2 + 1) = %v, want 11", got)
	}
}

func TestEval_UndefinedVariable(t *testing.T) {
	p, err := parser.New("y + 1")
	if err != nil {
		t.Fatalf("unexpected parser error: %v", err)
	}
	node, err := p.ParseExpression()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	_, err = Eval(node, NewEnvironment())
	if err == nil {
		t.Fatal("expected error for undefined variable, got nil")
	}
}

func TestEval_BuiltinFunctions(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"sqrt(16)", 4},
		{"abs(-5)", 5},
		{"pow(2, 10)", 1024},
		{"max(3, 7)", 7},
		{"min(3, 7)", 3},
		{"floor(3.7)", 3},
		{"ceil(3.2)", 4},
		{"round(3.5)", 4},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := evalString(t, tt.input, NewEnvironment())
			if got != tt.want {
				t.Errorf("Eval(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestEval_TrigFunctions(t *testing.T) {
	// sin(0) = 0, cos(0) = 1 - проверяем с допуском на погрешность float.
	got := evalString(t, "sin(0)", NewEnvironment())
	if math.Abs(got-0) > 1e-9 {
		t.Errorf("sin(0) = %v, want 0", got)
	}

	got = evalString(t, "cos(0)", NewEnvironment())
	if math.Abs(got-1) > 1e-9 {
		t.Errorf("cos(0) = %v, want 1", got)
	}
}

func TestEval_UnknownFunction(t *testing.T) {
	p, err := parser.New("foobar(1)")
	if err != nil {
		t.Fatalf("unexpected parser error: %v", err)
	}
	node, err := p.ParseExpression()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	_, err = Eval(node, NewEnvironment())
	if err == nil {
		t.Fatal("expected error for unknown function, got nil")
	}
}

func TestEval_WrongArgCount(t *testing.T) {
	tests := []string{"sqrt(1, 2)", "pow(2)"}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			p, err := parser.New(input)
			if err != nil {
				t.Fatalf("unexpected parser error: %v", err)
			}
			node, err := p.ParseExpression()
			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}

			_, err = Eval(node, NewEnvironment())
			if err == nil {
				t.Errorf("Eval(%q) expected error for wrong argument count, got nil", input)
			}
		})
	}
}

func TestEnvironment_GetSet(t *testing.T) {
	env := NewEnvironment()

	if _, ok := env.Get("missing"); ok {
		t.Error("Get() on empty environment should return ok=false")
	}

	env.Set("x", 42)
	v, ok := env.Get("x")
	if !ok || v != 42 {
		t.Errorf("Get(x) = (%v, %v), want (42, true)", v, ok)
	}
}
