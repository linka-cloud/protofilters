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
	"errors"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"

	"go.linka.cloud/protofilters/filters"
)

// WKType represents a google.protobuf well-known type
type WKType string

// String implements the Stringer interface
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

// Matcher provides a way to match proto.Message against protofilters.Filter
type Matcher interface {
	// Match matches to proto.Message against the protofilters.FieldsFilter
	// It returns an error if one of the field path or FieldFilter is invalid
	Match(m proto.Message, f *filters.FieldsFilter) (bool, error)
	// MatchFilters matches to proto.Message against the protofilters.FieldFilter slice
	// It returns an error if one of the field path or FieldFilter is invalid
	MatchFilters(m proto.Message, fs ...*filters.FieldFilter) (bool, error)

	MatchExpression(msg proto.Message, expr *filters.Expression) (bool, error)
}

// CachingMatcher is a Matcher that cache messages field path lookup results
type CachingMatcher interface {
	Matcher
	// Clear clears the lookup cache
	Clear()
}

// NewMatcher creates a CachingMatcher
func NewMatcher() CachingMatcher {
	return &matcher{cache: make(map[string]pref.FieldDescriptor)}
}

var defaultMatcher = NewMatcher()

// Match is a convenient method calling Match on the defaultMatcher
func Match(msg proto.Message, f *filters.FieldsFilter) (bool, error) {
	return defaultMatcher.Match(msg, f)
}

// MatchFilters is a convenient calling MatchFilters on the defaultMatcher
func MatchFilters(msg proto.Message, fs ...*filters.FieldFilter) (bool, error) {
	return defaultMatcher.MatchFilters(msg, fs...)
}

func MatchExpression(msg proto.Message, expr *filters.Expression) (bool, error) {
	return defaultMatcher.MatchExpression(msg, expr)
}

type matcher struct {
	mu    sync.RWMutex
	cache map[string]pref.FieldDescriptor
}

func (m *matcher) MatchExpression(msg proto.Message, expr *filters.Expression) (bool, error) {
	ok, err := m.Match(msg, filters.New(expr.Condition))
	if err != nil {
		return false, err
	}
	if (expr.OrExprs == nil || expr.AndExprs != nil) && !ok {
		return false, nil
	}
	if expr.OrExprs != nil {
		if expr.AndExprs == nil && ok {
			return true, nil
		}
		for _, v := range expr.OrExprs {
			orOk, err := m.MatchExpression(msg, v)
			if err != nil {
				return false, err
			}
			if orOk && expr.AndExprs == nil {
				return true, nil
			}
		}
	}
	if expr.AndExprs != nil {
		for _, v := range expr.AndExprs {
			andOk, err := m.MatchExpression(msg, v)
			if err != nil {
				return false, err
			}
			if !andOk {
				return false, nil
			}
		}
	}
	return true, nil
}

func (m *matcher) Match(msg proto.Message, f *filters.FieldsFilter) (bool, error) {
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
		rval := msg.ProtoReflect().Get(fd)
		if fd.IsList() {
			list := rval.List()
			for i := 0; i < list.Len(); i++ {
				match, err := match(list.Get(i), fd, filter)
				if err != nil {
					return false, err
				}
				if filter.GetNot() && !match {
					return false, nil
				}
				if !filter.GetNot() && match {
					return true, nil
				}
			}
			return false, nil
		}
		if fd.IsMap() {
			return false, errors.New("matching against map is not supported")
		}
		ok, err := match(rval, fd, filter)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func (m *matcher) MatchFilters(msg proto.Message, fs ...*filters.FieldFilter) (bool, error) {
	f := filters.New(fs...)
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
