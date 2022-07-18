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

package test

import (
	"go.linka.cloud/protofilters/filters"
)

type TestFilter struct {
	StringField filters.StringFilterer
	NumberField filters.NumberFilterer
	BoolField   filters.BoolFilterer
	// EnumField
	// MessageField
	// RepeatedStringField
	// RepeatedMessageField
	NumberValueField    filters.NullableNumberFilterer
	StringValueField    filters.NullableStringFilterer
	BoolValueField      filters.NullableBoolFilterer
	TimeValueField      filters.TimeFilterer
	DurationValueField  filters.DurationFilterer
	OptionalStringField filters.StringFilterer
	OptionalNumberField filters.NumberFilterer
	OptionalBoolField   filters.BoolFilterer
	// OptionalEnumField
}

func TestWhere(fn func(f TestFilter) *filters.Expression) *filters.Expression {
	return fn(TestFilters)
}

var TestFilters = TestFilter{
	StringField:         filters.StringField(TestFields.StringField),
	NumberField:         filters.NumberField(TestFields.NumberField),
	BoolField:           filters.BoolField(TestFields.BoolField),
	NumberValueField:    filters.NullableNumberField(TestFields.NumberValueField),
	StringValueField:    filters.NullableStringField(TestFields.StringValueField),
	BoolValueField:      filters.NullableBoolField(TestFields.BoolValueField),
	TimeValueField:      filters.TimeField(TestFields.TimeValueField),
	DurationValueField:  filters.DurationField(TestFields.DurationValueField),
	OptionalStringField: filters.NullableStringField(TestFields.OptionalStringField),
	OptionalNumberField: filters.NullableNumberField(TestFields.OptionalNumberField),
	OptionalBoolField:   filters.NullableBoolField(TestFields.OptionalBoolField),
}
