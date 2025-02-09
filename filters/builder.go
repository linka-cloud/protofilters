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

type Builder interface {
	FieldFilterer
	And(field ...string) Builder
	Or(e FieldFilterer) Builder
	StringEquals(s string) Builder
	StringNotEquals(s string) Builder
	StringNotIEquals(s string) Builder
	StringIEquals(s string) Builder
	StringHasPrefix(s string) Builder
	StringNotHasPrefix(s string) Builder
	StringIHasPrefix(s string) Builder
	StringNotIHasPrefix(s string) Builder
	StringHasSuffix(s string) Builder
	StringNotHasSuffix(s string) Builder
	StringIHasSuffix(s string) Builder
	StringNotIHasSuffix(s string) Builder
	StringRegex(s string) Builder
	StringNotRegex(s string) Builder
	StringIN(s ...string) Builder
	StringInf(s string) Builder
	StringSup(s string) Builder
	StringIInf(s string) Builder
	StringISup(s string) Builder
	StringNotIN(s ...string) Builder
	NumberEquals(n float64) Builder
	NumberNotEquals(n float64) Builder
	NumberInf(n float64) Builder
	NumberSup(n float64) Builder
	NumberIN(n ...float64) Builder
	NumberNotIN(n ...float64) Builder
	True() Builder
	False() Builder
	Null() Builder
	NotNull() Builder
	DurationEquals(d time.Duration) Builder
	DurationNotEquals(d time.Duration) Builder
	DurationSup(d time.Duration) Builder
	DurationInf(d time.Duration) Builder
	TimeEquals(t time.Time) Builder
	TimeNotEquals(t time.Time) Builder
	TimeAfter(t time.Time) Builder
	TimeBefore(t time.Time) Builder

	Clone() Builder
}

func Where(field string) Builder {
	r := &Expression{
		Condition: &FieldFilter{
			Field: field,
		},
	}
	return &builder{
		r: r,
		c: r,
	}
}

type builder struct {
	r *Expression
	c *Expression
}

func (b *builder) Expr() *Expression {
	return b.r
}

func (b *builder) And(field ...string) Builder {
	e := &Expression{
		Condition: &FieldFilter{
			Field: Field(field...),
		},
	}
	b.c = e
	b.r.AndExprs = append(b.r.AndExprs, b.c)
	return b
}

func (b *builder) Or(e FieldFilterer) Builder {
	b.r.OrExprs = append(b.r.OrExprs, e.Expr())
	b.c = b.r
	return b
}

// StringEquals constructs a string equals filter
func (b *builder) StringEquals(s string) Builder {
	b.c.Condition.Filter = StringEquals(s)
	return b
}

// StringNotEquals constructs a string not equals filter
func (b *builder) StringNotEquals(s string) Builder {
	b.c.Condition.Filter = StringNotEquals(s)
	return b
}

// StringNotIEquals constructs a case insensitive string not equals filter
func (b *builder) StringNotIEquals(s string) Builder {
	b.c.Condition.Filter = StringNotIEquals(s)
	return b
}

// StringIEquals constructs a case insensitive string equals filter
func (b *builder) StringIEquals(s string) Builder {
	b.c.Condition.Filter = StringIEquals(s)
	return b
}

// StringHasPrefix constructs a string match prefix filter
func (b *builder) StringHasPrefix(s string) Builder {
	b.c.Condition.Filter = StringHasPrefix(s)
	return b
}

// StringNotHasPrefix constructs a string not match prefix filter
func (b *builder) StringNotHasPrefix(s string) Builder {
	b.c.Condition.Filter = StringNotHasPrefix(s)
	return b
}

// StringIHasPrefix constructs a case insensitive string match prefix filter
func (b *builder) StringIHasPrefix(s string) Builder {
	b.c.Condition.Filter = StringIHasPrefix(s)
	return b
}

// StringNotIHasPrefix constructs a case insensitive string not match prefix filter
func (b *builder) StringNotIHasPrefix(s string) Builder {
	b.c.Condition.Filter = StringNotIHasPrefix(s)
	return b
}

// StringHasSuffix constructs a string match suffix filter
func (b *builder) StringHasSuffix(s string) Builder {
	b.c.Condition.Filter = StringHasSuffix(s)
	return b
}

// StringNotHasSuffix constructs a string not match suffix filter
func (b *builder) StringNotHasSuffix(s string) Builder {
	b.c.Condition.Filter = StringNotHasSuffix(s)
	return b
}

// StringIHasSuffix constructs a case insensitive string match suffix filter
func (b *builder) StringIHasSuffix(s string) Builder {
	b.c.Condition.Filter = StringIHasSuffix(s)
	return b
}

