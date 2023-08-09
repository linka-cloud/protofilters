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

	"github.com/dgraph-io/sroar"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"go.linka.cloud/protofilters/filters"
	"go.linka.cloud/protofilters/reflect"
)

func All(_ context.Context, _ ...protoreflect.FieldDescriptor) (bool, error) {
	return true, nil
}

// Func is a function that is called to determine if a field should be indexed
// It takes a context and a list of field descriptors that represent the path to the field
// It returns a bool indicating if the field should be indexed or not and an error if any
type Func func(ctx context.Context, fds ...protoreflect.FieldDescriptor) (bool, error)

// Index is a protobuf message index
type Index interface {
	// Insert inserts and indexes the given message with the given key
	Insert(ctx context.Context, k string, m proto.Message) error
	// Update updates the index for the given message with the given key
	Update(ctx context.Context, k string, m proto.Message) error
	// Remove removes the given key from the index
	Remove(ctx context.Context, k string) error
	// Find returns the keys of the messages that match the given FieldFilterer
	Find(ctx context.Context, t protoreflect.FullName, f filters.FieldFilterer) ([]string, error)
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
		s = make(store)
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

func (i *index) index(ctx context.Context, tx Tx, k string, m proto.Message, fds ...protoreflect.FieldDescriptor) error {
	f := m.ProtoReflect().Descriptor().Fields()
	for j := 0; j < f.Len(); j++ {
		fd := f.Get(j)
		fds := append(fds, fd)
		ok, err := i.fn(ctx, fds...)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		rval := m.ProtoReflect().Get(fd)
		if fd.IsList() {
			// we don't index lists of messages
			if fd.Kind() == protoreflect.MessageKind {
				continue
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
		if fd.Kind() == protoreflect.MessageKind {
			if !rval.Message().IsValid() {
				continue
			}
			m := rval.Message().Interface()
			if err := i.index(ctx, tx, k, m.(proto.Message), fds...); err != nil {
				return err
			}
			continue
		}
		if fd.HasOptionalKeyword() && !m.ProtoReflect().Has(fd) {
			rval = protoreflect.Value{}
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
	if err := i.index(ctx, tx, k, m); err != nil {
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
	if err := i.index(ctx, tx, k, m); err != nil {
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

func (i *index) doFind(ctx context.Context, tx Tx, r *keyReg, t protoreflect.FullName, f *filters.FieldFilter) (*sroar.Bitmap, error) {
	fds, err := tx.For(ctx, t)
	if err != nil {
		return nil, err
	}

	fit, ok, err := fds.Get(ctx, protoreflect.Name(f.Field))
	if !ok {
		return nil, nil
	}
	b := sroar.NewBitmap()
	for fit.Next() {
		v, err := fit.Value()
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
		it, err := v.Keys(ctx)
		if err != nil {
			return nil, err
		}
		for it.Next() {
			v, err := it.Value()
			if err != nil {
				return nil, err
			}
			b.Set(r.index(v))
		}
	}
	return b, nil
}

func (i *index) find(ctx context.Context, tx Tx, r *keyReg, t protoreflect.FullName, f filters.FieldFilterer) (*sroar.Bitmap, error) {
	expr := f.Expr()
	b, err := i.doFind(ctx, tx, r, t, expr.Condition)
	if err != nil {
		return nil, err
	}
	for _, v := range expr.AndExprs {
		b2, err := i.find(ctx, tx, r, t, v)
		if err != nil {
			return nil, err
		}
		b.And(b2)
	}
	for _, v := range expr.OrExprs {
		b2, err := i.find(ctx, tx, r, t, v)
		if err != nil {
			return nil, err
		}
		b.Or(b2)
	}
	return b, nil
}

func (i *index) Find(ctx context.Context, t protoreflect.FullName, f filters.FieldFilterer) ([]string, error) {
	if f == nil || f.Expr() == nil {
		return nil, nil
	}
	tx, err := i.store.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Close()
	r := newKeyReg()
	b, err := i.find(ctx, tx, r, t, f)
	if err != nil {
		return nil, err
	}
	return r.keysFor(b.ToArray()), nil
}

type keyReg struct {
	keys []string
}

func newKeyReg() *keyReg {
	return &keyReg{
		// size one to skip index 0
		keys: make([]string, 1),
	}
}

func (r *keyReg) index(k string) uint64 {
	for i, v := range r.keys {
		if v == k {
			return uint64(i)
		}
	}
	r.keys = append(r.keys, k)
	return uint64(len(r.keys) - 1)
}

func (r *keyReg) key(i uint64) string {
	return r.keys[i]
}

func (r *keyReg) keysFor(i []uint64) []string {
	out := make([]string, len(i))
	for j, v := range i {
		out[j] = r.key(v)
	}
	return out
}
