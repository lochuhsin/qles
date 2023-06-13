package qles

import (
	"fmt"

	"github.com/xwb1989/sqlparser"
)

func BuildSQL(sqlQuery string) (sqlparser.SQLNode, error) {
	parsed, err := sqlparser.Parse(sqlQuery)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func ReverseNot(node sqlparser.SQLNode) sqlparser.SQLNode {
	selectStatement := node.(*sqlparser.Select)
	whereNode := selectStatement.Where

	output := reverseTree(whereNode.Expr, false)
	whereNode.Expr = output
	selectStatement.Where = whereNode
	return selectStatement
}

/*Current only support single table query and expect no alias*/
func BuildESQuery(ast sqlparser.SQLNode, pathMap map[string]string) Query {
	selectAST := ast.(*sqlparser.Select)
	_ = selectAST.From
	_ = selectAST.OrderBy

	orderby := selectAST.OrderBy

	for _, o := range orderby {
		fmt.Println(o.Direction, o.Expr.(*sqlparser.ColName).Name.String())
	}
	// order by
	mainQ := convertWhereToES(selectAST.Where, pathMap)
	return mainQ
}