// StringNotIHasSuffix constructs a case insensitive string not match suffix filter
func (b *builder) StringNotIHasSuffix(s string) Builder {
	b.c.Condition.Filter = StringNotIHasSuffix(s)
	return b
}

// StringRegex constructs a string match regex filter
func (b *builder) StringRegex(s string) Builder {
	b.c.Condition.Filter = StringRegex(s)
	return b
}

// StringNotRegex constructs a string not match regex filter
func (b *builder) StringNotRegex(s string) Builder {
	b.c.Condition.Filter = StringNotRegex(s)
	return b
}

// StringIN constructs a string in slice filter
func (b *builder) StringIN(s ...string) Builder {
	b.c.Condition.Filter = StringIN(s...)
	return b
}

// StringNotIN constructs a string not in slice filter
func (b *builder) StringNotIN(s ...string) Builder {
	b.c.Condition.Filter = StringNotIN(s...)
	return b
}

// StringInf constructs a string inferior filter
func (b *builder) StringInf(s string) Builder {
	b.c.Condition.Filter = StringInf(s)
	return b
}

// StringSup constructs a string superior filter
func (b *builder) StringSup(s string) Builder {
	b.c.Condition.Filter = StringSup(s)
	return b
}

// StringIInf constructs a case insensitive string inferior filter
func (b *builder) StringIInf(s string) Builder {
	b.c.Condition.Filter = StringIInf(s)
	return b
}

// StringISup constructs a case insensitive string superior filter
func (b *builder) StringISup(s string) Builder {
	b.c.Condition.Filter = StringISup(s)
	return b
}

// NumberEquals constructs a number equals filter
func (b *builder) NumberEquals(n float64) Builder {
	b.c.Condition.Filter = NumberEquals(n)
	return b
}

// NumberNotEquals constructs a number not equals filter
func (b *builder) NumberNotEquals(n float64) Builder {
	b.c.Condition.Filter = NumberNotEquals(n)
	return b
}

// NumberInf constructs a number inferior filter
func (b *builder) NumberInf(n float64) Builder {
	b.c.Condition.Filter = NumberInf(n)
	return b
}

// NumberSup constructs a number superior filter
func (b *builder) NumberSup(n float64) Builder {
	b.c.Condition.Filter = NumberSup(n)
	return b
}

// NumberIN constructs a number in slice filter
func (b *builder) NumberIN(n ...float64) Builder {
	b.c.Condition.Filter = NumberIN(n...)
	return b
}

// NumberNotIN constructs a number not in slice filter
func (b *builder) NumberNotIN(n ...float64) Builder {
	b.c.Condition.Filter = NumberNotIN(n...)
	return b
}

// True constructs a bool is true filter
func (b *builder) True() Builder {
	b.c.Condition.Filter = True()
	return b
}

// False constructs a bool is false filter
func (b *builder) False() Builder {
	b.c.Condition.Filter = False()
	return b
}

// Null constructs a null check filter
func (b *builder) Null() Builder {
	b.c.Condition.Filter = Null()
	return b
}

// NotNull constructs a not null check filter
func (b *builder) NotNull() Builder {
	b.c.Condition.Filter = NotNull()
	return b
}

// DurationEquals constructs a duration equals filter
func (b *builder) DurationEquals(d time.Duration) Builder {
	b.c.Condition.Filter = DurationEquals(d)
	return b
}

// DurationNotEquals constructs a duration not equals filter
func (b *builder) DurationNotEquals(d time.Duration) Builder {
	b.c.Condition.Filter = DurationNotEquals(d)
	return b
}

// DurationSup constructs a duration superior filter
func (b *builder) DurationSup(d time.Duration) Builder {
	b.c.Condition.Filter = DurationSup(d)
	return b
}

// DurationInf constructs a duration inferior filter
func (b *builder) DurationInf(d time.Duration) Builder {
	b.c.Condition.Filter = DurationInf(d)
	return b
}

// TimeEquals constructs a time equals filter
func (b *builder) TimeEquals(t time.Time) Builder {
	b.c.Condition.Filter = TimeEquals(t)
	return b
}

// TimeNotEquals constructs a time not equals filter
func (b *builder) TimeNotEquals(t time.Time) Builder {
	b.c.Condition.Filter = TimeNotEquals(t)
	return b
}

// TimeAfter constructs a time after filter
func (b *builder) TimeAfter(t time.Time) Builder {
	b.c.Condition.Filter = TimeAfter(t)
	return b
}

// TimeBefore constructs a time before filter
func (b *builder) TimeBefore(t time.Time) Builder {
	b.c.Condition.Filter = TimeBefore(t)
	return b
}

func (b *builder) Clone() Builder {
	if b == nil {
		return nil
	}
	return &builder{
		r: b.r.CloneVT(),
		c: b.c.CloneVT(),
	}
}
