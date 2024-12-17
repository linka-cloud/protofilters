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
	"sort"
	"strings"

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

func where(field string, f *Filter) *FieldFilter {
	return &FieldFilter{
		Field:  field,
		Filter: f,
	}
}

// Field joins the parts as un field path, e.g. Field("message", "string_field") returns "message.string_field"
func Field(parts ...string) string {
	return strings.Join(parts, ".")
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
			return strings.EqualFold(x.GetEquals(), value), nil
		}
		return value == x.GetEquals(), nil
	case *StringFilter_HasPrefix:
		if insensitive {
			return strings.HasPrefix(strings.ToLower(value), strings.ToLower(x.GetHasPrefix())), nil
		}
		return strings.HasPrefix(value, x.GetHasPrefix()), nil
	case *StringFilter_HasSuffix:
		if insensitive {
			return strings.HasSuffix(strings.ToLower(value), strings.ToLower(x.GetHasSuffix())), nil
		}
		return strings.HasSuffix(value, x.GetHasSuffix()), nil
	case *StringFilter_Regex:
		reg, err := regexp.Compile(x.GetRegex())
		if err != nil {
			return false, err
		}
		return reg.MatchString(value), nil
	case *StringFilter_In_:
		for _, v := range x.GetIn().GetValues() {
			if (insensitive && strings.EqualFold(v, value)) || v == value {
				return true, nil
			}
		}
	case *StringFilter_Inf:
		if insensitive {
			return strings.ToLower(value) < strings.ToLower(x.GetInf()), nil
		}
		return value < x.GetInf(), nil
	case *StringFilter_Sup:
		if insensitive {
			return strings.ToLower(value) > strings.ToLower(x.GetSup()), nil
		}
		return value > x.GetSup(), nil
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

func (x *Expression) Fields() (fields []string) {
	if x == nil || x.Condition == nil {
		return nil
	}
	m := make(map[string]struct{})
	m[x.Condition.Field] = struct{}{}
	for _, v := range x.AndExprs {
		for _, v := range v.Fields() {
			m[v] = struct{}{}
		}
	}
	for _, v := range x.OrExprs {
		for _, v := range v.Fields() {
			m[v] = struct{}{}
		}
	}
	for k := range m {
		fields = append(fields, k)
	}
	sort.Strings(fields)
	return fields
}

func (x *Expression) FieldFilters() (fieldFilters []*FieldFilter) {
	if x == nil || x.Condition == nil {
		return nil
	}
	fieldFilters = append(fieldFilters, x.Condition)
	for _, v := range x.AndExprs {
		fieldFilters = append(fieldFilters, v.FieldFilters()...)
	}
	for _, v := range x.OrExprs {
		fieldFilters = append(fieldFilters, v.FieldFilters()...)
	}
	return fieldFilters
}

// Expr is a convenient method so that both Expression and FieldFilter
// implement the FieldFilterer interface
func (x *Expression) Expr() *Expression {
	return x
}

func (x *FieldsFilter) Expr() *Expression {
	if len(x.Filters) == 0 {
		return nil
	}
	e := &Expression{}
	for k, v := range x.Filters {
		if e.Condition == nil {
			e.Condition = &FieldFilter{Field: k, Filter: v}
			continue
		}
		e.AndExprs = append(e.AndExprs, &Expression{Condition: &FieldFilter{Field: k, Filter: v}})
	}
	return e
}
