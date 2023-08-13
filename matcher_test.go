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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"go.linka.cloud/protofilters/filters"
	test "go.linka.cloud/protofilters/tests/pb"
)

func TestOptionalsNil(t *testing.T) {
	tests := []struct {
		name string
		m    proto.Message
		f    filters.FieldFilterer
		ok   bool
	}{
		{
			name: "empty optional null",
			m:    &test.Test{},
			f:    filters.Where("optional_bool_field").Null(),
			ok:   true,
		},
		{
			name: "empty optional not null",
			m:    &test.Test{},
			f:    filters.Where("optional_bool_field").NotNull(),
			ok:   false,
		},
		{
			name: "empty optional true",
			m:    &test.Test{},
			f:    filters.Where("optional_bool_field").True(),
			ok:   false,
		},
		{
			name: "empty optional false",
			m:    &test.Test{},
			f:    filters.Where("optional_bool_field").False(),
			ok:   false,
		},
		{
			name: "false optional not null",
			m:    &test.Test{OptionalBoolField: proto.Bool(false)},
			f:    filters.Where("optional_bool_field").NotNull(),
			ok:   true,
		},
		{
			name: "false optional null",
			m:    &test.Test{OptionalBoolField: proto.Bool(false)},
			f:    filters.Where("optional_bool_field").Null(),
			ok:   false,
		},
		{
			name: "false optional true",
			m:    &test.Test{OptionalBoolField: proto.Bool(false)},
			f:    filters.Where("optional_bool_field").True(),
			ok:   false,
		},
		{
			name: "false optional false",
			m:    &test.Test{OptionalBoolField: proto.Bool(false)},
			f:    filters.Where("optional_bool_field").False(),
			ok:   true,
		},

		{
			name: "empty number optional not null",
			m:    &test.Test{},
			f:    filters.Where("optional_number_field").NotNull(),
			ok:   false,
		},
		{
			name: "empty number optional null",
			m:    &test.Test{},
			f:    filters.Where("optional_number_field").Null(),
			ok:   true,
		},
		{
			name: "number optional not null",
			m:    &test.Test{OptionalNumberField: proto.Int64(0)},
			f:    filters.Where("optional_number_field").NotNull(),
			ok:   true,
		},
		{
			name: "number optional null",
			m:    &test.Test{OptionalNumberField: proto.Int64(0)},
			f:    filters.Where("optional_number_field").Null(),
			ok:   false,
		},
		{
			name: "number optional equals",
			m:    &test.Test{OptionalNumberField: proto.Int64(0)},
			f:    filters.Where("optional_number_field").NumberEquals(0),
			ok:   true,
		},
		{
			name: "number optional not equals",
			m:    &test.Test{OptionalNumberField: proto.Int64(0)},
			f:    filters.Where("optional_number_field").NumberNotEquals(0),
			ok:   false,
		},
		{
			name: "string optional not null",
			m:    &test.Test{},
			f:    filters.Where("optional_string_field").NotNull(),
			ok:   false,
		},
		{
			name: "string optional null",
			m:    &test.Test{},
			f:    filters.Where("optional_string_field").Null(),
			ok:   true,
		},
		{
			name: "empty string optional not null",
			m:    &test.Test{OptionalStringField: proto.String("")},
			f:    filters.Where("optional_string_field").NotNull(),
			ok:   true,
		},
		{
			name: "empty string optional null",
			m:    &test.Test{OptionalStringField: proto.String("")},
			f:    filters.Where("optional_string_field").Null(),
			ok:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := Match(tt.m, tt.f)
			require.NoError(t, err)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func TestFieldFilter(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{StringField: "ok"}
	ok, err := Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"noop": filters.StringEquals("ok"),
	}})
	assert.Error(err)
	assert.False(ok)
	ok, err = Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"messageField": filters.Null(),
	}})
	assert.Error(err)
	assert.False(ok)

	assert.True(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"message_field": filters.Null(),
	}}))
	m.MessageField = m
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"message_field": filters.Null(),
	}}))
	ok, err = Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"message_field.string_field.message_field": filters.Null(),
	}})
	assert.Error(err)
	assert.False(ok)
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"message_field.message_field.message_field": filters.Null(),
	}}))
	assert.True(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"message_field.message_field.message_field.string_field": filters.StringIN("ok"),
	}}))
}

