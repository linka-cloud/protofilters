package protofilters



func New(filters ...*FieldFilter) *FieldsFilter {
	out := make(map[string]*Filter)
	for _, v := range filters {
		if v == nil {
			continue
		}
		out[v.Field] = v.Filter
	}
	return &FieldsFilter{Filters: out}
}

func StringEquals(s string) *Filter {
	return newStringFilter(&StringFilter{
		Condition: &StringFilter_Equals{
			Equals: s,
		},
	})
}

func StringNotEquals(s string) *Filter {
	return newStringFilter(&StringFilter{
		Condition: &StringFilter_Equals{
			Equals: s,
		},
		Not: true,
	})
}

func StringNotIEquals(s string) *Filter {
	return newStringFilter(&StringFilter{
		Condition: &StringFilter_Equals{
			Equals: s,
		},
		CaseInsensitive: true,
		Not:             true,
	})
}

func StringIEquals(s string) *Filter {
	return newStringFilter(&StringFilter{
		Condition: &StringFilter_Equals{
			Equals: s,
		},
		CaseInsensitive: true,
	})
}

func StringRegex(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_Regex{
				Regex: s,
			},
		},
	)
}

func StringNotRegex(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_Regex{
				Regex: s,
			},
			Not: true,
		},
	)
}

func newStringFilter(f *StringFilter) *Filter {
	return &Filter{
		Match: &Filter_String_{
			String_: f,
		},
	}
}

func StringIN(s ...string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_In_{
				In: &StringFilter_In{
					Values: s,
				},
			},
		})
}

func StringNotIN(s ...string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_In_{
				In: &StringFilter_In{
					Values: s,
				},
			},
			Not: true,
		})
}

func newNumberFilter(f *NumberFilter) *Filter {
	return &Filter{
		Match: &Filter_Number{
			Number: f,
		},
	}
}

func NumberEquals(n float64) *Filter {
	return newNumberFilter(
		&NumberFilter{
			Condition: &NumberFilter_Equals{
				Equals: n,
			},
		},
	)
}

func NumberNotEquals(n float64) *Filter {
	return newNumberFilter(
		&NumberFilter{
			Condition: &NumberFilter_Equals{
				Equals: n,
			},
			Not: true,
		},
	)
}

func NumberIN(n ...float64) *Filter {
	return newNumberFilter(
		&NumberFilter{
			Condition: &NumberFilter_In_{
				In: &NumberFilter_In{
					Values: n,
				},
			},
		},
	)
}

func NumberNotIN(n ...float64) *Filter {
	return newNumberFilter(
		&NumberFilter{
			Condition: &NumberFilter_In_{
				In: &NumberFilter_In{
					Values: n,
				},
			},
			Not: true,
		},
	)
}


func True() *Filter {
	return NewBoolFilter(&BoolFilter{Equals: true})
}

func False() *Filter {
	return NewBoolFilter(&BoolFilter{Equals: false})
}

func NewBoolFilter(f *BoolFilter) *Filter {
	return &Filter{
		Match: &Filter_Bool{
			Bool: f,
		},
	}
}

func Null() *Filter {
	return NewNullFilter(&NullFilter{Not: false})
}

func NotNull() *Filter {
	return NewNullFilter(&NullFilter{Not: true})
}

func NewNullFilter(f *NullFilter) *Filter {
	return &Filter{
		Match: &Filter_Null{
			Null: f,
		},
	}
}
