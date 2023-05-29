package main

import (
	"fmt"
	"qles/qles"
	"strings"

	"github.com/xwb1989/sqlparser"
)

type foo interface {
	ifoo()
}

func (goo) ifoo() {}

type goo struct {
}

func test(foo) {

}

func main() {
	query := "SELECT * FROM Temp WHERE a IN (1, 2, 3, 4, 5)"
	_, _ = qles.Build(query)
	tree, _ := sqlparser.Parse(query)
	s := tree.(*sqlparser.Select)

	Traverse(s, 1)

	tmp := sqlparser.SQLVal{}
	tmp.Val = []byte("1")
	tmp.Type = sqlparser.StrVal

	test(goo{})

}

func Traverse(node sqlparser.SQLNode, depth int) {
	indent := strings.Repeat("\t", depth)
	switch n := node.(type) {
	case *sqlparser.Select:
		fmt.Println(indent, "SELECT")
		Traverse(n.From, depth+1)
		Traverse(n.Where, depth+1)
	case *sqlparser.TableName:
		fmt.Println(indent, "Table Name:", n.Name)
	case *sqlparser.Where:
		fmt.Println(indent, "WHERE")
		fmt.Printf("the type of where child is %T", n.Expr)
		Traverse(n.Expr, depth+1)

	case *sqlparser.AndExpr:
		fmt.Println(indent, "And Expr:", n.Left, n.Right)
		Traverse(n.Left, depth+1)
		Traverse(n.Right, depth+1)

	case *sqlparser.OrExpr:
		fmt.Println(indent, "Or Expr:", n.Left, n.Right)
		Traverse(n.Left, depth+1)
		Traverse(n.Right, depth+1)

	case *sqlparser.ComparisonExpr:
		fmt.Println(indent, "Comparison Expr:", n.Operator)
		fmt.Println(indent, n.Right)
		fmt.Printf("%T", n.Right)
		// Traverse(n.Left, depth+1)
		// Traverse(n.Right, depth+1)

	case *sqlparser.ParenExpr:
		fmt.Println(indent, "Paren", "(")
		Traverse(n.Expr, depth+1)
	}

}
