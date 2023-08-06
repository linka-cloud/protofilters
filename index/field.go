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

package index

import (
	"context"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// Field is an indexed field in the index
type Field interface {
	// Value returns the value of the field
	Value() protoreflect.Value
	// Keys returns an iterator over the keys of the messages that have this value for this field
	Keys(ctx context.Context) (Iterator[string], error)
	// Descriptors returns the field descriptors for this field
	Descriptors() []protoreflect.FieldDescriptor
}

// FieldReader is an interface for reading a type fields from the index
type FieldReader interface {
	// Get returns the field for the given field descriptor
	Get(ctx context.Context, f protoreflect.Name) (Iterator[Field], bool, error)
}

type field struct {
	value       protoreflect.Value
	keys        []string
	descriptors []protoreflect.FieldDescriptor
}

func (f *field) Value() protoreflect.Value {
	return f.value
}

func (f *field) Keys(_ context.Context) (Iterator[string], error) {
	return &sliceIterator[string]{slice: f.keys}, nil
}

func (f *field) Descriptors() []protoreflect.FieldDescriptor {
	return f.descriptors
}

func (f *field) add(k string) {
	f.keys = append(f.keys, k)
}

func (f *field) remove(k string) {
	for i, v := range f.keys {
		if k == v {
			f.keys = append(f.keys[:i], f.keys[i+1:]...)
			return
		}
	}
}
