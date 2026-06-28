// Package repl реализует интерактивный цикл "прочитать строку - вычислить - вывести"
// (read-eval-print loop) для калькулятора.
package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/emityoffwhite/go-calculator/internal/evaluator"
	"github.com/emityoffwhite/go-calculator/internal/parser"
)

const ansVarName = "ans"

// REPL хранит состояние интерактивной сессии: окружение переменных и историю команд.
type REPL struct {
	env     *evaluator.Environment
	history []HistoryEntry
	out     io.Writer
}

// HistoryEntry - одна запись истории вычислений: введённое выражение и результат.
type HistoryEntry struct {
	Input  string
	Result float64
}

// New создаёт новую REPL-сессию, пишущую вывод в out.
func New(out io.Writer) *REPL {
	return &REPL{
		env: evaluator.NewEnvironment(),
		out: out,
	}
}

// History возвращает накопленную историю вычислений.
func (r *REPL) History() []HistoryEntry {
	return r.history
}

// Eval разбирает и вычисляет одну строку ввода, сохраняя результат
// в переменную "ans" и в историю. Не зависит от потоков ввода/вывода,
// что делает его легко тестируемым без эмуляции stdin.
func (r *REPL) Eval(line string) (float64, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return 0, fmt.Errorf("empty input")
	}

	p, err := parser.New(line)
	if err != nil {
		return 0, err
	}

	node, err := p.ParseExpression()
	if err != nil {
		return 0, err
	}

	result, err := evaluator.Eval(node, r.env)
	if err != nil {
		return 0, err
	}

	r.env.Set(ansVarName, result)
	r.history = append(r.history, HistoryEntry{Input: line, Result: result})

	return result, nil
}

// Run запускает интерактивный цикл чтения команд из in до EOF (Ctrl+D / Ctrl+C)
// или команды "exit"/"quit".
func (r *REPL) Run(in io.Reader) {
	scanner := bufio.NewScanner(in)

	fmt.Fprintln(r.out, "Go Calculator. Введите выражение или 'help' для справки, 'exit' для выхода.")

	for {
		fmt.Fprint(r.out, "> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())

		switch line {
		case "":
			continue
		case "exit", "quit":
			return
		case "help":
			r.printHelp()
			continue
		case "history":
			r.printHistory()
			continue
		}

		result, err := r.Eval(line)
		if err != nil {
			fmt.Fprintf(r.out, "ошибка: %v\n", err)
			continue
		}

		fmt.Fprintln(r.out, formatResult(result))
	}
}

func (r *REPL) printHelp() {
	fmt.Fprint(r.out, `Поддерживаемые операции:
  + - * / %  ^         арифметические операторы (^ - возведение в степень)
  ( )                   группировка выражений
  x = 5                 присваивание переменной
  ans                   результат последнего вычисления
  sqrt(x) abs(x)        квадратный корень, модуль
  sin(x) cos(x) tan(x)  тригонометрические функции (радианы)
  log(x) log10(x)       натуральный и десятичный логарифм
  floor(x) ceil(x) round(x)  округление
  pow(x, y) max(x, y) min(x, y)  функции от двух аргументов
  history                показать историю вычислений
  exit / quit            выйти из программы
`)
}

func (r *REPL) printHistory() {
	if len(r.history) == 0 {
		fmt.Fprintln(r.out, "история пуста")
		return
	}
	for i, entry := range r.history {
		fmt.Fprintf(r.out, "%d: %s = %s\n", i+1, entry.Input, formatResult(entry.Result))
	}
}

func formatResult(v float64) string {
	// %g даёт компактное представление: 5 вместо 5.000000, но 3.14 для дробных.
	return fmt.Sprintf("%g", v)
}
