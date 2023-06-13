package qles

import (
	"errors"
	"strconv"

	"github.com/xwb1989/sqlparser"
)

/* a basic interface for all es query types*/
type Query interface {
	iQuery()
}

func (TermQuery) iQuery()   {}
func (RangeQuery) iQuery()  {}
func (ExistQuery) iQuery()  {}
func (BoolQuery) iQuery()   {}
func (NestedQuery) iQuery() {}
func (SortQuery) iQuery()   {}

type SortComponent struct {
}

type SortField struct {
}

type SortQuery struct {
	Sort []map[string]SortField `json:"sort"`
}

func GetSortQuery(components []map[string]SortField) SortQuery {
	return SortQuery{Sort: components}
}

type TermQuery struct {
	Term map[string]any `json:"term"`
}

func GetTermQuery(f string, v any) TermQuery {
	return TermQuery{Term: map[string]any{f: v}}
}

type rangeType int

const (
	GT rangeType = iota
	GTE
	LT
	LTE
)

type rangeComponent struct {
	GT  any `json:"gt,omitempty"`
	GTE any `json:"gte,omitempty"`
	LT  any `json:"lt,omitempty"`
	LTE any `json:"lte,omitempty"`
}

type RangeQuery struct {
	Range map[string]rangeComponent `json:"range"`
}

func (rq *RangeQuery) Add(f string, v any, op rangeType) {
	if _, ok := rq.Range[f]; !ok {
		panic("field not found in range q")
	}
	obj := rq.Range[f]
	switch op {
	case GT:
		obj.GT = v
		rq.Range[f] = obj
	case GTE:
		obj.GT = v
		rq.Range[f] = obj
	case LT:
		obj.GT = v
		rq.Range[f] = obj
	case LTE:
		obj.GT = v
		rq.Range[f] = obj
	}
}

func GetRangeQuery(f string, v any, op rangeType) RangeQuery {
	var comp rangeComponent
	switch op {
	case GT:
		comp = rangeComponent{GT: v}
	case GTE:
		comp = rangeComponent{GTE: v}
	case LT:
		comp = rangeComponent{LT: v}
	case LTE:
		comp = rangeComponent{LTE: v}
	default:
		panic("Unsupported range type operation")
	}
	return RangeQuery{Range: map[string]rangeComponent{f: comp}}
}

type ExistQuery struct {
	Exist map[string]string `json:"exists"`
}

func GetExistQuery(field string) ExistQuery {
	return ExistQuery{Exist: map[string]string{"field": field}}
}

type BoolType int

const (
	Filter BoolType = iota
	Should
	Must
	MustNot
)

type boolComponents struct {
	Filter  []Query `json:"filter,omitempty"`
	Should  []Query `json:"should,omitempty"`
	Must    []Query `json:"must,omitempty"`
	MustNot []Query `json:"must_not,omitempty"`
}

type BoolQuery struct {
	Bool boolComponents
}

func GetBoolQuery(query []Query, queryType BoolType) BoolQuery {
	var bq BoolQuery
	switch queryType {
	case Filter:
		bq = BoolQuery{Bool: boolComponents{Filter: query}}
	case Should:
		bq = BoolQuery{Bool: boolComponents{Should: query}}
	case Must:
		bq = BoolQuery{Bool: boolComponents{Must: query}}
	case MustNot:
		bq = BoolQuery{Bool: boolComponents{MustNot: query}}
	}
	return bq
}

func (bq *BoolQuery) Add(query []Query, queryType BoolType) {
	switch queryType {
	case Filter:
		if bq.Bool.Filter != nil {
			bq.Bool.Filter = append(bq.Bool.Filter, query...)
		} else {
			bq.Bool.Filter = query
		}

	case Should:
		if bq.Bool.Should != nil {
			bq.Bool.Should = append(bq.Bool.Should, query...)
		} else {
			bq.Bool.Should = query
		}
	case Must:
		if bq.Bool.Must != nil {
			bq.Bool.Must = append(bq.Bool.Must, query...)
		} else {
			bq.Bool.Must = query
		}
	case MustNot:
		if bq.Bool.MustNot != nil {
			bq.Bool.MustNot = append(bq.Bool.MustNot, query...)
		} else {
			bq.Bool.MustNot = query
		}
	}
}

type nestedComponent struct {
	Path  string `json:"path"`
	Query Query  `json:"query"`
}

type NestedQuery struct {
	Nested nestedComponent
}

func GetNestedQuery(path string, query Query) NestedQuery {
	return NestedQuery{Nested: nestedComponent{Path: path, Query: query}}
}

