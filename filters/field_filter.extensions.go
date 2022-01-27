/*
 Copyright 2021 Linka Cloud  All rights reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package filters

import (
	"regexp"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Filters is a map containing field path associated to a Filter
type Filters map[string]*Filter

// New constructs a *FieldsField based on the provided *FieldFilter slice
func New(filters ...*FieldFilter) *FieldsFilter {
	out := make(map[string]*Filter)
	for _, v := range filters {
		if v == nil {
			continue
		}
		// TODO(adphi):  what if there is more than one filter for this field ?
		out[v.Field] = v.Filter
	}
	return &FieldsFilter{Filters: out}
}

func Where(field string, filter *Filter) *FieldFilter {
	return &FieldFilter{Field: field, Filter: filter}
}

// Field joins the parts as un field path, e.g. Field("message", "string_field") returns "message.string_field"
func Field(parts ...string) string {
	return strings.Join(parts, ".")
}

// StringEquals constructs a string equals filter
func StringEquals(s string) *Filter {
	return newStringFilter(&StringFilter{
		Condition: &StringFilter_Equals{
			Equals: s,
		},
	})
}

// StringNotEquals constructs a string not equals filter
func StringNotEquals(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_Equals{
				Equals: s,
			},
		},
		true,
	)
}

// StringNotIEquals constructs a case insensitive string not equals filter
func StringNotIEquals(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_Equals{
				Equals: s,
			},
			CaseInsensitive: true,
		},
		true,
	)
}

// StringIEquals constructs a case insensitive string equals filter
func StringIEquals(s string) *Filter {
	return newStringFilter(&StringFilter{
		Condition: &StringFilter_Equals{
			Equals: s,
		},
		CaseInsensitive: true,
	})
}

// StringRegex constructs a string match regex filter
func StringRegex(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_Regex{
				Regex: s,
			},
		},
	)
}

// StringNotRegex constructs a string not match regex filter
func StringNotRegex(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_Regex{
				Regex: s,
			},
		},
		true,
	)
}

func newStringFilter(f *StringFilter, not ...bool) *Filter {
	return &Filter{
		Match: &Filter_String_{
			String_: f,
		},
		Not: len(not) > 0 && not[0],
	}
}

// StringIN constructs a string in slice filter
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

// StringNotIN constructs a string not in slice filter
func StringNotIN(s ...string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_In_{
				In: &StringFilter_In{
					Values: s,
				},
			},
		},
		true)
}

func newNumberFilter(f *NumberFilter, not ...bool) *Filter {
	return &Filter{
		Match: &Filter_Number{
			Number: f,
		},
		Not: len(not) > 0 && not[0],
	}
}

// NumberEquals constructs a number equals filter
func NumberEquals(n float64) *Filter {
	return newNumberFilter(
		&NumberFilter{
			Condition: &NumberFilter_Equals{
				Equals: n,
			},
		},
	)
}

// NumberInf constructs a number inferior filter
func NumberInf(n float64) *Filter {
	return newNumberFilter(
		&NumberFilter{
			Condition: &NumberFilter_Inf{
				Inf: n,
			},
		},
	)
}

// NumberSup constructs a number superior filter
func NumberSup(n float64) *Filter {
	return newNumberFilter(
		&NumberFilter{
			Condition: &NumberFilter_Sup{
				Sup: n,
			},
		},
	)
}

// NumberIN constructs a number in slice filter
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

// NumberNotIN constructs a number not in slice filter
func NumberNotIN(n ...float64) *Filter {
	return newNumberFilter(
		&NumberFilter{
			Condition: &NumberFilter_In_{
				In: &NumberFilter_In{
					Values: n,
				},
			},
		},
		true,
	)
}

// True constructs a bool is true filter
func True() *Filter {
	return newBoolFilter(&BoolFilter{Equals: true})
}

// False constructs a bool is false filter
func False() *Filter {
	return newBoolFilter(&BoolFilter{Equals: false})
}

func newBoolFilter(f *BoolFilter, not ...bool) *Filter {
	return &Filter{
		Match: &Filter_Bool{
			Bool: f,
		},
		Not: len(not) > 0 && not[0],
	}
}

// Null constructs a null check filter
func Null() *Filter {
	return newNullFilter(&NullFilter{})
}

// NotNull constructs a not null check filter
func NotNull() *Filter {
	return newNullFilter(&NullFilter{}, true)
}

func newNullFilter(f *NullFilter, not ...bool) *Filter {
	return &Filter{
		Match: &Filter_Null{
			Null: f,
		},
		Not: len(not) > 0 && not[0],
	}
}

// DurationEquals constructs a duration equals filter
func DurationEquals(d time.Duration) *Filter {
	return newDurationFilter(
		&DurationFilter{
			Condition: &DurationFilter_Equals{
				Equals: durationpb.New(d),
			},
		},
	)
}

// DurationNotEquals constructs a duration not equals filter
func DurationNotEquals(d time.Duration) *Filter {
	return newDurationFilter(
		&DurationFilter{
			Condition: &DurationFilter_Equals{
				Equals: durationpb.New(d),
			},
		},
		true,
	)
}

// DurationSup constructs a duration superior filter
func DurationSup(d time.Duration) *Filter {
	return newDurationFilter(
		&DurationFilter{
			Condition: &DurationFilter_Sup{
				Sup: durationpb.New(d),
			},
		},
	)
}

// DurationInf constructs a duration inferior filter
func DurationInf(d time.Duration) *Filter {
	return newDurationFilter(
		&DurationFilter{
			Condition: &DurationFilter_Inf{
				Inf: durationpb.New(d),
			},
		},
	)
}

func newDurationFilter(f *DurationFilter, not ...bool) *Filter {
	return &Filter{
		Match: &Filter_Duration{
			Duration: f,
		},
		Not: len(not) > 0 && not[0],
	}
}

// TimeEquals constructs a time equals filter
func TimeEquals(t time.Time) *Filter {
	return newTimeFilter(
		&TimeFilter{
			Condition: &TimeFilter_Equals{
				Equals: timestamppb.New(t),
			},
		},
	)
}

// TimeNotEquals constructs a time not equals filter
func TimeNotEquals(t time.Time) *Filter {
	return newTimeFilter(
		&TimeFilter{
			Condition: &TimeFilter_Equals{
				Equals: timestamppb.New(t),
			},
		},
		true,
	)
}

// TimeAfter constructs a time after filter
func TimeAfter(t time.Time) *Filter {
	return newTimeFilter(
		&TimeFilter{
			Condition: &TimeFilter_After{
				After: timestamppb.New(t),
			},
		},
	)
}

// TimeBefore constructs a time before filter
func TimeBefore(t time.Time) *Filter {
	return newTimeFilter(
		&TimeFilter{
			Condition: &TimeFilter_Before{
				Before: timestamppb.New(t),
			},
		},
	)
}

func newTimeFilter(f *TimeFilter, not ...bool) *Filter {
	return &Filter{
		Match: &Filter_Time{
			Time: f,
		},
		Not: len(not) > 0 && not[0],
	}
}

func (x *FieldFilter) AndF(f *FieldFilter) *Expression {
	return &Expression{Condition: x, AndExprs: []*Expression{{Condition: f}}}
}

func (x *FieldFilter) OrF(f *FieldFilter) *Expression {
	return &Expression{Condition: x, OrExprs: []*Expression{{Condition: f}}}
}

func (x *FieldFilter) And(fd string, ft *Filter) *Expression {
	return &Expression{Condition: x, AndExprs: []*Expression{{Condition: Where(fd, ft)}}}
}

func (x *FieldFilter) Or(fd string, ft *Filter) *Expression {
	return &Expression{Condition: x, OrExprs: []*Expression{{Condition: Where(fd, ft)}}}
}

func (x *Expression) And(fd string, ft *Filter) *Expression {
	x.AndExprs = append(x.AndExprs, &Expression{Condition: Where(fd, ft)})
	return x
}

func (x *Expression) Or(fd string, ft *Filter) *Expression {
	x.AndExprs = append(x.OrExprs, &Expression{Condition: Where(fd, ft)})
	return x
}

func (x *Expression) AndF(f *FieldFilter) *Expression {
	x.AndExprs = append(x.AndExprs, &Expression{Condition: f})
	return x
}

func (x *Expression) OrF(f *FieldFilter) *Expression {
	x.AndExprs = append(x.OrExprs, &Expression{Condition: f})
	return x
}

func (x *Expression) AndE(e *Expression) *Expression {
	x.AndExprs = append(x.AndExprs, e)
	return x
}

func (x *Expression) OrE(e *Expression) *Expression {
	x.OrExprs = append(x.OrExprs, e)
	return x
}

// Match applies the filter against the provided string pointer
func (x *StringFilter) Match(v *string) (bool, error) {
	if v == nil {
		return false, nil
	}
	value := *v
	insensitive := x.GetCaseInsensitive()
	switch x.GetCondition().(type) {
	case *StringFilter_Equals:
		if insensitive {
			return strings.ToLower(x.GetEquals()) == strings.ToLower(value), nil
		}
		return value == x.GetEquals(), nil
	case *StringFilter_Regex:
		reg, err := regexp.Compile(x.GetRegex())
		if err != nil {
			return false, err
		}
		return reg.MatchString(value), nil
	case *StringFilter_In_:
		for _, v := range x.GetIn().GetValues() {
			if (insensitive && strings.ToLower(v) == strings.ToLower(value)) || v == value {
				return true, nil
			}
		}
	}
	return false, nil
}

// Match applies the filter against the provided number pointer
func (x *NumberFilter) Match(v *float64) (bool, error) {
	if v == nil {
		return false, nil
	}
	val := *v
	switch x.GetCondition().(type) {
	case *NumberFilter_Equals:
		return val == x.GetEquals(), nil
	case *NumberFilter_Inf:
		return val < x.GetInf(), nil
	case *NumberFilter_Sup:
		return val > x.GetSup(), nil
	case *NumberFilter_In_:
		for _, v := range x.GetIn().GetValues() {
			if val == v {
				return true, nil
			}
		}
	}
	return false, nil
}

// Match applies the filter against the provided bool pointer
func (x *BoolFilter) Match(v *bool) (bool, error) {
	if v == nil {
		return false, nil
	}
	return *v == x.GetEquals(), nil
}

// Match applies the filter against the provided message
func (x *NullFilter) Match(v interface{}) (bool, error) {
	return v == nil, nil
}

// Match applies the filter against the provided Timestamp pointer
func (x *TimeFilter) Match(v *timestamppb.Timestamp) (bool, error) {
	if v == nil {
		return false, nil
	}
	t1 := v.AsTime()
	switch x.GetCondition().(type) {
	case *TimeFilter_Equals:
		return t1.Equal(x.GetEquals().AsTime().UTC()), nil
	case *TimeFilter_Before:
		return t1.Before(x.GetBefore().AsTime().UTC()), nil
	case *TimeFilter_After:
		return t1.After(x.GetAfter().AsTime().UTC()), nil
	}
	return false, nil
}

// Match applies the filter against the provided Duration pointer
func (x *DurationFilter) Match(v *durationpb.Duration) (bool, error) {
	if v == nil {
		return false, nil
	}
	d1 := v.AsDuration()
	switch x.GetCondition().(type) {
	case *DurationFilter_Equals:
		return d1 == x.GetEquals().AsDuration(), nil
	case *DurationFilter_Inf:
		return d1 < x.GetInf().AsDuration(), nil
	case *DurationFilter_Sup:
		return d1 > x.GetSup().AsDuration(), nil
	}
	return false, nil
}
