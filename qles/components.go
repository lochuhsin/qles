package qles

/* a basic interface for all es query types*/
type Query interface {
	iQuery()
}

type TermLevelQuery interface {
	Query
	iTermQ()
}

func (TermQuery) iQuery() {}
func (TermQuery) iTermQ() {}

func (RangeQuery) iQuery() {}
func (RangeQuery) iTermQ() {}

func (ExistQuery) iQuery() {}
func (ExistQuery) iTermQ() {}

func (NestedQuery) iQuery() {}

type ComposeLevelQuery interface {
	Query
	iComposeQ()
}

func (BoolQuery) iQuery()    {}
func (BoolQuery) iComposeQ() {}

type TermQuery struct {
	Term map[string]any `json:"term"`
}

func GetTermQuery(f string, v any) TermQuery {
	return TermQuery{Term: map[string]any{f: v}}
}

type RangeQuery struct {
	Range map[string]map[string]any `json:"range"`
}

func GetRangeQuery(f string, v any, op string) RangeQuery {
	return RangeQuery{Range: map[string]map[string]any{f: {op: v}}}
}

type ExistQuery struct {
	Exist map[string]string `json:"exists"`
}

func GetExistQuery(field string) ExistQuery {
	return ExistQuery{Exist: map[string]string{"field": field}}
}

type BoolQuery struct {
	Bool map[string]ComposeLevelQuery `json:"bool"`
}

func GetBoolQuery() BoolQuery {
	return BoolQ
}

func (bq *BoolQuery) AddFilter(conditions []map[string]Query) {
	bq.Filter = append(bq.Filter, conditions...)
}

func (bq *BoolQuery) AddShould(conditions []map[string]Query) {
	bq.Should = append(bq.Should, conditions...)
}

func (bq *BoolQuery) AddMust(conditions []map[string]Query) {
	bq.Must = append(bq.Must, conditions...)
}

func (bq *BoolQuery) AddMustNot(conditions []map[string]Query) {
	bq.MustNot = append(bq.MustNot, conditions...)
}

type NestedQuery struct {
	Path  string           `json:"path"`
	Query map[string]Query `json:"query"`
}

func GetNestedQuery(path string, query map[string]Query) map[string]NestedQuery {
	return map[string]NestedQuery{"nested": {Path: path, Query: query}}
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

func ConvertGT(field string, value any) RangeQuery {
	var RangeQ RangeQuery
	switch v := value.(type) {
	case int:
		RangeQ = GetRangeQuery(field, v, "GT")
	case float32:
		RangeQ = GetRangeQuery(field, v, "GT")
	case float64:
		RangeQ = GetRangeQuery(field, v, "GT")
	default:
		panic("un support range type")
	}
	return RangeQ
}

func ConvertGTE(field string, value any) RangeQuery {
	var RangeQ RangeQuery
	switch v := value.(type) {
	case int:
		RangeQ = GetRangeQuery(field, v, "GTE")
	case float32:
		RangeQ = GetRangeQuery(field, v, "GTE")
	case float64:
		RangeQ = GetRangeQuery(field, v, "GTE")
	default:
		panic("un support range type")
	}
	return RangeQ
}

func ConvertLT(field string, value any) RangeQuery {
	var RangeQ RangeQuery
	switch v := value.(type) {
	case int:
		RangeQ = GetRangeQuery(field, v, "LT")
	case float32:
		RangeQ = GetRangeQuery(field, v, "LT")
	case float64:
		RangeQ = GetRangeQuery(field, v, "LT")
	default:
		panic("un support range type")
	}
	return RangeQ
}

func ConvertLTE(field string, value any) RangeQuery {
	var RangeQ RangeQuery
	switch v := value.(type) {
	case int:
		RangeQ = GetRangeQuery(field, v, "LTE")
	case float32:
		RangeQ = GetRangeQuery(field, v, "LTE")
	case float64:
		RangeQ = GetRangeQuery(field, v, "LTE")
	default:
		panic("un support range type")
	}
	return RangeQ
}

func ConvertExist(field string, value any) {}

func ConvertNotExist(field string, value any) {}
