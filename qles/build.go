package qles

import (
	"github.com/xwb1989/sqlparser"
)

func Build(sqlQuery string) (sqlparser.SQLNode, error) {
	parsed, err := sqlparser.Parse(sqlQuery)
	if err != nil {
		return nil, err
	}

	reversedAST := ReverseNot(parsed)
	return reversedAST, nil
}

func ReverseNot(node sqlparser.SQLNode) sqlparser.SQLNode {
	selectStatement := node.(*sqlparser.Select)
	whereNode := selectStatement.Where

	output := reverseTree(whereNode.Expr, false)
	whereNode.Expr = output
	selectStatement.Where = whereNode
	return selectStatement
}
