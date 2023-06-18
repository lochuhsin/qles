package qles

import (
	"errors"

	"github.com/xwb1989/sqlparser"
)

func mergeAnd(left, right []Object) []Object {
	if left != nil && right != nil {
		newArr := []Object{}
		for _, l := range left {
			for _, r := range right {
				newArr = append(newArr, mergeObj(l, r))
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

func mergeOr(left, right []Object) []Object {
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

func mergeObj(obj1 Object, obj2 Object) Object {
	for k, v := range obj2.PathConditions {
		obj1.PathConditions[k] = append(obj1.PathConditions[k], v...)
	}
	return obj1
}

func buildESRelation(node sqlparser.Expr, pathMap map[string]string) ([]Object, error) {
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
			obj := GetObject(pathMap)
			obj.AddCondition(n, path)
			return []Object{obj}, nil

		case "in":
			valTuple := n.Right.(sqlparser.ValTuple)
			field := n.Left.(*sqlparser.ColName).Name.String()
			path := pathMap[field]
			objs := []Object{}
			for _, expr := range valTuple {
				comp := sqlparser.ComparisonExpr{Operator: "=", Left: n.Left, Right: expr}
				obj := GetObject(pathMap)
				obj.AddCondition(&comp, path)
				objs = append(objs, obj)
			}
			return objs, nil
		case "not in":
			// exprs := n.Right
			valTuple := n.Right.(sqlparser.ValTuple)
			field := n.Left.(*sqlparser.ColName).Name.String()
			path := pathMap[field]
			obj := GetObject(pathMap)
			for _, expr := range valTuple {
				comp := sqlparser.ComparisonExpr{Operator: "!=", Left: n.Left, Right: expr}
				obj.AddCondition(&comp, path)
			}
			return []Object{obj}, nil
		}

	// name IS NULL
	case *sqlparser.IsExpr:
		column := n.Expr.(*sqlparser.ColName).Name.String()
		path := pathMap[column]
		obj := GetObject(pathMap)
		obj.AddCondition(n, path)
		return []Object{obj}, nil

	// Between
	case *sqlparser.RangeCond:
		column := n.Left.(*sqlparser.ColName).Name.String()
		path := pathMap[column]
		obj := GetObject(pathMap)
		obj.AddCondition(n, path)
		return []Object{obj}, nil

	default:
		return nil, errors.New("un-supported sql operator")
	}

	return nil, nil
}

type Object struct {
	PathConditions map[string][]sqlparser.Expr
}

func (obj *Object) AddCondition(expr sqlparser.Expr, path string) {
	obj.PathConditions[path] = append(obj.PathConditions[path], expr)
}

func (obj *Object) ToQuery() Query {
	outerCond := []Query{}
	for path, tokens := range obj.PathConditions {
		if len(tokens) == 0 {
			continue
		}
		conditions := []Query{}
		for _, t := range tokens {
			conditions = append(conditions, toESQuery(t))
		}

		var query Query
		if len(conditions) == 1 {
			query = conditions[0]
		} else {
			query = GetBoolQuery(conditions, Filter)
		}

		if path != "" {
			query = GetNestedQuery(path, query)
		}
		outerCond = append(outerCond, query)
	}
	if len(outerCond) == 1 {
		return outerCond[0]
	}
	return GetBoolQuery(outerCond, Filter)
}

func GetObject(pathMap map[string]string) Object {
	m := map[string][]sqlparser.Expr{}
	for _, v := range pathMap {
		m[v] = []sqlparser.Expr{}
	}
	return Object{PathConditions: m}
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
