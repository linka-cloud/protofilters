package filters

var ExpressionFields = struct {
	Condition string
	AndExprs  string
	OrExprs   string
}{
	Condition: "condition",
	AndExprs:  "and_exprs",
	OrExprs:   "or_exprs",
}

var FieldsFilterFields = struct {
	Filters string
}{
	Filters: "filters",
}

var FieldFilterFields = struct {
	Field  string
	Filter string
}{
	Field:  "field",
	Filter: "filter",
}

var FilterFields = struct {
	String_  string
	Number   string
	Bool     string
	Null     string
	Time     string
	Duration string
	Not      string
}{
	String_:  "string",
	Number:   "number",
	Bool:     "bool",
	Null:     "null",
	Time:     "time",
	Duration: "duration",
	Not:      "not",
}

var StringFilterFields = struct {
	Equals          string
	Regex           string
	In              string
	CaseInsensitive string
}{
	Equals:          "equals",
	Regex:           "regex",
	In:              "in",
	CaseInsensitive: "case_insensitive",
}

var NumberFilterFields = struct {
	Equals string
	Sup    string
	Inf    string
	In     string
}{
	Equals: "equals",
	Sup:    "sup",
	Inf:    "inf",
	In:     "in",
}

var NullFilterFields = struct {
}{}

var BoolFilterFields = struct {
	Equals string
}{
	Equals: "equals",
}

var TimeFilterFields = struct {
	Equals string
	Before string
	After  string
}{
	Equals: "equals",
	Before: "before",
	After:  "after",
}

var DurationFilterFields = struct {
	Equals string
	Sup    string
	Inf    string
}{
	Equals: "equals",
	Sup:    "sup",
	Inf:    "inf",
}

var StringFilter_InFields = struct {
	Values string
}{
	Values: "values",
}

var NumberFilter_InFields = struct {
	Values string
}{
	Values: "values",
}
