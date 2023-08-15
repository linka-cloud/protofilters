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
	"strings"

	pref "google.golang.org/protobuf/reflect/protoreflect"
)

func Lookup(msg pref.Message, path string) ([]pref.FieldDescriptor, error) {
	md0 := msg.Descriptor()
	md := md0
	fds, ok := rangeFields(path, func(field string) (pref.FieldDescriptor, bool) {
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
		// may be nil
		md = fd.Message()

		if (fd.IsList() && fd.Kind() != pref.MessageKind) || fd.IsMap() {
			md = nil
		}
		return fd, true
	})
	if !ok {
		return nil, fmt.Errorf("%s does not contain '%s'", md0.FullName(), path)
	}
	return fds, nil
}

// rangeFields is like strings.Split(path, "."), but avoids allocations by
// iterating over each field in place and calling a iterator function.
// (taken from "google.golang.org/protobuf/types/known/fieldmaskpb")
func rangeFields(path string, f func(field string) (pref.FieldDescriptor, bool)) ([]pref.FieldDescriptor, bool) {
	var fds []pref.FieldDescriptor
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
		fds = append(fds, v)
		if len(path) == 0 {
			return fds, true
		}
		path = strings.TrimPrefix(path, ".")
	}
}
