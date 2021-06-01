package matcher

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"

	pf "go.linka.cloud/protofilters"
)

type Matcher interface {
	Match(m proto.Message, f *pf.FieldsFilter) (bool, error)
	MatchFilters(m proto.Message, fs ...*pf.FieldFilter) (bool, error)
}

type CachingMatcher interface {
	Matcher
	ResetCache()
}

var defaultMatcher CachingMatcher = &matcher{cache: make(map[string]pref.FieldDescriptor)}

func Match(m proto.Message, f *pf.FieldsFilter) (bool, error) {
	return defaultMatcher.Match(m, f)
}

func MatchFilters(m proto.Message, fs ...*pf.FieldFilter) (bool, error) {
	return defaultMatcher.MatchFilters(m, fs...)
}

type matcher struct {
	mu sync.RWMutex
	cache map[string]pref.FieldDescriptor
}

func (x *matcher) Match(m proto.Message, f *pf.FieldsFilter) (bool, error) {
	if m == nil {
		return false, errors.New("message is null")
	}
	if f == nil {
		return true, nil
	}
	for path, filter := range f.Filters {
		fd, err := x.lookup(m, path)
		if err != nil {
			return false, err
		}
		ok, err := match(m, fd, filter)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func (x *matcher) MatchFilters(m proto.Message, fs ...*pf.FieldFilter) (bool, error) {
	f := pf.New(fs...)
	return x.Match(m, f)
}

func (x *matcher) ResetCache() {
	x.mu.Lock()
	x.cache = make(map[string]pref.FieldDescriptor)
	x.mu.Unlock()
}

func (x *matcher) lookup(m proto.Message, path string) (pref.FieldDescriptor, error) {
	if x.cache == nil {
		x.mu.Lock()
		x.cache = make(map[string]pref.FieldDescriptor)
		x.mu.Unlock()
	}
	key := fmt.Sprintf("%s.%s", m.ProtoReflect().Descriptor().FullName(), path)
	x.mu.RLock()
	fd, ok := x.cache[key]
	x.mu.RUnlock()
	if ok {
		return fd, nil
	}
	md0 := m.ProtoReflect().Descriptor()
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
	x.mu.Lock()
	x.cache[key] = fd
	x.mu.Unlock()
	return fd, nil
}

func match(m proto.Message, fd pref.FieldDescriptor, f *pf.Filter) (bool, error) {
	switch f.GetMatch().(type) {
	case *pf.Filter_String_:
		return matchString(m, fd, f)
	case *pf.Filter_Number:
		return matchNumber(m, fd, f)
	case *pf.Filter_Bool:
		return matchBool(m, fd, f)
	case *pf.Filter_Null:
		return matchNull(m, fd, f)
	}
	return false, nil
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
	if fd.Kind() != pref.BoolKind {
		return false, fmt.Errorf("cannot use bool filter on %s", fd.Kind().String())
	}
	return m.ProtoReflect().Get(fd).Bool() == f.GetBool().GetEquals(), nil
}

func matchString(m proto.Message, fd pref.FieldDescriptor, f *pf.Filter) (bool, error) {
	if fd.Kind() != pref.StringKind && fd.Kind() != pref.EnumKind {
		return false, fmt.Errorf("cannot use string filter on %s", fd.Kind().String())
	}
	insensitive := f.GetString_().GetCaseInsensitive()
	rval := m.ProtoReflect().Get(fd)
	value := rval.String()
	if fd.Kind() == pref.EnumKind {
		e := fd.Enum().Values().ByNumber(rval.Enum())
		if e == nil {
			return false, nil
		}
		value = string(e.Name())
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
	default:
		return false, fmt.Errorf("cannot use number filter on %s", fd.Kind().String())
	}
	var match bool
	switch f.GetNumber().GetCondition().(type) {
	case *pf.NumberFilter_Equals:
		match = val == f.GetNumber().GetEquals()
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

