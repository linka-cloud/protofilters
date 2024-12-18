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

	"github.com/cespare/xxhash/v2"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Field is an indexed field in the index
type Field interface {
	// Value returns the value of the field
	Value() protoreflect.Value
	// Bitmap returns a Bitmap of the keys that have this value
	Bitmap(ctx context.Context) (*Bitmap, error)
	// Descriptors returns the field descriptors for this field
	Descriptors() []protoreflect.FieldDescriptor
}

// FieldReader is an interface for reading a type fields from the index
type FieldReader interface {
	// Get returns the field for the given field descriptor
	Get(ctx context.Context, f protoreflect.Name) (Iterator[Field], bool, error)
}

func newField(v protoreflect.Value, fds []protoreflect.FieldDescriptor) *field {
	return &field{
		value:       v,
		bitmap:      NewBitmapWith(1024),
		descriptors: fds,
	}
}

type field struct {
	value       protoreflect.Value
	bitmap      *Bitmap
	descriptors []protoreflect.FieldDescriptor
}

func (f *field) Value() protoreflect.Value {
	return f.value
}

func (f *field) Bitmap(_ context.Context) (*Bitmap, error) {
	return f.bitmap, nil
}

func (f *field) Descriptors() []protoreflect.FieldDescriptor {
	return f.descriptors
}

func (f *field) add(k string) uint64 {
	h := xxhash.Sum64String(k)
	f.bitmap.Set(h)
	return h
}

func (f *field) remove(k string) {
	h := xxhash.Sum64String(k)
	f.bitmap.Remove(h)
}
