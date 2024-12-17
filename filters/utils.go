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
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

// StringHasPrefix constructs a string match prefix filter
func StringHasPrefix(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_HasPrefix{
				HasPrefix: s,
			},
		},
	)
}

// StringNotHasPrefix constructs a string not match prefix filter
func StringNotHasPrefix(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_HasPrefix{
				HasPrefix: s,
			},
		},
		true,
	)
}

// StringIHasPrefix constructs a case insensitive string match prefix filter
func StringIHasPrefix(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_HasPrefix{
				HasPrefix: s,
			},
			CaseInsensitive: true,
		},
	)
}

// StringNotIHasPrefix constructs a case insensitive string not match prefix filter
func StringNotIHasPrefix(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_HasPrefix{
				HasPrefix: s,
			},
			CaseInsensitive: true,
		},
		true,
	)
}

// StringHasSuffix constructs a string match suffix filter
func StringHasSuffix(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_HasSuffix{
				HasSuffix: s,
			},
		},
	)
}

// StringNotHasSuffix constructs a string not match suffix filter
func StringNotHasSuffix(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_HasSuffix{
				HasSuffix: s,
			},
		},
		true,
	)
}

// StringIHasSuffix constructs a case insensitive string match suffix filter
func StringIHasSuffix(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_HasSuffix{
				HasSuffix: s,
			},
			CaseInsensitive: true,
		},
	)
}

// StringNotIHasSuffix constructs a case insensitive string not match suffix filter
func StringNotIHasSuffix(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_HasSuffix{
				HasSuffix: s,
			},
			CaseInsensitive: true,
		},
		true,
	)
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

func StringInf(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_Inf{
				Inf: s,
			},
		})
}

func StringSup(s string) *Filter {
	return newStringFilter(
		&StringFilter{
			Condition: &StringFilter_Sup{
				Sup: s,
			},
		})
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

// NumberNotEquals constructs a number not equals filter
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
