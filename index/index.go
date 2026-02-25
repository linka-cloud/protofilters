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
	"reflect"
	"strconv"
	"strings"

	"github.com/cespare/xxhash/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"go.linka.cloud/protofilters/filters"
	"go.linka.cloud/protofilters/index/bitmap"
	preflect "go.linka.cloud/protofilters/reflect"
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
	// Update updates the index for the given key using old and new messages
	Update(ctx context.Context, k string, old, m proto.Message) error
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
		if fd.Kind() == protoreflect.MessageKind && !preflect.IsWKType(fd.Message().FullName()) {
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

func (i *index) Update(ctx context.Context, k string, old, m proto.Message) error {
	tx, err := i.store.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close()
	oldValues := map[string]fieldValues{}
	newValues := map[string]fieldValues{}
	if old != nil {
		if oldValues, err = i.collectValues(ctx, old.ProtoReflect()); err != nil {
			return err
		}
	}
	if m != nil {
		if newValues, err = i.collectValues(ctx, m.ProtoReflect()); err != nil {
			return err
		}
	}
	if len(oldValues) > 0 {
		fr, err := tx.For(ctx, old.ProtoReflect().Descriptor().FullName())
		if err != nil {
			return err
		}
		if err := applyDiff(ctx, fr, tx, k, oldValues, newValues); err != nil {
			return err
		}
	} else {
		for _, fv := range newValues {
			for _, v := range fv.values {
				if err := tx.Add(ctx, k, v, fv.fds...); err != nil {
					return err
				}
			}
		}
	}
	return tx.Commit(ctx)
}

type fieldValues struct {
	fds    []protoreflect.FieldDescriptor
	values []protoreflect.Value
}

func (i *index) collectValues(ctx context.Context, m protoreflect.Message) (map[string]fieldValues, error) {
	values := map[string]fieldValues{}
	if err := i.collectValuesInto(ctx, values, m); err != nil {
		return nil, err
	}
	return values, nil
}

func (i *index) collectValuesInto(ctx context.Context, out map[string]fieldValues, m protoreflect.Message, fds ...protoreflect.FieldDescriptor) error {
	f := m.Descriptor().Fields()
	name := m.Descriptor().FullName()
	for j := 0; j < f.Len(); j++ {
		fd := f.Get(j)
		path := append(fds, fd)
		ok, err := i.fn(ctx, name, path...)
		if err != nil {
			return err
		}
		rval := m.Get(fd)
		if fd.IsList() {
			if fd.Kind() == protoreflect.MessageKind {
				for j2 := 0; j2 < rval.List().Len(); j2++ {
					if err := i.collectValuesInto(ctx, out, rval.List().Get(j2).Message(), path...); err != nil {
						return err
					}
				}
			}
			list := rval.List()
			for j2 := 0; j2 < list.Len(); j2++ {
				out = appendValue(out, path, list.Get(j2))
			}
			continue
		}
		if fd.IsMap() {
			continue
		}
		if fd.Kind() == protoreflect.MessageKind && !preflect.IsWKType(fd.Message().FullName()) {
			if !rval.Message().IsValid() {
				continue
			}
			if err := i.collectValuesInto(ctx, out, rval.Message(), path...); err != nil {
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
		out = appendValue(out, path, rval)
	}
	return nil
}

func appendValue(out map[string]fieldValues, fds []protoreflect.FieldDescriptor, v protoreflect.Value) map[string]fieldValues {
	key := string(joinFieldNames(fds))
	fv := out[key]
	if fv.fds == nil {
		fv.fds = append([]protoreflect.FieldDescriptor(nil), fds...)
	}
	for _, existing := range fv.values {
		if valueEqual(existing, v) {
			out[key] = fv
			return out
		}
	}
	fv.values = append(fv.values, v)
	out[key] = fv
	return out
}

func applyDiff(ctx context.Context, fr FieldReader, tx Tx, k string, oldValues, newValues map[string]fieldValues) error {
	seen := map[string]struct{}{}
	for key := range oldValues {
		seen[key] = struct{}{}
	}
	for key := range newValues {
		seen[key] = struct{}{}
	}
	for key := range seen {
		ov := oldValues[key]
		nv := newValues[key]
		remove, add := diffValues(ov.values, nv.values)
		for _, v := range remove {
			if err := removeValue(ctx, fr, k, v, ov.fds...); err != nil {
				return err
			}
		}
		for _, v := range add {
			if err := tx.Add(ctx, k, v, nv.fds...); err != nil {
				return err
			}
		}
	}
	return nil
}

func diffValues(oldValues, newValues []protoreflect.Value) (remove []protoreflect.Value, add []protoreflect.Value) {
	used := make([]bool, len(newValues))
	for _, ov := range oldValues {
		found := false
		for i, nv := range newValues {
			if used[i] {
				continue
			}
			if valueEqual(ov, nv) {
				used[i] = true
				found = true
				break
			}
		}
		if !found {
			remove = append(remove, ov)
		}
	}
	for i, nv := range newValues {
		if !used[i] {
			add = append(add, nv)
		}
	}
	return remove, add
}

func removeValue(ctx context.Context, fr FieldReader, k string, v protoreflect.Value, fds ...protoreflect.FieldDescriptor) error {
	if len(fds) == 0 {
		return nil
	}
	name := joinFieldNames(fds)
	h := keyHash(k)
	for f, err := range fr.Get(ctx, name) {
		if err != nil {
			return err
		}
		if !valueEqual(f.Value(), v) {
			continue
		}
		b, err := f.Bitmap(ctx)
		if err != nil {
			return err
		}
		b.Remove(h)
		return nil
	}
	return nil
}

func joinFieldNames(fds []protoreflect.FieldDescriptor) protoreflect.Name {
	if len(fds) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(string(fds[0].Name()))
	for _, fd := range fds[1:] {
		b.WriteByte('.')
		b.WriteString(string(fd.Name()))
	}
	return protoreflect.Name(b.String())
}

func valueEqual(a, b protoreflect.Value) bool {
	av := a.Interface()
	bv := b.Interface()
	if av == nil || bv == nil {
		return av == bv
	}
	if reflect.TypeOf(av) != reflect.TypeOf(bv) {
		return false
	}
	if reflect.TypeOf(av).Comparable() {
		return av == bv
	}
	return reflect.DeepEqual(av, bv)
}

func keyHash(k string) uint64 {
	if i, err := strconv.ParseUint(k, 10, 64); err == nil {
		return i
	}
	return xxhash.Sum64String(k)
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
		ok, err := preflect.Match(v.Value(), fd, f.Filter)
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
