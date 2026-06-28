// Package evaluator вычисляет значение AST-дерева, построенного парсером.
package evaluator

import (
	"fmt"
	"math"

	"github.com/emityoffwhite/go-calculator/internal/ast"
)

// builtinFunc - сигнатура встроенной функции калькулятора.
type builtinFunc func(args []float64) (float64, error)

// builtins содержит все доступные в калькуляторе функции.
// Объявлены как переменная пакета (не константа), что позволяет легко
// добавлять новые функции без изменения логики Eval.
var builtins = map[string]builtinFunc{
	"sqrt":  unary(math.Sqrt),
	"abs":   unary(math.Abs),
	"sin":   unary(math.Sin),
	"cos":   unary(math.Cos),
	"tan":   unary(math.Tan),
	"log":   unary(math.Log), // натуральный логарифм
	"log10": unary(math.Log10),
	"floor": unary(math.Floor),
	"ceil":  unary(math.Ceil),
	"round": unary(math.Round),
	"pow": func(args []float64) (float64, error) {
		if len(args) != 2 {
			return 0, fmt.Errorf("pow expects 2 arguments, got %d", len(args))
		}
		return math.Pow(args[0], args[1]), nil
	},
	"max": func(args []float64) (float64, error) {
		if len(args) != 2 {
			return 0, fmt.Errorf("max expects 2 arguments, got %d", len(args))
		}
		return math.Max(args[0], args[1]), nil
	},
	"min": func(args []float64) (float64, error) {
		if len(args) != 2 {
			return 0, fmt.Errorf("min expects 2 arguments, got %d", len(args))
		}
		return math.Min(args[0], args[1]), nil
	},
}

// unary адаптирует функцию вида func(float64) float64 (как math.Sqrt)
// к сигнатуре builtinFunc, проверяя количество аргументов.
func unary(f func(float64) float64) builtinFunc {
	return func(args []float64) (float64, error) {
		if len(args) != 1 {
			return 0, fmt.Errorf("expected 1 argument, got %d", len(args))
		}
		return f(args[0]), nil
	}
}

// Environment хранит переменные между вычислениями в рамках одной REPL-сессии.
type Environment struct {
	vars map[string]float64
}

// NewEnvironment создаёт пустое окружение переменных.
func NewEnvironment() *Environment {
	return &Environment{vars: make(map[string]float64)}
}

// Get возвращает значение переменной и флаг её существования.
func (e *Environment) Get(name string) (float64, bool) {
	v, ok := e.vars[name]
	return v, ok
}

// Set устанавливает значение переменной.
func (e *Environment) Set(name string, value float64) {
	e.vars[name] = value
}

// Eval вычисляет значение AST-узла в контексте переданного окружения.
func Eval(node ast.Node, env *Environment) (float64, error) {
	switch n := node.(type) {
	case *ast.NumberLiteral:
		return n.Value, nil

	case *ast.Identifier:
		v, ok := env.Get(n.Name)
		if !ok {
			return 0, fmt.Errorf("undefined variable: %s", n.Name)
		}
		return v, nil

	case *ast.UnaryExpr:
		return evalUnary(n, env)

	case *ast.BinaryExpr:
		return evalBinary(n, env)

	case *ast.CallExpr:
		return evalCall(n, env)

	case *ast.AssignExpr:
		value, err := Eval(n.Value, env)
		if err != nil {
			return 0, err
		}
		env.Set(n.Name, value)
		return value, nil

	default:
		return 0, fmt.Errorf("unknown AST node: %T", node)
	}
}

func evalUnary(n *ast.UnaryExpr, env *Environment) (float64, error) {
	v, err := Eval(n.Operand, env)
	if err != nil {
		return 0, err
	}
	switch n.Op {
	case "-":
		return -v, nil
	case "+":
		return v, nil
	default:
		return 0, fmt.Errorf("unknown unary operator: %s", n.Op)
	}
}

func evalBinary(n *ast.BinaryExpr, env *Environment) (float64, error) {
	left, err := Eval(n.Left, env)
	if err != nil {
		return 0, err
	}
	right, err := Eval(n.Right, env)
	if err != nil {
		return 0, err
	}

	switch n.Op {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	case "%":
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return math.Mod(left, right), nil
	case "^":
		return math.Pow(left, right), nil
	default:
		return 0, fmt.Errorf("unknown binary operator: %s", n.Op)
	}
}

func evalCall(n *ast.CallExpr, env *Environment) (float64, error) {
	fn, ok := builtins[n.Function]
	if !ok {
		return 0, fmt.Errorf("unknown function: %s", n.Function)
	}

	args := make([]float64, len(n.Args))
	for i, argNode := range n.Args {
		v, err := Eval(argNode, env)
		if err != nil {
			return 0, err
		}
		args[i] = v
	}

	result, err := fn(args)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", n.Function, err)
	}
	return result, nil
}