func TestString(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{StringField: "ok"}
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"string_value_field": filters.StringEquals("ok"),
	}}))
	m.StringValueField = wrapperspb.String("pointer")
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"string_value_field": filters.StringEquals("ok"),
	}}))
	assert.True(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"string_value_field": filters.StringEquals("pointer"),
	}}))
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"string_value_field": nil,
	}}))
	assert.True(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"string_field": filters.StringEquals("ok"),
	}}))
	assert.True(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"string_field": filters.StringIN("other", "ok"),
	}}))
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"string_field": filters.StringIN("other", "noop"),
	}}))
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"string_field": filters.StringNotRegex(`[a-z](.+)`),
	}}))
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"string_value_field": filters.StringNotRegex(`[a-z](.+)`),
	}}))
	m.StringValueField = wrapperspb.String("whatever")
	assert.True(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"string_value_field": filters.StringRegex(`[a-z](.+)`),
	}}))
	assert.True(Match(m, filters.Where("string_value_field").StringHasPrefix("what")))
	assert.False(Match(m, filters.Where("string_value_field").StringHasPrefix("noop")))
	assert.True(Match(m, filters.Where("string_value_field").StringHasSuffix("ever")))
	assert.False(Match(m, filters.Where("string_value_field").StringHasSuffix("noop")))
}

func TestEnum(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{EnumField: test.Test_Type(42)}
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"enum_field": filters.StringIN("OTHER"),
	}}))
	assert.True(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"enum_field": filters.NumberIN(0, 42),
	}}))
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"enum_field": filters.NumberNotIN(0, 42),
	}}))
	m.EnumField = 0
	assert.True(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"string_field": filters.StringNotIN(),
		"enum_field":   filters.StringIN("NONE"),
	}}))
}

func TestNumber(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{NumberField: 42}
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"number_field": filters.NumberEquals(0),
	}}))
	assert.True(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"number_field": filters.NumberEquals(42),
	}}))
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"number_field": filters.NumberIN(0, 22),
	}}))
	assert.True(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"number_field": filters.NumberInf(43),
	}}))
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"number_field": filters.NumberInf(41),
	}}))
	assert.True(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"number_field": filters.NumberSup(41),
	}}))
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"number_field": filters.NumberSup(43),
	}}))
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"number_value_field": filters.NumberSup(41),
	}}))
	m.NumberValueField = wrapperspb.Int64(42)
	assert.False(Match(m, &filters.FieldsFilter{Filters: filters.Filters{
		"number_value_field": filters.NumberSup(43),
	}}))
}

func TestDuration(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{DurationValueField: durationpb.New(42)}
	assert.True(MatchFilters(m, &filters.FieldFilter{
		Field:  "duration_value_field",
		Filter: filters.DurationEquals(42),
	}))
	assert.True(MatchFilters(m, &filters.FieldFilter{
		Field:  "duration_value_field",
		Filter: filters.DurationInf(43),
	}))
	assert.False(MatchFilters(m, &filters.FieldFilter{
		Field:  "duration_value_field",
		Filter: filters.DurationInf(41),
	}))
	assert.False(MatchFilters(m, &filters.FieldFilter{
		Field:  "duration_value_field",
		Filter: filters.DurationSup(43),
	}))
	assert.True(MatchFilters(m, &filters.FieldFilter{
		Field:  "duration_value_field",
		Filter: filters.DurationSup(41),
	}))
}

func TestTime(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{TimeValueField: timestamppb.New(time.Now().UTC())}
	assert.True(MatchFilters(m, &filters.FieldFilter{
		Field:  "time_value_field",
		Filter: filters.TimeEquals(m.TimeValueField.AsTime()),
	}))
	assert.True(MatchFilters(m, &filters.FieldFilter{
		Field:  "time_value_field",
		Filter: filters.TimeAfter(m.TimeValueField.AsTime().Add(-time.Second)),
	}))
	assert.False(MatchFilters(m, &filters.FieldFilter{
		Field:  "time_value_field",
		Filter: filters.TimeAfter(m.TimeValueField.AsTime().Add(time.Second)),
	}))
	assert.False(MatchFilters(m, &filters.FieldFilter{
		Field:  "time_value_field",
		Filter: filters.TimeBefore(m.TimeValueField.AsTime().Add(-time.Second)),
	}))
	assert.True(MatchFilters(m, &filters.FieldFilter{
		Field:  "time_value_field",
		Filter: filters.TimeBefore(m.TimeValueField.AsTime().Add(time.Second)),
	}))
}

