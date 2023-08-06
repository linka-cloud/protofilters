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
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type Field struct {
	Value       protoreflect.Value
	Keys        []string
	Descriptors []protoreflect.FieldDescriptor
}

func (f *Field) add(k string) {
	f.Keys = append(f.Keys, k)
}

func (f *Field) remove(k string) {
	for i, v := range f.Keys {
		if k == v {
			f.Keys = append(f.Keys[:i], f.Keys[i+1:]...)
			return
		}
	}
}

type Txer interface {
	Tx(ctx context.Context) (Tx, error)
}

type Tx interface {
	Store
	Commit(ctx context.Context) error
	Close() error
}

type FieldReader interface {
	Get(ctx context.Context, f protoreflect.Name) ([]*Field, bool, error)
}

type Store interface {
	For(ctx context.Context, t protoreflect.FullName) (FieldReader, error)
	Add(ctx context.Context, k string, v protoreflect.Value, fds ...protoreflect.FieldDescriptor) error
	Remove(ctx context.Context, k string, f protoreflect.FieldDescriptor, v protoreflect.Value) error
	Clear(ctx context.Context, k string) error
}

type fakeTxer struct {
	Store
}

func (f fakeTxer) Tx(_ context.Context) (Tx, error) {
	return noopTX{Store: f.Store}, nil
}

type noopTX struct {
	Store
}

func (n noopTX) Commit(_ context.Context) error {
	return nil
}

func (n noopTX) Close() error {
	return nil
}

type fieldReader struct {
	m map[protoreflect.Name][]*Field
}

func (f *fieldReader) Get(_ context.Context, n protoreflect.Name) ([]*Field, bool, error) {
	out, ok := f.m[n]
	return out, ok, nil
}

type store map[protoreflect.FullName][]*Field

func (s store) For(_ context.Context, t protoreflect.FullName) (FieldReader, error) {
	out := make(map[protoreflect.Name][]*Field)
	for k, v := range s {
		if strings.HasPrefix(string(k), string(t)) {
			// +1 for the dot
			out[protoreflect.Name(k[len(t)+1:])] = v
		}
	}
	return &fieldReader{m: out}, nil
}

func (s store) Add(_ context.Context, k string, v protoreflect.Value, fds ...protoreflect.FieldDescriptor) error {
	if len(fds) == 0 {
		return nil
	}
	n := fds[0].FullName()
	for _, v := range fds[1:] {
		n = n.Append(v.Name())
	}
	if _, ok := s[n]; !ok {
		s[n] = make([]*Field, 0)
	}
	for _, fi := range s[n] {
		if fi.Value.Interface() == v.Interface() {
			fi.add(k)
			return nil
		}
	}
	s[n] = append(s[n], &Field{
		Value:       v,
		Keys:        []string{k},
		Descriptors: fds,
	})
	return nil
}

func (s store) Remove(_ context.Context, k string, f protoreflect.FieldDescriptor, v protoreflect.Value) error {
	if _, ok := s[f.FullName()]; !ok {
		return nil
	}
	for _, fi := range s[f.FullName()] {
		if fi.Value.Interface() == v.Interface() {
			fi.remove(k)
			return nil
		}
	}
	return nil
}

func (s store) Clear(_ context.Context, k string) error {
	for _, fis := range s {
		for _, fi := range fis {
			fi.remove(k)
		}
	}
	return nil
}
