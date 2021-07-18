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
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"go.linka.cloud/protofilters/filters"
	test "go.linka.cloud/protofilters/tests/pb"
)

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