func TestRepeated(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{
		StringField: "whatever",
		NumberField: 42,
		BoolField:   true,
		EnumField:   test.Test_ONE,
		MessageField: &test.Test{
			StringField:         "whatever2",
			NumberField:         43,
			EnumField:           test.Test_TWO,
			RepeatedStringField: []string{"three", "four"},
		},
		RepeatedStringField: []string{"one", "two"},
		TimeValueField:      timestamppb.Now(),
		DurationValueField:  durationpb.New(5 * time.Second),
	}
	assert.False(MatchFilters(m, &filters.FieldFilter{
		Field:  "repeated_string_field",
		Filter: filters.StringIN("four", "five"),
	}))
	assert.True(MatchFilters(m, &filters.FieldFilter{
		Field:  "repeated_string_field",
		Filter: filters.StringIN("two", "three"),
	}))
	assert.False(MatchFilters(m, &filters.FieldFilter{
		Field:  "repeated_string_field",
		Filter: filters.StringNotIN("two", "three"),
	}))
}

func TestSubField(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{
		MessageField: &test.Test{
			StringField:         "zero",
			NumberField:         43,
			EnumField:           test.Test_TWO,
			RepeatedStringField: []string{"three", "four"},
			RepeatedMessageField: []*test.Test{{
				StringField: "five",
			}},
		},
	}
	assert.False(MatchFilters(m, &filters.FieldFilter{
		Field:  "message_field.string_field",
		Filter: filters.StringIN("one", "two"),
	}))
	assert.True(MatchFilters(m, &filters.FieldFilter{
		Field:  "message_field.string_field",
		Filter: filters.StringIN("zero", "one"),
	}))
	assert.False(MatchFilters(m, &filters.FieldFilter{
		Field:  "message_field.number_field",
		Filter: filters.NumberEquals(42),
	}))
	assert.True(MatchFilters(m, &filters.FieldFilter{
		Field:  "message_field.number_field",
		Filter: filters.NumberEquals(43),
	}))
	assert.True(MatchFilters(m, &filters.FieldFilter{
		Field:  "message_field.repeated_message_field.string_field",
		Filter: filters.StringEquals("five"),
	}))
	assert.False(MatchFilters(m, &filters.FieldFilter{
		Field:  "message_field.repeated_message_field.string_field",
		Filter: filters.StringEquals("six"),
	}))
}

func TestSubFieldNil(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{}
	assert.False(MatchFilters(m, &filters.FieldFilter{
		Field:  "message_field.number_field",
		Filter: filters.NumberEquals(43),
	}))
	assert.False(MatchFilters(m, &filters.FieldFilter{
		Field:  "message_field.repeated_string_field",
		Filter: filters.StringIN("three", "four"),
	}))
}

func TestOptional(t *testing.T) {
	assert := assert.New(t)
	e := test.Test_Type(42)
	m := &test.Test{
		OptionalEnumField:   &e,
		OptionalStringField: proto.String("42"),
		OptionalBoolField:   proto.Bool(false),
		OptionalNumberField: proto.Int64(42),
	}
	assert.True(MatchFilters(m, &filters.FieldFilter{
		Field:  "optional_string_field",
		Filter: filters.StringIN("42", "43"),
	}))
	assert.True(MatchFilters(m, &filters.FieldFilter{
		Field:  "optional_bool_field",
		Filter: filters.False(),
	}))
	assert.False(MatchFilters(m, &filters.FieldFilter{
		Field:  "optional_bool_field",
		Filter: filters.True(),
	}))
	assert.True(MatchFilters(m, &filters.FieldFilter{
		Field:  "optional_number_field",
		Filter: filters.NumberInf(43),
	}))
}

