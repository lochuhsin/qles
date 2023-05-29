package qles

import (
	"github.com/xwb1989/sqlparser"
)

func reverseTree(node sqlparser.Expr, isNot bool) sqlparser.Expr {
	switch n := node.(type) {

	case *sqlparser.AndExpr:
		l, r := reverseTree(n.Left, isNot), reverseTree(n.Right, isNot)
		if isNot {
			return &sqlparser.OrExpr{Left: l, Right: r}
		}
		return &sqlparser.AndExpr{Left: l, Right: r}

	case *sqlparser.OrExpr:
		l, r := reverseTree(n.Left, isNot), reverseTree(n.Right, isNot)
		if isNot {
			return &sqlparser.AndExpr{Left: l, Right: r}
		}
		return &sqlparser.OrExpr{Left: l, Right: r}

	case *sqlparser.ParenExpr:
		return &sqlparser.ParenExpr{Expr: reverseTree(n.Expr, isNot)}

	case *sqlparser.NotExpr:
		return reverseTree(n.Expr, !isNot)

	case *sqlparser.ComparisonExpr:
		if isNot {
			n = reverseComparison(n)
		}
		return n

	case *sqlparser.IsExpr:
		if isNot {
			n = reverseIs(n)
		}
		return n

	case *sqlparser.RangeCond:
		if isNot {
			return reverseRange(n)
		}
		return n
	}

	return nil
}

func reverseIs(n *sqlparser.IsExpr) *sqlparser.IsExpr {
	op := n.Operator

	switch op {
	case "is":
		n.Operator = "is null"

	case "is null":
		n.Operator = "null"

	case "is true":
		n.Operator = "is false"

	case "is not true":
		n.Operator = "is true"

	case "is false":
		n.Operator = "is true"

	case "is not false":
		n.Operator = "is false"
	}

	return n
}

func reverseRange(n *sqlparser.RangeCond) *sqlparser.OrExpr {
	l := sqlparser.ComparisonExpr{Operator: "<", Left: n.Left, Right: n.From}
	r := sqlparser.ComparisonExpr{Operator: ">", Left: n.Left, Right: n.To}
	return &sqlparser.OrExpr{Left: &l, Right: &r}
}

func reverseComparison(n *sqlparser.ComparisonExpr) *sqlparser.ComparisonExpr {
	op := n.Operator
	switch op {
	case ">":
		n.Operator = "<="

	case ">=":
		n.Operator = "<"

	case "<":
		n.Operator = ">="

	case "<=":
		n.Operator = ">"

	case "=":
		n.Operator = "<>"

	case "<>":
		n.Operator = "="

	case "!=":
		n.Operator = "="

	case "<=>":
		n.Operator = "<>"

	case "in":
		n.Operator = "not in"

	case "not in":
		n.Operator = "in"

	case "like":
		n.Operator = "not like"

	case "not like":
		n.Operator = "like"
	}
	return n
}