// TODO: add some sort of type constrain
func ConvertEQ(field string, value any) TermQuery {
	var termQ TermQuery
	switch v := value.(type) {
	case int:
		termQ = GetTermQuery(field, v)
	case float32:
		termQ = GetTermQuery(field, v)
	case float64:
		termQ = GetTermQuery(field, v)
	case string:
		termQ = GetTermQuery(field, v)
	}
	return termQ
}

func ConvertNEQ(field string, value any) BoolQuery {
	eq := ConvertEQ(field, value)
	return GetBoolQuery([]Query{eq}, MustNot)
}

func ConvertGT(field string, value any) RangeQuery {
	var RangeQ RangeQuery
	switch v := value.(type) {
	case int:
		RangeQ = GetRangeQuery(field, v, GT)
	case float32:
		RangeQ = GetRangeQuery(field, v, GT)
	case float64:
		RangeQ = GetRangeQuery(field, v, GT)
	default:
		panic("un support range type")
	}
	return RangeQ
}

func ConvertGTE(field string, value any) RangeQuery {
	var RangeQ RangeQuery
	switch v := value.(type) {
	case int:
		RangeQ = GetRangeQuery(field, v, GTE)
	case float32:
		RangeQ = GetRangeQuery(field, v, GTE)
	case float64:
		RangeQ = GetRangeQuery(field, v, GTE)
	default:
		panic("un support range type")
	}
	return RangeQ
}

func ConvertLT(field string, value any) RangeQuery {
	var RangeQ RangeQuery
	switch v := value.(type) {
	case int:
		RangeQ = GetRangeQuery(field, v, LT)
	case float32:
		RangeQ = GetRangeQuery(field, v, LT)
	case float64:
		RangeQ = GetRangeQuery(field, v, LT)
	default:
		panic("un support range type")
	}
	return RangeQ
}

func ConvertLTE(field string, value any) RangeQuery {
	var RangeQ RangeQuery
	switch v := value.(type) {
	case int:
		RangeQ = GetRangeQuery(field, v, LTE)
	case float32:
		RangeQ = GetRangeQuery(field, v, LTE)
	case float64:
		RangeQ = GetRangeQuery(field, v, LTE)
	default:
		panic("un support range type")
	}
	return RangeQ
}

func ConvertExist(field string) ExistQuery {
	return GetExistQuery(field)
}

func ConvertNotExist(field string) BoolQuery {
	exist := GetExistQuery(field)
	return GetBoolQuery([]Query{exist}, MustNot)
}

func ConvertRangeExpr(token sqlparser.RangeCond) RangeQuery {
	field := token.Left.(*sqlparser.ColName).Name.String()

	low, err := ConvertToNativeType(token.Left.(*sqlparser.SQLVal))
	if err != nil {
		panic(err)
	}
	high, err := ConvertToNativeType(token.Left.(*sqlparser.SQLVal))
	if err != nil {
		panic(err)
	}

	rq := GetRangeQuery(field, low, LTE)
	rq.Add(field, high, GTE)
	return rq
}

func ConvertComparisonExpr(token sqlparser.ComparisonExpr) Query {
	operator := token.Operator
	field := token.Left.(*sqlparser.ColName).Name.String()
	value, err := ConvertToNativeType(token.Right.(*sqlparser.SQLVal))
	if err != nil {
		panic(err)
	}
	switch operator {
	case "=":
		return ConvertEQ(field, value)
	case "!=":
		return ConvertNEQ(field, value)
	case "<>":
		return ConvertNEQ(field, value)
	case "<=>":
		return ConvertNEQ(field, value)
	case ">":
		return ConvertGT(field, value)
	case ">=":
		return ConvertGTE(field, value)
	case "<":
		return ConvertLT(field, value)
	case "<=":
		return ConvertLTE(field, value)
	case "like":
		panic("like operator not implemented yet")
	case "not like":
		panic("like operator not implemented yet")
	default:
		panic("Un recognized comparison operator")
	}
}

func ConvertAndExpr(left, right Query) BoolQuery {
	bq := GetBoolQuery([]Query{left}, Filter)
	bq.Add([]Query{right}, Filter)
	return bq
}

func ConvertOrExpr(left, right Query) BoolQuery {
	bq := GetBoolQuery([]Query{left}, Should)
	bq.Add([]Query{right}, Should)
	return bq
}

/*
const (

	StrVal = ValType(iota)
	IntVal
	FloatVal
	HexNum
	HexVal
	ValArg
	BitVal

)
*/
func ConvertToNativeType(val *sqlparser.SQLVal) (any, error) {
	switch val.Type {
	case sqlparser.StrVal:
		return string(val.Val), nil
	case sqlparser.IntVal:
		return strconv.Atoi(string(val.Val))
	case sqlparser.FloatVal:
		return strconv.ParseFloat(string(val.Val), 64)
	default:
		return "", errors.New("Other type not implemented yet")
	}
}
