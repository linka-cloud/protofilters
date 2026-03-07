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
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/cespare/xxhash/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"go.linka.cloud/protofilters/filters"
)

func All(_ context.Context, _ protoreflect.FullName, _ ...protoreflect.FieldDescriptor) (bool, error) {
	return true, nil
}

// Func is a function that is called to determine if a field should be indexed.
type Func func(ctx context.Context, name protoreflect.FullName, fds ...protoreflect.FieldDescriptor) (bool, error)

// Index is a protobuf message index.
type Index interface {
	Insert(ctx context.Context, k string, m proto.Message) error
	Update(ctx context.Context, k string, old, m proto.Message) error
	Remove(ctx context.Context, k string) error
	Find(ctx context.Context, t protoreflect.FullName, f filters.FieldFilterer) ([]string, []string, error)
}

// New creates a compatibility key-based index backed by the UID index implementation.
func New(s Store, fn Func) Index {
	if fn == nil {
		fn = All
	}
	if s == nil {
		return &keyIndex{
			uid:      NewUID(nil, fn),
			resolver: newUIDKeys(),
		}
	}
	x, ok := any(s).(Txer)
	if !ok {
		x = &fakeTxer{Store: s}
	}
	return &keyIndex{
		uid:      newUIDFromTxer(uidTxer{Txer: x}, fn),
		store:    x,
		resolver: newUIDKeys(),
	}
}

type keyIndex struct {
	uid      UIDIndex
	store    Txer
	resolver *uidKeys
}

type uidKeys struct {
	mu      sync.RWMutex
	uidKeys map[uint64]map[string]struct{}
	keyUID  map[string]uint64
}

func newUIDKeys() *uidKeys {
	return &uidKeys{
		uidKeys: make(map[uint64]map[string]struct{}),
		keyUID:  make(map[string]uint64),
	}
}

func (r *uidKeys) add(key string, uid uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.keyUID[key] = uid
	if _, ok := r.uidKeys[uid]; !ok {
		r.uidKeys[uid] = make(map[string]struct{})
	}
	r.uidKeys[uid][key] = struct{}{}
}

func (r *uidKeys) remove(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	uid, ok := r.keyUID[key]
	if !ok {
		return
	}
	delete(r.keyUID, key)
	keys := r.uidKeys[uid]
	delete(keys, key)
	if len(keys) == 0 {
		delete(r.uidKeys, uid)
	}
}

func (r *uidKeys) keys(uid uint64) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	keys := r.uidKeys[uid]
	out := make([]string, 0, len(keys))
	for k := range keys {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func (i *keyIndex) Insert(ctx context.Context, k string, m proto.Message) error {
	uid := keyHash(k)
	if err := i.uid.Insert(ctx, uid, m); err != nil {
		return err
	}
	i.resolver.add(k, uid)
	return nil
}

func (i *keyIndex) Update(ctx context.Context, k string, old, m proto.Message) error {
	uid := keyHash(k)
	if err := i.uid.Update(ctx, uid, old, m); err != nil {
		return err
	}
	i.resolver.add(k, uid)
	return nil
}

func (i *keyIndex) Remove(ctx context.Context, k string) error {
	uid := keyHash(k)
	if err := i.uid.Remove(ctx, uid); err != nil {
		return err
	}
	i.resolver.remove(k)
	return nil
}

func (i *keyIndex) Find(ctx context.Context, t protoreflect.FullName, f filters.FieldFilterer) ([]string, []string, error) {
	var tx Tx
	if i.store != nil {
		var err error
		tx, err = i.store.Tx(ctx)
		if err != nil {
			return nil, nil, err
		}
		defer tx.Close()
	}

	var keys []string
	var collisions []string
	for uid, err := range i.uid.Find(ctx, t, f, FindOptions{}) {
		if err != nil {
			return nil, nil, err
		}
		single, many := i.resolver.oneOrMany(uid)
		if tx == nil {
			if many != nil {
				collisions = append(collisions, many...)
				continue
			}
			if single != "" {
				keys = append(keys, single)
			}
			continue
		}

		extra, err := tx.Keys(ctx, uid)
		if err != nil {
			return nil, nil, err
		}
		resolved := mergeResolvedKeys(uid, single, many, extra)
		if len(resolved) != 1 {
			collisions = append(collisions, resolved...)
			continue
		}
		keys = append(keys, resolved...)
	}
	return keys, collisions, nil
}

func mergeResolvedKeys(uid uint64, single string, many, extra []string) []string {
	if len(extra) == 0 {
		if many != nil {
			return many
		}
		if single == "" {
			return nil
		}
		return []string{single}
	}

	m := make(map[string]struct{}, 1+len(many)+len(extra))
	if single != "" {
		m[single] = struct{}{}
	}
	for _, v := range many {
		m[v] = struct{}{}
	}
	uidPlaceholder := strconv.FormatUint(uid, 10)
	for _, v := range extra {
		if v == uidPlaceholder {
			continue
		}
		m[v] = struct{}{}
	}
	out := make([]string, 0, len(m))
	for v := range m {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

func (r *uidKeys) oneOrMany(uid uint64) (string, []string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	keys := r.uidKeys[uid]
	if len(keys) == 0 {
		return "", nil
	}
	if len(keys) == 1 {
		for k := range keys {
			return k, nil
		}
	}
	out := make([]string, 0, len(keys))
	for k := range keys {
		out = append(out, k)
	}
	sort.Strings(out)
	return "", out
}

type fieldValues struct {
	fds    []protoreflect.FieldDescriptor
	values []protoreflect.Value
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

func isUnsetRealOneofField(m protoreflect.Message, fd protoreflect.FieldDescriptor) bool {
	oneof := fd.ContainingOneof()
	if oneof == nil || oneof.IsSynthetic() {
		return false
	}
	return !m.Has(fd)
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
