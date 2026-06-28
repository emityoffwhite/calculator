package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/emityoffwhite/go-calculator/internal/evaluator"
	"github.com/emityoffwhite/go-calculator/internal/parser"
	"github.com/emityoffwhite/go-calculator/internal/repl"
)

func main() {
	exprFlag := flag.String("expr", "", "вычислить одно выражение и выйти, например --expr \"2 + 2 * 2\"")
	flag.Parse()

	if *exprFlag != "" {
		if err := evalOnce(*exprFlag); err != nil {
			fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
			os.Exit(1)
		}
		return
	}

	r := repl.New(os.Stdout)
	r.Run(os.Stdin)
}

// evalOnce вычисляет одно выражение без запуска REPL - удобно для скриптов
// и пайпов: calc --expr "2+2" или результат можно сразу использовать в shell.
func evalOnce(expr string) error {
	p, err := parser.New(expr)
	if err != nil {
		return err
	}

	node, err := p.ParseExpression()
	if err != nil {
		return err
	}

	env := evaluator.NewEnvironment()
	result, err := evaluator.Eval(node, env)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}