func TestMatchExpression(t *testing.T) {
	assert := assert.New(t)
	e := test.Test_Type(42)
	m := &test.Test{
		StringField: "whatever",
		NumberField: 42,
		BoolField:   true,
		EnumField:   test.Test_ONE,
		MessageField: &test.Test{
			StringField:         "whatever2",
			NumberField:         43,
			EnumField:           test.Test_TWO,
			RepeatedStringField: []string{"three", "four"},
		},
		RepeatedStringField: []string{"one", "two"},
		TimeValueField:      timestamppb.Now(),
		DurationValueField:  durationpb.New(5 * time.Second),
		OptionalEnumField:   &e,
		OptionalStringField: proto.String("42"),
		OptionalBoolField:   proto.Bool(false),
		OptionalNumberField: proto.Int64(42),
	}
	assert.False(Match(m, &filters.Expression{
		Condition: &filters.FieldFilter{
			Field:  "string_field",
			Filter: filters.StringEquals("not whatever"),
		},
	}))
	assert.True(Match(m, &filters.Expression{
		Condition: &filters.FieldFilter{
			Field:  "string_field",
			Filter: filters.StringEquals("whatever"),
		},
	}))
	assert.True(Match(m, &filters.Expression{
		Condition: &filters.FieldFilter{
			Field:  "string_field",
			Filter: filters.StringEquals("whatever"),
		},
		AndExprs: []*filters.Expression{{
			Condition: &filters.FieldFilter{
				Field:  "string_field",
				Filter: filters.StringNotEquals("what"),
			},
		}},
	}))
	assert.False(Match(m, &filters.Expression{
		Condition: &filters.FieldFilter{
			Field:  "string_field",
			Filter: filters.StringEquals("whatever"),
		},
		AndExprs: []*filters.Expression{{
			Condition: &filters.FieldFilter{
				Field:  "string_field",
				Filter: filters.StringNotEquals("whatever"),
			},
		}},
	}))
	f1 := filters.Where("string_field").StringEquals("whatever").
		And("number_field").NumberIN(42, 43).
		Or(filters.Where("bool_field").False()).
		Or(filters.Where("optional_bool_field").False())

	assert.True(Match(m, f1))
	f2 := &filters.Expression{
		Condition: &filters.FieldFilter{
			Field:  "string_field",
			Filter: filters.StringEquals("whatever"),
		},
		AndExprs: []*filters.Expression{{
			Condition: &filters.FieldFilter{
				Field:  "number_field",
				Filter: filters.NumberIN(42, 43),
			},
		}},
		OrExprs: []*filters.Expression{{
			Condition: &filters.FieldFilter{
				Field:  "bool_field",
				Filter: filters.False(),
			}}, {
			Condition: &filters.FieldFilter{
				Field:  "optional_bool_field",
				Filter: filters.False(),
			}},
		},
	}
	assert.Equal(f1.Expr(), f2)
	assert.True(Match(m, f2))
	assert.Equal([]string{"bool_field", "number_field", "optional_bool_field", "string_field"}, f1.Expr().Fields())
	assert.True(Match(m, &filters.Expression{
		Condition: &filters.FieldFilter{
			Field:  "string_field",
			Filter: filters.StringEquals("whatever"),
		},
		AndExprs: []*filters.Expression{{
			Condition: &filters.FieldFilter{
				Field:  "number_field",
				Filter: filters.NumberNotIN(42, 43),
			},
		}},
		OrExprs: []*filters.Expression{{
			Condition: &filters.FieldFilter{
				Field:  "bool_field",
				Filter: filters.False(),
			},
			OrExprs: []*filters.Expression{{
				Condition: &filters.FieldFilter{
					Field:  "optional_bool_field",
					Filter: filters.False(),
				},
			}},
		}},
	}))
}

func TestSimpleOrFalse(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{
		StringField: "whatever",
		BoolField:   true,
	}
	assert.False(Match(m, &filters.Expression{
		Condition: &filters.FieldFilter{
			Field:  "string_field",
			Filter: filters.StringEquals("something"),
		},
		OrExprs: []*filters.Expression{{
			Condition: &filters.FieldFilter{
				Field:  "bool_field",
				Filter: filters.False(),
			},
		}},
	}))
}
