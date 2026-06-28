package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestREPL_Eval(t *testing.T) {
	r := New(&bytes.Buffer{})

	result, err := r.Eval("3 + 4")
	if err != nil {
		t.Fatalf("Eval() unexpected error: %v", err)
	}
	if result != 7 {
		t.Errorf("Eval() = %v, want 7", result)
	}
}

func TestREPL_Eval_EmptyInput(t *testing.T) {
	r := New(&bytes.Buffer{})

	_, err := r.Eval("   ")
	if err == nil {
		t.Error("Eval() on empty input should return error")
	}
}

func TestREPL_AnsVariable(t *testing.T) {
	r := New(&bytes.Buffer{})

	if _, err := r.Eval("5 + 5"); err != nil {
		t.Fatalf("first Eval() unexpected error: %v", err)
	}

	result, err := r.Eval("ans * 2")
	if err != nil {
		t.Fatalf("second Eval() unexpected error: %v", err)
	}
	if result != 20 {
		t.Errorf("ans * 2 = %v, want 20 (ans should hold previous result 10)", result)
	}
}

func TestREPL_History(t *testing.T) {
	r := New(&bytes.Buffer{})

	if _, err := r.Eval("1 + 1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := r.Eval("2 + 2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	history := r.History()
	if len(history) != 2 {
		t.Fatalf("len(History()) = %d, want 2", len(history))
	}
	if history[0].Input != "1 + 1" || history[0].Result != 2 {
		t.Errorf("history[0] = %+v, want {Input: \"1 + 1\", Result: 2}", history[0])
	}
	if history[1].Input != "2 + 2" || history[1].Result != 4 {
		t.Errorf("history[1] = %+v, want {Input: \"2 + 2\", Result: 4}", history[1])
	}
}

func TestREPL_Eval_DoesNotRecordFailedAttempts(t *testing.T) {
	r := New(&bytes.Buffer{})

	if _, err := r.Eval("3 + 4"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := r.Eval("1 / 0"); err == nil {
		t.Fatal("expected division by zero error")
	}

	history := r.History()
	if len(history) != 1 {
		t.Errorf("len(History()) = %d, want 1 (failed eval should not be recorded)", len(history))
	}
}

// TestREPL_Run_FullSession проверяет весь интерактивный цикл целиком:
// ввод нескольких команд через io.Reader и проверка итогового вывода.
// Это интеграционный тест, имитирующий реального пользователя за терминалом.
func TestREPL_Run_FullSession(t *testing.T) {
	input := strings.NewReader("x = 5\nx * 2 + 1\nans + 10\nexit\n")
	var out bytes.Buffer

	r := New(&out)
	r.Run(input)

	output := out.String()

	wantSubstrings := []string{"5", "11", "21"}
	for _, want := range wantSubstrings {
		if !strings.Contains(output, want) {
			t.Errorf("output missing expected substring %q\nfull output:\n%s", want, output)
		}
	}
}

func TestREPL_Run_HandlesErrors(t *testing.T) {
	input := strings.NewReader("1 / 0\nexit\n")
	var out bytes.Buffer

	r := New(&out)
	r.Run(input)

	if !strings.Contains(out.String(), "ошибка") {
		t.Errorf("expected error message in output, got:\n%s", out.String())
	}
}

func TestREPL_Run_HelpCommand(t *testing.T) {
	input := strings.NewReader("help\nexit\n")
	var out bytes.Buffer

	r := New(&out)
	r.Run(input)

	if !strings.Contains(out.String(), "sqrt") {
		t.Errorf("expected help text to mention sqrt, got:\n%s", out.String())
	}
}

func TestREPL_Run_HistoryCommand(t *testing.T) {
	input := strings.NewReader("1 + 1\nhistory\nexit\n")
	var out bytes.Buffer

	r := New(&out)
	r.Run(input)

	if !strings.Contains(out.String(), "1 + 1 = 2") {
		t.Errorf("expected history output to contain '1 + 1 = 2', got:\n%s", out.String())
	}
}
