package qles

import (
	"errors"

	"github.com/xwb1989/sqlparser"
)

func buildSort(sqlparser.OrderBy) {

}

func convertWhereToES(whereAST *sqlparser.Where, pathMap map[string]string) Query {
	esQueryObj, err := buildESRelation(whereAST.Expr, pathMap)
	if err != nil {
		panic(err)
	}

	shouldCond := []Query{}
	for _, esObj := range esQueryObj {
		shouldCond = append(shouldCond, esObj.ToQuery())
	}

	boolQ := GetBoolQuery(shouldCond, Should)
	return boolQ
}

func mergeAnd(left, right []AndObject) []AndObject {
	if left != nil && right != nil {
		newArr := make([]AndObject, len(left)*len(right))
		i := 0
		for _, l := range left {
			for _, r := range right {
				newArr[i] = mergeObj(l, r)
			}
		}
		return newArr
	}

	if left != nil {
		return left
	}
	if right != nil {
		return right
	}
	return nil
}

func mergeOr(left, right []AndObject) []AndObject {
	if left != nil && right != nil {
		return append(left, right...)
	}
	if left != nil {
		return left
	}
	if right != nil {
		return right
	}
	return nil
}

func mergeObj(obj1 AndObject, obj2 AndObject) AndObject {
	for k, v := range obj2.PathConditions {
		obj1.PathConditions[k] = append(obj1.PathConditions[k], v...)
	}
	return obj1
}

func buildESRelation(node sqlparser.Expr, pathMap map[string]string) ([]AndObject, error) {
	switch n := node.(type) {

	case *sqlparser.AndExpr:
		left, err := buildESRelation(n.Left, pathMap)
		if err != nil {
			return nil, err
		}
		right, err := buildESRelation(n.Right, pathMap)
		if err != nil {
			return nil, err
		}
		return mergeAnd(left, right), nil

	case *sqlparser.OrExpr:
		left, err := buildESRelation(n.Left, pathMap)
		if err != nil {
			return nil, err
		}
		right, err := buildESRelation(n.Right, pathMap)
		if err != nil {
			return nil, err
		}
		return mergeOr(left, right), nil

	case *sqlparser.ParenExpr:
		return buildESRelation(n.Expr, pathMap)

	case *sqlparser.NotExpr:
		return nil, errors.New("shouldn't contain any single not operator in tree")

	case *sqlparser.ComparisonExpr:
		op := n.Operator
		switch op {
		case "=", "!=", "<>", ">=", ">", "<=", "<", "like", "not like":
			column := n.Left.(*sqlparser.ColName).Name.String()
			path := pathMap[column]
			obj := GetAndObject(pathMap)
			obj.AddCondition(n, path)
			return []AndObject{obj}, nil

		case "in":
			valTuple := n.Right.(sqlparser.ValTuple)
			field := n.Left.(*sqlparser.ColName).Name.String()
			path := pathMap[field]
			objs := []AndObject{}
			for _, expr := range valTuple {
				comp := sqlparser.ComparisonExpr{Operator: "=", Left: n.Left, Right: expr}
				obj := GetAndObject(pathMap)
				obj.AddCondition(&comp, path)
				objs = append(objs, obj)
			}
			return objs, nil
		case "not in":
			// exprs := n.Right
			valTuple := n.Right.(sqlparser.ValTuple)
			field := n.Left.(*sqlparser.ColName).Name.String()
			path := pathMap[field]
			obj := GetAndObject(pathMap)
			for _, expr := range valTuple {
				comp := sqlparser.ComparisonExpr{Operator: "!=", Left: n.Left, Right: expr}
				obj.AddCondition(&comp, path)
			}
			return []AndObject{obj}, nil
		}

	// name IS NULL
	case *sqlparser.IsExpr:
		column := n.Expr.(*sqlparser.ColName).Name.String()
		path := pathMap[column]
		obj := GetAndObject(pathMap)
		obj.AddCondition(n, path)
		return []AndObject{obj}, nil

	// Between
	case *sqlparser.RangeCond:
		column := n.Left.(*sqlparser.ColName).Name.String()
		path := pathMap[column]
		obj := GetAndObject(pathMap)
		obj.AddCondition(n, path)
		return []AndObject{obj}, nil

	default:
		return nil, errors.New("un-supported sql operator")
	}

	return nil, nil
}

type AndObject struct {
	PathConditions map[string][]sqlparser.Expr
}

func (obj *AndObject) AddCondition(expr sqlparser.Expr, path string) {
	obj.PathConditions[path] = append(obj.PathConditions[path], expr)
}

func (obj *AndObject) ToQuery() Query {
	outerCond := []Query{}
	for path, tokens := range obj.PathConditions {
		if len(tokens) == 0 {
			continue
		}
		conditions := []Query{}
		for _, t := range tokens {
			conditions = append(conditions, toESQuery(t))
		}
		var query Query = GetBoolQuery(conditions, Filter)
		if path != "" {
			query = GetNestedQuery(path, query)
		}
		outerCond = append(outerCond, query)
	}
	return GetBoolQuery(outerCond, Filter)
}

func GetAndObject(pathMap map[string]string) AndObject {
	m := map[string][]sqlparser.Expr{}
	for _, v := range pathMap {
		m[v] = []sqlparser.Expr{}
	}
	return AndObject{PathConditions: m}
}

func toESQuery(token sqlparser.Expr) Query {

	switch node := token.(type) {

	case *sqlparser.AndExpr:
		left := toESQuery(node.Left)
		right := toESQuery(node.Right)
		return ConvertAndExpr(left, right)

	case *sqlparser.OrExpr:
		left := toESQuery(node.Left)
		right := toESQuery(node.Right)
		return ConvertOrExpr(left, right)

	case *sqlparser.ParenExpr:
		return toESQuery(node.Expr)

	case *sqlparser.ComparisonExpr:
		return ConvertComparisonExpr(*node)

	case *sqlparser.RangeCond:
		return ConvertRangeExpr(*node)
	}

	return nil
}
