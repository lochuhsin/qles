package qles

import (
	"github.com/xwb1989/sqlparser"
)

/*Current only support single table query and expect no alias*/
func BuildES(sqlQuery string, pathMap map[string]string) (Query, error) {

	ast, err := preprocess(sqlQuery)
	if err != nil {
		return nil, err
	}
	selectAST := ast.(*sqlparser.Select)
	// search
	searchQ := ConvertWhereToES(selectAST.Where, pathMap)
	// order by
	sortQ := ConvertOrderByToES(selectAST.Where, selectAST.OrderBy, pathMap)
	// fields
	fields := ConvertSelectColumnsToES(selectAST.SelectExprs)
	return GetESQuery(searchQ, sortQ, fields), nil
}

func preprocess(sqlQuery string) (sqlparser.SQLNode, error) {
	parsed, err := sqlparser.Parse(sqlQuery)
	if err != nil {
		return nil, err
	}
	selectStatement := parsed.(*sqlparser.Select)
	whereNode := selectStatement.Where
	if whereNode != nil {

		whereNode.Expr = reverseTree(whereNode.Expr, false)
	}
	selectStatement.Where = whereNode
	return selectStatement, nil
}

func ConvertSelectColumnsToES(selectExprs sqlparser.SelectExprs) []string {
	columns := []string{}
	for _, c := range selectExprs {
		str := sqlparser.String(c)
		if str == "*" {
			return nil
		}
		columns = append(columns, str)
	}
	return columns
}

func ConvertOrderByToES(whereAST *sqlparser.Where, orderbyAST sqlparser.OrderBy, pathMap map[string]string) []map[string]SortObject {
	orderFields := map[string]string{}
	for _, fieldObj := range orderbyAST {
		field := fieldObj.Expr.(*sqlparser.ColName).Name.String()
		orderFields[field] = fieldObj.Direction
	}

	components := []map[string]SortObject{}
	for field, order := range orderFields {
		// normal field
		if path, ok := pathMap[field]; ok && path == "" {
			sortObj := GetSortObj(field, "avg", order, nil)
			components = append(components, sortObj)
		} else if ok && path != "" {
			newPathMap := map[string]string{}
			for k, v := range pathMap {
				if v == path {
					newPathMap[k] = ""
				} else if v != "" {
					newPathMap[k] = v
				}
			}
			query := ConvertWhereToES(whereAST, newPathMap)
			var sortObj map[string]SortObject
			if query == nil {
				sortObj = GetSortObj(field, "avg", order, nil)
			} else {
				nestedSortQ := GetNestedSortQuery(path, query)
				sortObj = GetSortObj(field, "avg", order, &nestedSortQ)
			}
			components = append(components, sortObj)
		}
	}
	return components
}

func ConvertWhereToES(whereAST *sqlparser.Where, pathMap map[string]string) Query {
	if whereAST == nil {
		return nil
	}
	esQueryObj, err := buildESRelation(whereAST.Expr, pathMap)
	if err != nil {
		panic(err)
	}
	if len(esQueryObj) == 1 {
		return esQueryObj[0].ToQuery()
	}

	shouldCond := []Query{}
	for _, esObj := range esQueryObj {
		shouldCond = append(shouldCond, esObj.ToQuery())
	}
	return GetBoolQuery(shouldCond, Should)
}
