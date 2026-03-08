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

package reflect

import (
	"fmt"
	"regexp"
	"strings"

	pref "google.golang.org/protobuf/reflect/protoreflect"

	"go.linka.cloud/protofilters/filters"
)

// WKType represents a google.protobuf well-known type
type WKType string

// String implements the Stringer interface
func (t WKType) String() string {
	return string(t)
}

var WKTypes = []WKType{
	Timestamp,
	Duration,
	DoubleValue,
	FloatValue,
	Int64Value,
	UInt64Value,
	Int32Value,
	UInt32Value,
	BoolValue,
	StringValue,
	BytesValue,
}

func IsWKType(t pref.FullName) bool {
	for _, wt := range WKTypes {
		if wt == WKType(t) {
			return true
		}
	}
	return false
}

const (
	Timestamp   WKType = "google.protobuf.Timestamp"
	Duration    WKType = "google.protobuf.Duration"
	DoubleValue WKType = "google.protobuf.DoubleValue"
	FloatValue  WKType = "google.protobuf.FloatValue"
	Int64Value  WKType = "google.protobuf.Int64Value"
	UInt64Value WKType = "google.protobuf.UInt64Value"
	Int32Value  WKType = "google.protobuf.Int32Value"
	UInt32Value WKType = "google.protobuf.UInt32Value"
	BoolValue   WKType = "google.protobuf.BoolValue"
	StringValue WKType = "google.protobuf.StringValue"
	BytesValue  WKType = "google.protobuf.BytesValue"
)

func Match(val pref.Value, fd pref.FieldDescriptor, f *filters.Filter) (bool, error) {
	switch f.GetMatch().(type) {
	case *filters.Filter_String_:
		return matchString(val, fd, f)
	case *filters.Filter_Number:
		return matchNumber(val, fd, f)
	case *filters.Filter_Bool:
		return matchBool(val, fd, f)
	case *filters.Filter_Null:
		return matchNull(val, fd, f)
	case *filters.Filter_Time:
		return matchTime(val, fd, f)
	case *filters.Filter_Duration:
		return matchDuration(val, fd, f)
	}
	return false, nil
}

func matchString(rval pref.Value, fd pref.FieldDescriptor, f *filters.Filter) (bool, error) {
	var value string
	hasValue := true
	valueSet := false
	if fd.Kind() != pref.StringKind && fd.Kind() != pref.EnumKind {
		if fd.Kind() != pref.MessageKind || WKType(fd.Message().FullName()) != StringValue {
			return false, fmt.Errorf("cannot use string filter on %s", fd.Kind().String())
		}
		// return early as the condition will always be false
		if !rval.IsValid() {
			return checkNot(f, false, nil)
		}
		value = rval.Message().Get(fd.Message().Fields().Get(0)).String()
		valueSet = true
	}
	if fd.Kind() == pref.EnumKind {
		e := fd.Enum().Values().ByNumber(rval.Enum())
		if e == nil {
			return false, nil
		}
		value = string(e.Name())
		valueSet = true
	} else if !valueSet {
		value = rval.String()
		if !rval.IsValid() {
			hasValue = false
		}
	}
	match, err := matchStringFilter(f.GetString_(), value, hasValue)
	return checkNot(f, match, err)
}

