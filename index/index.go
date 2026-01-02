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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"go.linka.cloud/protofilters/filters"
	"go.linka.cloud/protofilters/index/bitmap"
	"go.linka.cloud/protofilters/reflect"
)

func All(_ context.Context, _ protoreflect.FullName, _ ...protoreflect.FieldDescriptor) (bool, error) {
	return true, nil
}

// Func is a function that is called to determine if a field should be indexed
// It takes a context and a list of field descriptors that represent the path to the field
// It returns a bool indicating if the field should be indexed or not and an error if any
type Func func(ctx context.Context, name protoreflect.FullName, fds ...protoreflect.FieldDescriptor) (bool, error)

// Index is a protobuf message index
type Index interface {
	// Insert inserts and indexes the given message with the given key
	Insert(ctx context.Context, k string, m proto.Message) error
	// Update updates the index for the given message with the given key
	Update(ctx context.Context, k string, m proto.Message) error
	// Remove removes the given key from the index
	Remove(ctx context.Context, k string) error
	// Find returns the keys of the messages that match the given FieldFilterer
	// It returns the keys that match the filter and the keys that have hash collisions
	// They need to be resolved by checking the actual values
	Find(ctx context.Context, t protoreflect.FullName, f filters.FieldFilterer) ([]string, []string, error)
}

type index struct {
	store Txer
	fn    Func
}

// New creates a new index using the given store and index function
// If the store is nil, a new in-memory store is created
// If the index function is nil, all fields are indexed
func New(s Store, fn Func) Index {
	if fn == nil {
		fn = All
	}
	if s == nil {
		s = newStore()
	}
	x, ok := any(s).(Txer)
	if !ok {
		x = &fakeTxer{Store: s}
	}
	return &index{
		store: x,
		fn:    fn,
	}
}

func (i *index) index(ctx context.Context, tx Tx, k string, m protoreflect.Message, fds ...protoreflect.FieldDescriptor) error {
	f := m.Descriptor().Fields()
	name := m.Descriptor().FullName()
	if len(fds) > 0 {
		_ = fds
	}
	for j := 0; j < f.Len(); j++ {
		fd := f.Get(j)
		fds := append(fds, fd)
		ok, err := i.fn(ctx, name, fds...)
		if err != nil {
			return err
		}
		rval := m.Get(fd)
		if fd.IsList() {
			// we don't index lists of messages
			if fd.Kind() == protoreflect.MessageKind {
				for j2 := 0; j2 < rval.List().Len(); j2++ {
					m := rval.List().Get(j2).Message()
					if err := i.index(ctx, tx, k, m, fds...); err != nil {
						return err
					}
				}
			}
			list := rval.List()
			for j2 := 0; j2 < list.Len(); j2++ {
				if err := tx.Add(ctx, k, list.Get(j2), fds...); err != nil {
					return err
				}
			}
			continue
		}
		if fd.IsMap() {
			continue
		}
		if fd.Kind() == protoreflect.MessageKind && !reflect.IsWKType(fd.Message().FullName()) {
			if !rval.Message().IsValid() {
				continue
			}
			if err := i.index(ctx, tx, k, rval.Message(), fds...); err != nil {
				return err
			}
			continue
		}
		if fd.HasOptionalKeyword() && !m.Has(fd) {
			rval = protoreflect.Value{}
		}
		if !ok {
			continue
		}
		if err := tx.Add(ctx, k, rval, fds...); err != nil {
			return err
		}
	}
	return nil
}

// Insert indexes a message fields with the given key.
func (i *index) Insert(ctx context.Context, k string, m proto.Message) error {
	tx, err := i.store.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close()
	if err := i.index(ctx, tx, k, m.ProtoReflect()); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (i *index) Update(ctx context.Context, k string, m proto.Message) error {
	tx, err := i.store.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close()
	f := m.ProtoReflect().Descriptor().Fields()
	for j := 0; j < f.Len(); j++ {
		fd := f.Get(j)
		if !fd.IsList() {
			if err := tx.Remove(ctx, k, fd, m.ProtoReflect().Get(fd)); err != nil {
				return err
			}
		}
	}
	if err := i.index(ctx, tx, k, m.ProtoReflect()); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (i *index) Remove(ctx context.Context, k string) error {
	tx, err := i.store.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close()
	if err := tx.Clear(ctx, k); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (i *index) doFind(ctx context.Context, tx Tx, t protoreflect.FullName, f *filters.FieldFilter) (bitmap.Bitmap, error) {
	fds, err := tx.For(ctx, t)
	if err != nil {
		return nil, err
	}

	b := bitmap.NewWith(1024)
	for v, err := range fds.Get(ctx, protoreflect.Name(f.Field)) {
		if err != nil {
			return nil, err
		}
		ds := v.Descriptors()
		fd := ds[len(ds)-1]
		ok, err := reflect.Match(v.Value(), fd, f.Filter)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		b2, err := v.Bitmap(ctx)
		if err != nil {
			return nil, err
		}
		b.Or(b2)
	}
	return b, nil
}

func (i *index) find(ctx context.Context, tx Tx, t protoreflect.FullName, f filters.FieldFilterer) (bitmap.Bitmap, error) {
	expr := f.Expr()
	b, err := i.doFind(ctx, tx, t, expr.Condition)
	if err != nil {
		return nil, err
	}
	for _, v := range expr.AndExprs {
		b2, err := i.find(ctx, tx, t, v)
		if err != nil {
			return nil, err
		}
		b.And(b2)
	}
	for _, v := range expr.OrExprs {
		b2, err := i.find(ctx, tx, t, v)
		if err != nil {
			return nil, err
		}
		b.Or(b2)
	}
	return b, nil
}

func (i *index) Find(ctx context.Context, t protoreflect.FullName, f filters.FieldFilterer) ([]string, []string, error) {
	if f == nil || f.Expr() == nil {
		return nil, nil, nil
	}
	tx, err := i.store.Tx(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Close()
	b, err := i.find(ctx, tx, t, f)
	if err != nil {
		return nil, nil, err
	}
	var keys []string
	var collisions []string
	for v := range b.Iter() {
		ks, err := tx.Keys(ctx, v)
		if err != nil {
			return nil, nil, err
		}
		if len(ks) != 1 {
			collisions = append(collisions, ks...)
			continue
		}
		keys = append(keys, ks...)
	}
	return keys, collisions, nil
}
