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
)

type FieldFilterer interface {
	Expr() *Expression
}

func NullField(field string) NullFilterer {
	return nullFieldFilterer{field: field}
}

type NullFilterer interface {
	Null() *FieldFilter
}

func StringField(field string) StringFilterer {
	return stringFieldFilter{field: field}
}

func NullableStringField(field string) NullableStringFilterer {
	return stringFieldFilter{field: field}
}

type StringFilterer interface {
	Equals(s string) *FieldFilter
	NotEquals(s string) *FieldFilter
	IEquals(s string) *FieldFilter
	NotIEquals(s string) *FieldFilter
	Regex(s string) *FieldFilter
	NotRegex(s string) *FieldFilter
	IN(s ...string) *FieldFilter
	NotIN(s ...string) *FieldFilter
}

type NullableStringFilterer interface {
	StringFilterer
	NullFilterer
}

func NumberField(field string) NumberFilterer {
	return numberFieldFilterer{field: field}
}

func NullableNumberField(field string) NullableNumberFilterer {
	return numberFieldFilterer{field: field}
}

type NumberFilterer interface {
	Equals(n float64) *FieldFilter
	Inf(n float64) *FieldFilter
	Sup(n float64) *FieldFilter
	IN(n ...float64) *FieldFilter
	NotIN(n ...float64) *FieldFilter
}

type NullableNumberFilterer interface {
	NumberFilterer
	NullFilterer
}

func BoolField(field string) BoolFilterer {
	return boolFieldFilterer{field: field}
}

func NullableBoolField(field string) NullableBoolFilterer {
	return boolFieldFilterer{field: field}
}

type BoolFilterer interface {
	True() *FieldFilter
	False() *FieldFilter
}

type NullableBoolFilterer interface {
	BoolFilterer
	NullFilterer
}

func DurationField(field string) DurationFilterer {
	return durationFielFilterer{field: field}
}

func NullableDurationField(field string) NullableDurationFilterer {
	return durationFielFilterer{field: field}
}

type DurationFilterer interface {
	Equals(d time.Duration) *FieldFilter
	NotEquals(d time.Duration) *FieldFilter
	Sup(d time.Duration) *FieldFilter
	Inf(d time.Duration) *FieldFilter
}

type NullableDurationFilterer interface {
	DurationFilterer
	NullFilterer
}

func TimeField(field string) TimeFilterer {
	return timeFieldFilterer{field: field}
}

func NullableTimeField(field string) NullableTimeFilterer {
	return timeFieldFilterer{field: field}
}

type TimeFilterer interface {
	Equals(t time.Time) *FieldFilter
	NotEquals(t time.Time) *FieldFilter
	After(t time.Time) *FieldFilter
	Before(t time.Time) *FieldFilter
}

type NullableTimeFilterer interface {
	TimeFilterer
	NullFilterer
}

type stringFieldFilter struct {
	field string
}

func (f stringFieldFilter) Equals(s string) *FieldFilter {
	return where(f.field, StringEquals(s))
}

func (f stringFieldFilter) NotEquals(s string) *FieldFilter {
	return where(f.field, StringNotEquals(s))
}

func (f stringFieldFilter) IEquals(s string) *FieldFilter {
	return where(f.field, StringIEquals(s))
}

func (f stringFieldFilter) NotIEquals(s string) *FieldFilter {
	return where(f.field, StringNotIEquals(s))
}

func (f stringFieldFilter) Regex(s string) *FieldFilter {
	return where(f.field, StringRegex(s))
}

func (f stringFieldFilter) NotRegex(s string) *FieldFilter {
	return where(f.field, StringNotRegex(s))
}

func (f stringFieldFilter) IN(s ...string) *FieldFilter {
	return where(f.field, StringIN(s...))
}

func (f stringFieldFilter) NotIN(s ...string) *FieldFilter {
	return where(f.field, StringNotIN(s...))
}

func (f stringFieldFilter) Null() *FieldFilter {
	return where(f.field, Null())
}

type numberFieldFilterer struct {
	field string
}

func (f numberFieldFilterer) Equals(n float64) *FieldFilter {
	return where(f.field, NumberEquals(n))
}

func (f numberFieldFilterer) Inf(n float64) *FieldFilter {
	return where(f.field, NumberInf(n))
}

func (f numberFieldFilterer) Sup(n float64) *FieldFilter {
	return where(f.field, NumberSup(n))
}

func (f numberFieldFilterer) IN(n ...float64) *FieldFilter {
	return where(f.field, NumberIN(n...))
}

func (f numberFieldFilterer) NotIN(n ...float64) *FieldFilter {
	return where(f.field, NumberNotIN(n...))
}

func (f numberFieldFilterer) Null() *FieldFilter {
	return where(f.field, Null())
}

type boolFieldFilterer struct {
	field string
}

func (f boolFieldFilterer) True() *FieldFilter {
	return where(f.field, True())
}

func (f boolFieldFilterer) False() *FieldFilter {
	return where(f.field, False())
}

func (f boolFieldFilterer) Null() *FieldFilter {
	return where(f.field, Null())
}

type nullFieldFilterer struct {
	field string
}

func (n nullFieldFilterer) Null() *FieldFilter {
	return where(n.field, Null())
}

type durationFielFilterer struct {
	field string
}

func (f durationFielFilterer) Equals(d time.Duration) *FieldFilter {
	return where(f.field, DurationEquals(d))
}

func (f durationFielFilterer) NotEquals(d time.Duration) *FieldFilter {
	return where(f.field, DurationNotEquals(d))
}

func (f durationFielFilterer) Sup(d time.Duration) *FieldFilter {
	return where(f.field, DurationSup(d))
}

func (f durationFielFilterer) Inf(d time.Duration) *FieldFilter {
	return where(f.field, DurationInf(d))
}

func (f durationFielFilterer) Null() *FieldFilter {
	return where(f.field, Null())
}

type timeFieldFilterer struct {
	field string
}

func (f timeFieldFilterer) Equals(t time.Time) *FieldFilter {
	return where(f.field, TimeEquals(t))
}

func (f timeFieldFilterer) NotEquals(t time.Time) *FieldFilter {
	return where(f.field, TimeNotEquals(t))
}

func (f timeFieldFilterer) After(t time.Time) *FieldFilter {
	return where(f.field, TimeAfter(t))
}

func (f timeFieldFilterer) Before(t time.Time) *FieldFilter {
	return where(f.field, TimeBefore(t))
}

func (f timeFieldFilterer) Null() *FieldFilter {
	return where(f.field, Null())
}