func matchNumber(rval pref.Value, fd pref.FieldDescriptor, f *filters.Filter) (bool, error) {
	// fast path for float64
	if val, ok := rval.Interface().(float64); ok {
		match, err := matchNumberFilter(f.GetNumber(), val, true)
		return checkNot(f, match, err)
	}
	var val float64
	hasValue := true
	if !rval.IsValid() {
		if !fd.HasOptionalKeyword() {
			return false, fmt.Errorf("cannot use number filter on %s", fd.Kind().String())
		}
		hasValue = false
	}
	if hasValue {
		switch fd.Kind() {
		case pref.Int32Kind,
			pref.Sint32Kind,
			pref.Int64Kind,
			pref.Sint64Kind,
			pref.Sfixed32Kind,
			pref.Fixed32Kind,
			pref.Sfixed64Kind,
			pref.Fixed64Kind:
			val = float64(rval.Int())
		case pref.Uint32Kind, pref.Uint64Kind:
			val = float64(rval.Uint())
		case pref.FloatKind, pref.DoubleKind:
			val = rval.Float()
		case pref.EnumKind:
			val = float64(rval.Enum())
		case pref.MessageKind:
			switch WKType(fd.Message().FullName()) {
			case DoubleValue, FloatValue:
				val = rval.Message().Get(fd.Message().Fields().Get(0)).Float()
			case Int64Value, Int32Value:
				val = float64(rval.Message().Get(fd.Message().Fields().Get(0)).Int())
			case UInt64Value, UInt32Value:
				val = float64(rval.Message().Get(fd.Message().Fields().Get(0)).Uint())
			default:
				return false, fmt.Errorf("cannot use number filter on %s", fd.Kind().String())
			}
		default:
			return false, fmt.Errorf("cannot use number filter on %s", fd.Kind().String())
		}
	}
	match, err := matchNumberFilter(f.GetNumber(), val, hasValue)
	return checkNot(f, match, err)
}

func matchBool(rval pref.Value, fd pref.FieldDescriptor, f *filters.Filter) (bool, error) {
	var val bool
	hasValue := true
	if fd.Kind() != pref.BoolKind {
		if fd.Kind() != pref.MessageKind || WKType(fd.Message().FullName()) != BoolValue {
			return false, fmt.Errorf("cannot use bool filter on %s", fd.Kind().String())
		}
		// return early as the condition will always be false
		if !rval.IsValid() {
			return checkNot(f, false, nil)
		}
		val = rval.Message().Get(fd.Message().Fields().Get(0)).Bool()
	} else if !fd.HasOptionalKeyword() || rval.IsValid() {
		val = rval.Bool()
	} else {
		hasValue = false
	}
	match, err := matchBoolFilter(f.GetBool(), val, hasValue)
	return checkNot(f, match, err)
}

func matchNull(rval pref.Value, fd pref.FieldDescriptor, f *filters.Filter) (bool, error) {
	var match bool
	switch fd.Kind() {
	case pref.MessageKind:
		match = !rval.Message().IsValid()
	case pref.GroupKind:
		match = rval.List().Len() == 0
	default:
		if !fd.HasOptionalKeyword() {
			return false, fmt.Errorf("cannot use null filter on %s", fd.Kind().String())
		}
		match = !rval.IsValid()
	}
	return checkNot(f, match, nil)
}

func matchTime(rval pref.Value, fd pref.FieldDescriptor, f *filters.Filter) (bool, error) {
	if fd.Kind() != pref.MessageKind || WKType(fd.Message().FullName()) != Timestamp {
		return false, fmt.Errorf("cannot use time filter on %s", fd.Kind().String())
	}
	if !rval.IsValid() {
		return checkNot(f, false, nil)
	}
	seconds := rval.Message().Get(fd.Message().Fields().Get(0)).Int()
	nanos := int32(rval.Message().Get(fd.Message().Fields().Get(1)).Int())
	match, err := matchTimeFilter(f.GetTime(), seconds, nanos, true)
	return checkNot(f, match, err)
}

func matchDuration(rval pref.Value, fd pref.FieldDescriptor, f *filters.Filter) (bool, error) {
	if fd.Kind() != pref.MessageKind || WKType(fd.Message().FullName()) != Duration {
		return false, fmt.Errorf("cannot use duration filter on %s", fd.Kind().String())
	}
	if !rval.IsValid() {
		return checkNot(f, false, nil)
	}
	seconds := rval.Message().Get(fd.Message().Fields().Get(0)).Int()
	nanos := int32(rval.Message().Get(fd.Message().Fields().Get(1)).Int())
	match, err := matchDurationFilter(f.GetDuration(), seconds, nanos, true)
	return checkNot(f, match, err)
}

