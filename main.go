package main

import (
	"encoding/json"
	"fmt"
	"qles/qles"

	"github.com/xwb1989/sqlparser"
)

type constrain interface {
	int | float32
}

func add[T constrain](a T) T {
	return a + a
}

// fix dot convertion shit
func main() {

	query := "SELECT * FROM Temp WHERE path.b IN (7, 8, 9) ORDER BY test"
	ast, _ := qles.BuildSQL(query)
	selectStatement := ast.(*sqlparser.Select)
	tables := selectStatement.From
	for _, t := range tables {
		fmt.Println(t)
	}

	orderby := selectStatement.OrderBy
	for _, o := range orderby {
		fmt.Println(o.Direction, o.Expr.(*sqlparser.ColName).Name.String())
	}

	reverse := qles.ReverseNot(ast)
	boolq := qles.BuildESQuery(reverse, map[string]string{"a": "", "b": "path"})
	_, _ = json.Marshal(boolq)
	// fmt.Println(string(j))
}

// func Traverse(node sqlparser.SQLNode, depth int) {
// 	indent := strings.Repeat("\t", depth)
// 	switch n := node.(type) {
// 	case *sqlparser.Select:
// 		fmt.Println(indent, "SELECT")
// 		Traverse(n.From, depth+1)
// 		Traverse(n.Where, depth+1)
// 	case *sqlparser.TableName:
// 		fmt.Println(indent, "Table Name:", n.Name)
// 	case *sqlparser.Where:
// 		fmt.Println(indent, "WHERE")
// 		fmt.Printf("the type of where child is %T", n.Expr)
// 		Traverse(n.Expr, depth+1)

// 	case *sqlparser.AndExpr:
// 		fmt.Println(indent, "And Expr:", n.Left, n.Right)
// 		Traverse(n.Left, depth+1)
// 		Traverse(n.Right, depth+1)

// 	case *sqlparser.OrExpr:
// 		fmt.Println(indent, "Or Expr:", n.Left, n.Right)
// 		Traverse(n.Left, depth+1)
// 		Traverse(n.Right, depth+1)

// 	case *sqlparser.ComparisonExpr:
// 		fmt.Println(indent, "Comparison Expr:", n.Operator)
// 		fmt.Println(indent, n.Right)
// 		fmt.Printf("%T", n.Right)
// 		// Traverse(n.Left, depth+1)
// 		// Traverse(n.Right, depth+1)

// 	case *sqlparser.ParenExpr:
// 		fmt.Println(indent, "Paren", "(")
// 		Traverse(n.Expr, depth+1)
// 	}

// }
