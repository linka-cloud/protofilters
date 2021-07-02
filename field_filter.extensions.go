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

package protofilters

import (
	"regexp"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Filters map[string]*Filter

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
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_Equals{
				Equals: s,
			},
		},
		true,
	)
}

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
		},
		true,
	)
}

func NumberInf(n float64) *Filter {
	return newNumberFilter(
		&NumberFilter{
			Condition: &NumberFilter_Inf{
				Inf: n,
			},
		},
	)
}

func NumberSup(n float64) *Filter {
	return newNumberFilter(
		&NumberFilter{
			Condition: &NumberFilter_Sup{
				Sup: n,
			},
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
		},
		true,
	)
}

func True() *Filter {
	return newBoolFilter(&BoolFilter{Equals: true})
}

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

func Null() *Filter {
	return newNullFilter(&NullFilter{})
}

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

func DurationEquals(d time.Duration) *Filter {
	return newDurationFilter(
		&DurationFilter{
			Condition: &DurationFilter_Equals{
				Equals: durationpb.New(d),
			},
		},
	)
}

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

func DurationSup(d time.Duration) *Filter {
	return newDurationFilter(
		&DurationFilter{
			Condition: &DurationFilter_Sup{
				Sup: durationpb.New(d),
			},
		},
	)
}

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

func TimeEquals(t time.Time) *Filter {
	return newTimeFilter(
		&TimeFilter{
			Condition: &TimeFilter_Equals{
				Equals: timestamppb.New(t),
			},
		},
	)
}

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

func TimeAfter(t time.Time) *Filter {
	return newTimeFilter(
		&TimeFilter{
			Condition: &TimeFilter_After{
				After: timestamppb.New(t),
			},
		},
	)
}

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

func (x *BoolFilter) Match(v *bool) (bool, error) {
	if v == nil {
		return false, nil
	}
	return *v == x.GetEquals(), nil
}

func (x *NullFilter) Match(v interface{}) (bool, error) {
	return v == nil, nil
}

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