func matchStringFilter(f *filters.StringFilter, value string, hasValue bool) (bool, error) {
	if !hasValue {
		return false, nil
	}
	insensitive := f.GetCaseInsensitive()
	switch f.GetCondition().(type) {
	case *filters.StringFilter_Equals:
		if insensitive {
			return strings.EqualFold(f.GetEquals(), value), nil
		}
		return value == f.GetEquals(), nil
	case *filters.StringFilter_HasPrefix:
		if insensitive {
			return strings.HasPrefix(strings.ToLower(value), strings.ToLower(f.GetHasPrefix())), nil
		}
		return strings.HasPrefix(value, f.GetHasPrefix()), nil
	case *filters.StringFilter_HasSuffix:
		if insensitive {
			return strings.HasSuffix(strings.ToLower(value), strings.ToLower(f.GetHasSuffix())), nil
		}
		return strings.HasSuffix(value, f.GetHasSuffix()), nil
	case *filters.StringFilter_Regex:
		reg, err := regexp.Compile(f.GetRegex())
		if err != nil {
			return false, err
		}
		return reg.MatchString(value), nil
	case *filters.StringFilter_In_:
		for _, v := range f.GetIn().GetValues() {
			if (insensitive && strings.EqualFold(v, value)) || v == value {
				return true, nil
			}
		}
	case *filters.StringFilter_Inf:
		if insensitive {
			return strings.ToLower(value) < strings.ToLower(f.GetInf()), nil
		}
		return value < f.GetInf(), nil
	case *filters.StringFilter_Sup:
		if insensitive {
			return strings.ToLower(value) > strings.ToLower(f.GetSup()), nil
		}
		return value > f.GetSup(), nil
	}
	return false, nil
}

func matchNumberFilter(f *filters.NumberFilter, value float64, hasValue bool) (bool, error) {
	if !hasValue {
		return false, nil
	}
	switch f.GetCondition().(type) {
	case *filters.NumberFilter_Equals:
		return value == f.GetEquals(), nil
	case *filters.NumberFilter_Inf:
		return value < f.GetInf(), nil
	case *filters.NumberFilter_Sup:
		return value > f.GetSup(), nil
	case *filters.NumberFilter_In_:
		for _, v := range f.GetIn().GetValues() {
			if value == v {
				return true, nil
			}
		}
	}
	return false, nil
}

func matchBoolFilter(f *filters.BoolFilter, value bool, hasValue bool) (bool, error) {
	if !hasValue {
		return false, nil
	}
	return value == f.GetEquals(), nil
}

func matchTimeFilter(f *filters.TimeFilter, seconds int64, nanos int32, hasValue bool) (bool, error) {
	if !hasValue {
		return false, nil
	}
	switch f.GetCondition().(type) {
	case *filters.TimeFilter_Equals:
		t := f.GetEquals()
		return seconds == t.GetSeconds() && nanos == t.GetNanos(), nil
	case *filters.TimeFilter_Before:
		t := f.GetBefore()
		if seconds != t.GetSeconds() {
			return seconds < t.GetSeconds(), nil
		}
		return nanos < t.GetNanos(), nil
	case *filters.TimeFilter_After:
		t := f.GetAfter()
		if seconds != t.GetSeconds() {
			return seconds > t.GetSeconds(), nil
		}
		return nanos > t.GetNanos(), nil
	}
	return false, nil
}

func matchDurationFilter(f *filters.DurationFilter, seconds int64, nanos int32, hasValue bool) (bool, error) {
	if !hasValue {
		return false, nil
	}
	switch f.GetCondition().(type) {
	case *filters.DurationFilter_Equals:
		d := f.GetEquals()
		return seconds == d.GetSeconds() && nanos == d.GetNanos(), nil
	case *filters.DurationFilter_Inf:
		d := f.GetInf()
		if seconds != d.GetSeconds() {
			return seconds < d.GetSeconds(), nil
		}
		return nanos < d.GetNanos(), nil
	case *filters.DurationFilter_Sup:
		d := f.GetSup()
		if seconds != d.GetSeconds() {
			return seconds > d.GetSeconds(), nil
		}
		return nanos > d.GetNanos(), nil
	}
	return false, nil
}

func checkNot(f *filters.Filter, match bool, err error) (bool, error) {
	if f.GetNot() {
		return !match, err
	}
	return match, err
}
