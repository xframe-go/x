package requests

type QueryParams struct {
	Keyword  string            `json:"keyword"`
	Filters  Filters           `json:"filters"`
	Sorter   map[string]string `json:"sorter"`
	Preload  []string          `json:"preload"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

type Filters []Filter

type Filter struct {
	Field    string      `json:"field"`
	Operator Operator    `json:"operator"`
	Value    FilterValue `json:"value"`
}

type Operator string

const (
	Equal     Operator = "eq"
	NotEqual  Operator = "ne"
	Greater   Operator = "gt"
	Less      Operator = "lt"
	GreaterEq Operator = "gte"
	LessEq    Operator = "lte"
	In        Operator = "in"
	NotIn     Operator = "nin"
	Contains  Operator = "like"
	Between   Operator = "btw"
)
