// Package ast определяет узлы дерева разбора (Abstract Syntax Tree)
// для арифметических выражений калькулятора.
package ast

// Node - общий интерфейс для всех узлов дерева разбора.
type Node interface {
	String() string
}

// NumberLiteral представляет числовой литерал, например 3.14.
type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) String() string { return formatFloat(n.Value) }

// Identifier представляет имя переменной, например x.
type Identifier struct {
	Name string
}

func (i *Identifier) String() string { return i.Name }

// BinaryExpr представляет бинарную операцию: левый операнд, оператор, правый операнд.
// Например: 3 + 4 -> BinaryExpr{Left: 3, Op: "+", Right: 4}
type BinaryExpr struct {
	Left  Node
	Op    string
	Right Node
}

func (b *BinaryExpr) String() string {
	return "(" + b.Left.String() + " " + b.Op + " " + b.Right.String() + ")"
}

// UnaryExpr представляет унарную операцию: -5, +3.
type UnaryExpr struct {
	Op      string
	Operand Node
}

func (u *UnaryExpr) String() string {
	return "(" + u.Op + u.Operand.String() + ")"
}

// CallExpr представляет вызов функции: sqrt(16), pow(2, 10).
type CallExpr struct {
	Function string
	Args     []Node
}

func (c *CallExpr) String() string {
	s := c.Function + "("
	for i, arg := range c.Args {
		if i > 0 {
			s += ", "
		}
		s += arg.String()
	}
	return s + ")"
}

// AssignExpr представляет присваивание переменной: x = 5.
type AssignExpr struct {
	Name  string
	Value Node
}

func (a *AssignExpr) String() string {
	return a.Name + " = " + a.Value.String()
}
