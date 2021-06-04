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
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pf "go.linka.cloud/protofilters"
)

type WKType string

func (t WKType) String() string {
	return string(t)
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

type Matcher interface {
	Match(m proto.Message, f *pf.FieldsFilter) (bool, error)
	MatchFilters(m proto.Message, fs ...*pf.FieldFilter) (bool, error)
}

type CachingMatcher interface {
	Matcher
	Clear()
}

var defaultMatcher CachingMatcher = &matcher{cache: make(map[string]pref.FieldDescriptor)}

func Match(msg proto.Message, f *pf.FieldsFilter) (bool, error) {
	return defaultMatcher.Match(msg, f)
}

func MatchFilters(msg proto.Message, fs ...*pf.FieldFilter) (bool, error) {
	return defaultMatcher.MatchFilters(msg, fs...)
}

type matcher struct {
	mu    sync.RWMutex
	cache map[string]pref.FieldDescriptor
}

func (m *matcher) Match(msg proto.Message, f *pf.FieldsFilter) (bool, error) {
	if msg == nil {
		return false, errors.New("message is null")
	}
	if f == nil {
		return true, nil
	}
	for path, filter := range f.Filters {
		fd, err := m.lookup(msg, path)
		if err != nil {
			return false, err
		}
		ok, err := match(msg, fd, filter)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func (m *matcher) MatchFilters(msg proto.Message, fs ...*pf.FieldFilter) (bool, error) {
	f := pf.New(fs...)
	return m.Match(msg, f)
}

func (m *matcher) Clear() {
	m.mu.Lock()
	m.cache = make(map[string]pref.FieldDescriptor)
	m.mu.Unlock()
}

func (m *matcher) lookup(msg proto.Message, path string) (pref.FieldDescriptor, error) {
	if m.cache == nil {
		m.mu.Lock()
		m.cache = make(map[string]pref.FieldDescriptor)
		m.mu.Unlock()
	}
	key := fmt.Sprintf("%s.%s", msg.ProtoReflect().Descriptor().FullName(), path)
	m.mu.RLock()
	fd, ok := m.cache[key]
	m.mu.RUnlock()
	if ok {
		return fd, nil
	}
	md0 := msg.ProtoReflect().Descriptor()
	md := md0
	fd, ok = rangeFields(path, func(field string) (pref.FieldDescriptor, bool) {
		// Search the field within the message.
		if md == nil {
			return nil, false // not within a message
		}
		fd := md.Fields().ByName(pref.Name(field))
		// The real field name of a group is the message name.
		if fd == nil {
			gd := md.Fields().ByName(pref.Name(strings.ToLower(field)))
			if gd != nil && gd.Kind() == pref.GroupKind && string(gd.Message().Name()) == field {
				fd = gd
			}
		} else if fd.Kind() == pref.GroupKind && string(fd.Message().Name()) != field {
			fd = nil
		}
		if fd == nil {
			return nil, false // message does not have this field
		}
		// Identify the next message to search within.
		md = fd.Message() // may be nil

		// Repeated fields are only allowed at the last postion.
		if fd.IsList() || fd.IsMap() {
			md = nil
		}

		return fd, true
	})
	if !ok {
		return nil, fmt.Errorf("%s does not contain '%s'", md0.FullName(), path)
	}
	m.mu.Lock()
	m.cache[key] = fd
	m.mu.Unlock()
	return fd, nil
}

func match(msg proto.Message, fd pref.FieldDescriptor, f *pf.Filter) (bool, error) {
	switch f.GetMatch().(type) {
	case *pf.Filter_String_:
		return matchString(msg, fd, f)
	case *pf.Filter_Number:
		return matchNumber(msg, fd, f)
	case *pf.Filter_Bool:
		return matchBool(msg, fd, f)
	case *pf.Filter_Null:
		return matchNull(msg, fd, f)
	case *pf.Filter_Time:
		return matchTime(msg, fd, f)
	case *pf.Filter_Duration:
		return matchDuration(msg, fd, f)
	}
	return false, nil
}

func matchDuration(m proto.Message, fd pref.FieldDescriptor, f *pf.Filter) (bool, error) {
	if fd.Kind() != pref.MessageKind || WKType(fd.Message().FullName()) != Duration {
		return false, fmt.Errorf("cannot use duration filter on %s", fd.Kind().String())
	}
	rval := m.ProtoReflect().Get(fd)
	if !m.ProtoReflect().Has(fd) {
		if f.GetDuration().GetNot() {
			return true, nil
		}
		return false, nil
	}
	rval.Message().Get(fd.Message().Fields().Get(0))
	t1 := (&durationpb.Duration{
		Seconds: rval.Message().Get(fd.Message().Fields().Get(0)).Int(),
		Nanos:   int32(rval.Message().Get(fd.Message().Fields().Get(1)).Int()),
	}).AsDuration()
	var match bool
	switch f.GetDuration().GetCondition().(type) {
	case *pf.DurationFilter_Equals:
		match = t1 == f.GetDuration().GetEquals().AsDuration()
	case *pf.DurationFilter_Inf:
		match = t1 < f.GetDuration().GetInf().AsDuration()
	case *pf.DurationFilter_Sup:
		match = t1 > f.GetDuration().GetSup().AsDuration()
	}
	if f.GetDuration().GetNot() {
		return !match, nil
	}
	return match, nil
}

func matchTime(m proto.Message, fd pref.FieldDescriptor, f *pf.Filter) (bool, error) {
	if fd.Kind() != pref.MessageKind || WKType(fd.Message().FullName()) != Timestamp {
		return false, fmt.Errorf("cannot use time filter on %s", fd.Kind().String())
	}
	rval := m.ProtoReflect().Get(fd)
	if !m.ProtoReflect().Has(fd) {
		if f.GetTime().GetNot() {
			return true, nil
		}
		return false, nil
	}
	t1 := (&timestamppb.Timestamp{
		Seconds: rval.Message().Get(fd.Message().Fields().Get(0)).Int(),
		Nanos:   int32(rval.Message().Get(fd.Message().Fields().Get(1)).Int()),
	}).AsTime()
	var match bool
	switch f.GetTime().GetCondition().(type) {
	case *pf.TimeFilter_Equals:
		match = t1.Equal(f.GetTime().GetEquals().AsTime().UTC())
	case *pf.TimeFilter_Before:
		match = t1.Before(f.GetTime().GetBefore().AsTime().UTC())
	case *pf.TimeFilter_After:
		match = t1.After(f.GetTime().GetAfter().AsTime().UTC())
	}
	if f.GetTime().GetNot() {
		return !match, nil
	}
	return match, nil
}

func matchNull(m proto.Message, fd pref.FieldDescriptor, f *pf.Filter) (bool, error) {
	var match bool
	switch fd.Kind() {
	case pref.MessageKind:
		match = !m.ProtoReflect().Has(fd)
	case pref.GroupKind:
		match = m.ProtoReflect().Get(fd).List().Len() == 0
	default:
		return false, fmt.Errorf("cannot use null filter on %s", fd.Kind().String())
	}
	if f.GetNull().GetNot() {
		return !match, nil
	}
	return match, nil
}

func matchBool(m proto.Message, fd pref.FieldDescriptor, f *pf.Filter) (bool, error) {
	rval := m.ProtoReflect().Get(fd)
	var val *bool
	if fd.Kind() != pref.BoolKind {
		if fd.Kind() != pref.MessageKind || WKType(fd.Message().FullName()) != BoolValue {
			return false, fmt.Errorf("cannot use bool filter on %s", fd.Kind().String())
		}
		// return early as the condition will always be false
		if !m.ProtoReflect().Has(fd) {
			return false, nil
		}
		val = proto.Bool(rval.Message().Get(fd.Message().Fields().Get(0)).Bool())
	}
	if val == nil {
		val = proto.Bool(rval.Bool())
	}
	return *val == f.GetBool().GetEquals(), nil
}

func matchString(m proto.Message, fd pref.FieldDescriptor, f *pf.Filter) (bool, error) {
	rval := m.ProtoReflect().Get(fd)
	var stringValue *string
	if fd.Kind() != pref.StringKind && fd.Kind() != pref.EnumKind {
		if fd.Kind() != pref.MessageKind || WKType(fd.Message().FullName()) != StringValue {
			return false, fmt.Errorf("cannot use string filter on %s", fd.Kind().String())
		}
		// return early as the condition will always be false
		if !m.ProtoReflect().Has(fd) {
			if f.GetString_().GetNot() {
				return true, nil
			}
			return false, nil
		}
		stringValue = proto.String(rval.Message().Get(fd.Message().Fields().Get(0)).String())
	}
	insensitive := f.GetString_().GetCaseInsensitive()
	value := rval.String()
	if fd.Kind() == pref.EnumKind {
		e := fd.Enum().Values().ByNumber(rval.Enum())
		if e == nil {
			return false, nil
		}
		value = string(e.Name())
	} else if stringValue != nil {
		value = *stringValue
	}
	var match bool
	switch f.GetString_().GetCondition().(type) {
	case *pf.StringFilter_Equals:
		if insensitive {
			match = strings.ToLower(f.GetString_().GetEquals()) == strings.ToLower(value)
		} else {
			match = value == f.GetString_().GetEquals()
		}
	case *pf.StringFilter_Regex:
		reg, err := regexp.Compile(f.GetString_().GetRegex())
		if err != nil {
			return false, err
		}
		match = reg.MatchString(value)
	case *pf.StringFilter_In_:
	lookup:
		for _, v := range f.GetString_().GetIn().GetValues() {
			if (insensitive && strings.ToLower(v) == strings.ToLower(value)) || v == value {
				match = true
				break lookup
			}
		}
	}
	if f.GetString_().GetNot() {
		return !match, nil
	}
	return match, nil
}

func matchNumber(m proto.Message, fd pref.FieldDescriptor, f *pf.Filter) (bool, error) {
	rval := m.ProtoReflect().Get(fd)
	var val float64
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
			if !m.ProtoReflect().Has(fd) {
				return false, nil
			}
			val = rval.Message().Get(fd.Message().Fields().Get(0)).Float()
		case Int64Value, Int32Value:
			if !m.ProtoReflect().Has(fd) {
				return false, nil
			}
			val = float64(rval.Message().Get(fd.Message().Fields().Get(0)).Int())
		case UInt64Value, UInt32Value:
			if !m.ProtoReflect().Has(fd) {
				return false, nil
			}
			val = float64(rval.Message().Get(fd.Message().Fields().Get(0)).Uint())
		default:
			return false, fmt.Errorf("cannot use number filter on %s", fd.Kind().String())
		}
	default:
		return false, fmt.Errorf("cannot use number filter on %s", fd.Kind().String())
	}
	var match bool
	switch f.GetNumber().GetCondition().(type) {
	case *pf.NumberFilter_Equals:
		match = val == f.GetNumber().GetEquals()
	case *pf.NumberFilter_Inf:
		match = val < f.GetNumber().GetInf()
	case *pf.NumberFilter_Sup:
		match = val > f.GetNumber().GetSup()
	case *pf.NumberFilter_In_:
	lookup:
		for _, v := range f.GetNumber().GetIn().GetValues() {
			if val == v {
				match = true
				break lookup
			}
		}
	}
	if f.GetNumber().GetNot() {
		return !match, nil
	}
	return match, nil
}

// rangeFields is like strings.Split(path, "."), but avoids allocations by
// iterating over each field in place and calling a iterator function.
// (taken from "google.golang.org/protobuf/types/known/fieldmaskpb")
func rangeFields(path string, f func(field string) (pref.FieldDescriptor, bool)) (pref.FieldDescriptor, bool) {
	for {
		var field string
		if i := strings.IndexByte(path, '.'); i >= 0 {
			field, path = path[:i], path[i:]
		} else {
			field, path = path, ""
		}
		v, ok := f(field)
		if !ok {
			return nil, false
		}
		if len(path) == 0 {
			return v, true
		}
		path = strings.TrimPrefix(path, ".")
	}
}
