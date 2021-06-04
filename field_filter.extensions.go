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
			Not: true,
		},
	)
}

func True() *Filter {
	return newBoolFilter(&BoolFilter{Equals: true})
}

func False() *Filter {
	return newBoolFilter(&BoolFilter{Equals: false})
}

func newBoolFilter(f *BoolFilter) *Filter {
	return &Filter{
		Match: &Filter_Bool{
			Bool: f,
		},
	}
}

func Null() *Filter {
	return newNullFilter(&NullFilter{Not: false})
}

func NotNull() *Filter {
	return newNullFilter(&NullFilter{Not: true})
}

func newNullFilter(f *NullFilter) *Filter {
	return &Filter{
		Match: &Filter_Null{
			Null: f,
		},
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
			Not: true,
		},
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

func newDurationFilter(f *DurationFilter) *Filter {
	return &Filter{
		Match: &Filter_Duration{
			Duration: f,
		},
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
			Not: true,
		},
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

func newTimeFilter(f *TimeFilter) *Filter {
	return &Filter{
		Match: &Filter_Time{
			Time: f,
		},
	}
}
