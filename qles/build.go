package qles

import (
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
	if whereNode != nil {

		whereNode.Expr = reverseTree(whereNode.Expr, false)
	}
	selectStatement.Where = whereNode
	return selectStatement
}

/*Current only support single table query and expect no alias*/
func BuildESQuery(ast sqlparser.SQLNode, pathMap map[string]string) Query {
	selectAST := ast.(*sqlparser.Select)
	// order by
	searchQ := ConvertWhereToES(selectAST.Where, pathMap)
	sortQ := ConvertOrderByToES(selectAST.Where, selectAST.OrderBy, pathMap)
	return GetESQuery(searchQ, sortQ)
}
