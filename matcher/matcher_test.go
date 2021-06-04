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

package matcher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pf "go.linka.cloud/protofilters"
	test "go.linka.cloud/protofilters/tests/pb"
)

func TestFieldFilter(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{StringField: "ok"}
	ok, err := Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"noop": pf.StringEquals("ok"),
	}})
	assert.Error(err)
	assert.False(ok)
	ok, err = Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"messageField": pf.Null(),
	}})
	assert.Error(err)
	assert.False(ok)

	assert.True(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"message_field": pf.Null(),
	}}))
	m.MessageField = m
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"message_field": pf.Null(),
	}}))
	ok, err = Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"message_field.string_field.message_field": pf.Null(),
	}})
	assert.Error(err)
	assert.False(ok)
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"message_field.message_field.message_field": pf.Null(),
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"message_field.message_field.message_field.string_field": pf.StringIN("ok"),
	}}))
}

func TestString(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{StringField: "ok"}
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"string_value_field": pf.StringEquals("ok"),
	}}))
	m.StringValueField = wrapperspb.String("pointer")
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"string_value_field": pf.StringEquals("ok"),
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"string_value_field": pf.StringEquals("pointer"),
	}}))
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"string_value_field": nil,
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"string_field": pf.StringEquals("ok"),
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"string_field": pf.StringIN("other", "ok"),
	}}))
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"string_field": pf.StringIN("other", "noop"),
	}}))
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"string_field": pf.StringNotRegex(`[a-z](.+)`),
	}}))
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"string_value_field": pf.StringNotRegex(`[a-z](.+)`),
	}}))
	m.StringValueField = wrapperspb.String("whatever")
	assert.True(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"string_value_field": pf.StringRegex(`[a-z](.+)`),
	}}))
}

func TestEnum(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{EnumField: test.Test_Type(42)}
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"enum_field": pf.StringIN("OTHER"),
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"enum_field": pf.NumberIN(0, 42),
	}}))
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"enum_field": pf.NumberNotIN(0, 42),
	}}))
	m.EnumField = 0
	assert.True(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"string_field": pf.StringNotIN(),
		"enum_field":   pf.StringIN("NONE"),
	}}))
}

func TestNumber(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{NumberField: 42}
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"number_field": pf.NumberEquals(0),
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"number_field": pf.NumberEquals(42),
	}}))
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"number_field": pf.NumberIN(0, 22),
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"number_field": pf.NumberInf(43),
	}}))
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"number_field": pf.NumberInf(41),
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"number_field": pf.NumberSup(41),
	}}))
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"number_field": pf.NumberSup(43),
	}}))
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"number_value_field": pf.NumberSup(41),
	}}))
	m.NumberValueField = wrapperspb.Int64(42)
	assert.False(Match(m, &pf.FieldsFilter{Filters: pf.Filters{
		"number_value_field": pf.NumberSup(43),
	}}))
}

func TestDuration(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{DurationValueField: durationpb.New(42)}
	assert.True(MatchFilters(m, &pf.FieldFilter{
		Field:  "duration_value_field",
		Filter: pf.DurationEquals(42),
	}))
	assert.True(MatchFilters(m, &pf.FieldFilter{
		Field:  "duration_value_field",
		Filter: pf.DurationInf(43),
	}))
	assert.False(MatchFilters(m, &pf.FieldFilter{
		Field:  "duration_value_field",
		Filter: pf.DurationInf(41),
	}))
	assert.False(MatchFilters(m, &pf.FieldFilter{
		Field:  "duration_value_field",
		Filter: pf.DurationSup(43),
	}))
	assert.True(MatchFilters(m, &pf.FieldFilter{
		Field:  "duration_value_field",
		Filter: pf.DurationSup(41),
	}))
}

func TestTime(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{TimeValueField: timestamppb.New(time.Now().UTC())}
	assert.True(MatchFilters(m, &pf.FieldFilter{
		Field:  "time_value_field",
		Filter: pf.TimeEquals(m.TimeValueField.AsTime()),
	}))
	assert.True(MatchFilters(m, &pf.FieldFilter{
		Field:  "time_value_field",
		Filter: pf.TimeAfter(m.TimeValueField.AsTime().Add(-time.Second)),
	}))
	assert.False(MatchFilters(m, &pf.FieldFilter{
		Field:  "time_value_field",
		Filter: pf.TimeAfter(m.TimeValueField.AsTime().Add(time.Second)),
	}))
	assert.False(MatchFilters(m, &pf.FieldFilter{
		Field:  "time_value_field",
		Filter: pf.TimeBefore(m.TimeValueField.AsTime().Add(-time.Second)),
	}))
	assert.True(MatchFilters(m, &pf.FieldFilter{
		Field:  "time_value_field",
		Filter: pf.TimeBefore(m.TimeValueField.AsTime().Add(time.Second)),
	}))
}
