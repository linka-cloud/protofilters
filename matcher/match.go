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
	"fmt"

	"google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pf "go.linka.cloud/protofilters"
)

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

func matchString(m proto.Message, fd pref.FieldDescriptor, f *pf.Filter) (bool, error) {
	rval := m.ProtoReflect().Get(fd)
	var value *string
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
		value = proto.String(rval.Message().Get(fd.Message().Fields().Get(0)).String())
	}
	if fd.Kind() == pref.EnumKind {
		e := fd.Enum().Values().ByNumber(rval.Enum())
		if e == nil {
			return false, nil
		}
		value = proto.String(string(e.Name()))
	} else if value == nil {
		value = proto.String(rval.String())
	}
	return f.GetString_().Match(value)
}

func matchNumber(m proto.Message, fd pref.FieldDescriptor, f *pf.Filter) (bool, error) {
	rval := m.ProtoReflect().Get(fd)
	var val *float64
	switch fd.Kind() {
	case pref.Int32Kind,
		pref.Sint32Kind,
		pref.Int64Kind,
		pref.Sint64Kind,
		pref.Sfixed32Kind,
		pref.Fixed32Kind,
		pref.Sfixed64Kind,
		pref.Fixed64Kind:
		val = proto.Float64(float64(rval.Int()))
	case pref.Uint32Kind, pref.Uint64Kind:
		val = proto.Float64(float64(rval.Uint()))
	case pref.FloatKind, pref.DoubleKind:
		val = proto.Float64(rval.Float())
	case pref.EnumKind:
		val = proto.Float64(float64(rval.Enum()))
	case pref.MessageKind:
		if m.ProtoReflect().Has(fd) {
			switch WKType(fd.Message().FullName()) {
			case DoubleValue, FloatValue:
				val = proto.Float64(rval.Message().Get(fd.Message().Fields().Get(0)).Float())
			case Int64Value, Int32Value:
				val = proto.Float64(float64(rval.Message().Get(fd.Message().Fields().Get(0)).Int()))
			case UInt64Value, UInt32Value:
				val = proto.Float64(float64(rval.Message().Get(fd.Message().Fields().Get(0)).Uint()))
			default:
				return false, fmt.Errorf("cannot use number filter on %s", fd.Kind().String())
			}
		} else {
			val = nil
		}
	default:
		return false, fmt.Errorf("cannot use number filter on %s", fd.Kind().String())
	}
	return f.GetNumber().Match(val)
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
	return f.GetBool().Match(val)
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
	return f.GetTime().Match(&timestamppb.Timestamp{
		Seconds: rval.Message().Get(fd.Message().Fields().Get(0)).Int(),
		Nanos:   int32(rval.Message().Get(fd.Message().Fields().Get(1)).Int()),
	})
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

	return f.GetDuration().Match(&durationpb.Duration{
		Seconds: rval.Message().Get(fd.Message().Fields().Get(0)).Int(),
		Nanos:   int32(rval.Message().Get(fd.Message().Fields().Get(1)).Int()),
	})
}
