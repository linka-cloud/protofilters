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
	"iter"
	"strconv"
	"strings"
	"sync"

	"github.com/cespare/xxhash/v2"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Txer is an interface for a transactioner.
type Txer interface {
	// Tx returns a transaction.
	Tx(ctx context.Context) (Tx, error)
}

// Tx is an interface for a transaction.
type Tx interface {
	Store
	// Commit commits the transaction.
	Commit(ctx context.Context) error
	// Close closes the transaction.
	// If the transaction has not been committed, it will be rolled back.
	// If the transaction has been committed, it should be a no-op.
	Close() error
}

// UIDTxer is an interface for UID transactions.
type UIDTxer interface {
	Tx(ctx context.Context) (UIDTx, error)
}

// UIDTx is an interface for a UID transaction.
type UIDTx interface {
	UIDStore
	Commit(ctx context.Context) error
	Close() error
}

// Store is an interface for storing and retrieving protobuf message fields.
type Store interface {
	// For returns a FieldReader for the given message type.
	For(ctx context.Context, t protoreflect.FullName) (FieldReader, error)
	// Add adds a value to the store for the given key and field descriptor.
	Add(ctx context.Context, k string, v protoreflect.Value, fds ...protoreflect.FieldDescriptor) error
	// Remove removes a value from the store for the given key and field descriptor.
	Remove(ctx context.Context, k string, f protoreflect.FieldDescriptor, v protoreflect.Value) error
	// Keys returns the values for the given uint64 hash.
	// It may return multiple values in case of hash collisions.
	Keys(ctx context.Context, i uint64) ([]string, error)
	// Clear removes all values from the store for the given key.
	Clear(ctx context.Context, k string) error
}

// UIDStore is an interface for storing and retrieving protobuf message fields by UID.
type UIDStore interface {
	// For returns a FieldReader for the given message type.
	For(ctx context.Context, t protoreflect.FullName) (FieldReader, error)
	// AddUID adds a value to the store for the given UID and field descriptor.
	AddUID(ctx context.Context, uid uint64, v protoreflect.Value, fds ...protoreflect.FieldDescriptor) error
	// RemoveUID removes a value from the store for the given UID and field descriptor.
	RemoveUID(ctx context.Context, uid uint64, f protoreflect.FieldDescriptor, v protoreflect.Value) error
	// ClearUID removes all values from the store for the given UID.
	ClearUID(ctx context.Context, uid uint64) error
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

type fakeUIDTxer struct {
	UIDStore
}

func (f fakeUIDTxer) Tx(_ context.Context) (UIDTx, error) {
	return noopUIDTx{UIDStore: f.UIDStore}, nil
}

type noopUIDTx struct {
	UIDStore
}

func (n noopUIDTx) Commit(_ context.Context) error {
	return nil
}

func (n noopUIDTx) Close() error {
	return nil
}

type fieldReader struct {
	m map[protoreflect.Name][]*field
}

func (f *fieldReader) Get(_ context.Context, n protoreflect.Name) iter.Seq2[Field, error] {
	return func(yield func(Field, error) bool) {
		for _, v := range f.m[n] {
			if !yield(v, nil) {
				return
			}
		}
	}
}

func newStore() Store {
	return &store{
		fields:   make(map[protoreflect.FullName][]*field),
		hashKeys: make(map[uint64][]string),
		keyHash:  make(map[string]uint64),
	}
}

func newUIDStore() UIDStore {
	return &uidStore{
		fields: make(map[protoreflect.FullName][]*field),
	}
}

type store struct {
	fields   map[protoreflect.FullName][]*field
	hashKeys map[uint64][]string
	keyHash  map[string]uint64
	m        sync.RWMutex
}

type uidStore struct {
	fields map[protoreflect.FullName][]*field
	m      sync.RWMutex
}

type uidTxer struct {
	Txer
}

func (l uidTxer) Tx(ctx context.Context) (UIDTx, error) {
	tx, err := l.Txer.Tx(ctx)
	if err != nil {
		return nil, err
	}
	return uidTx{Tx: tx}, nil
}

type uidTx struct {
	Tx
}

func (l uidTx) AddUID(ctx context.Context, uid uint64, v protoreflect.Value, fds ...protoreflect.FieldDescriptor) error {
	return l.Tx.Add(ctx, strconv.FormatUint(uid, 10), v, fds...)
}

func (l uidTx) RemoveUID(ctx context.Context, uid uint64, f protoreflect.FieldDescriptor, v protoreflect.Value) error {
	return l.Tx.Remove(ctx, strconv.FormatUint(uid, 10), f, v)
}

func (l uidTx) ClearUID(ctx context.Context, uid uint64) error {
	return l.Tx.Clear(ctx, strconv.FormatUint(uid, 10))
}

func (s *store) Keys(_ context.Context, i uint64) ([]string, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	k := s.hashKeys[i]
	o := make([]string, len(k))
	copy(o, k)
	return o, nil
}

func (s *store) For(_ context.Context, t protoreflect.FullName) (FieldReader, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	out := make(map[protoreflect.Name][]*field)
	for k, v := range s.fields {
		if strings.HasPrefix(string(k), string(t)) {
			// +1 for the dot
			out[protoreflect.Name(k[len(t)+1:])] = v
		}
	}
	return &fieldReader{m: out}, nil
}

func (s *store) Add(_ context.Context, k string, v protoreflect.Value, fds ...protoreflect.FieldDescriptor) error {
	s.m.Lock()
	defer s.m.Unlock()
	if len(fds) == 0 {
		return nil
	}
	n := fds[0].FullName()
	for _, v := range fds[1:] {
		n = n.Append(v.Name())
	}
	if _, ok := s.fields[n]; !ok {
		s.fields[n] = make([]*field, 0)
	}
	for _, fi := range s.fields[n] {
		if fi.value.Interface() == v.Interface() {
			i := fi.add(k)
			s.addIndex(k, i)
			return nil
		}
	}
	fi := newField(v, fds)
	i := fi.add(k)
	s.addIndex(k, i)
	s.fields[n] = append(s.fields[n], fi)
	return nil
}

func (s *store) addIndex(k string, i uint64) {
	if _, ok := s.keyHash[k]; ok {
		return
	}
	s.hashKeys[i] = append(s.hashKeys[i], k)
	s.keyHash[k] = i
}

func (s *store) Remove(_ context.Context, k string, f protoreflect.FieldDescriptor, v protoreflect.Value) error {
	s.m.Lock()
	defer s.m.Unlock()
	if _, ok := s.fields[f.FullName()]; !ok {
		return nil
	}
	if _, ok := s.keyHash[k]; !ok {
		return nil
	}
	for _, fi := range s.fields[f.FullName()] {
		if fi.value.Interface() == v.Interface() {
			fi.remove(k)
			return nil
		}
	}
	return nil
}

func (s *store) Clear(_ context.Context, k string) error {
	s.m.Lock()
	defer s.m.Unlock()
	i := xxhash.Sum64String(k)
	delete(s.hashKeys, i)
	delete(s.keyHash, k)
	for _, fis := range s.fields {
		for _, fi := range fis {
			fi.remove(k)
		}
	}
	return nil
}

func (s *uidStore) For(_ context.Context, t protoreflect.FullName) (FieldReader, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	out := make(map[protoreflect.Name][]*field)
	for k, v := range s.fields {
		if strings.HasPrefix(string(k), string(t)) {
			out[protoreflect.Name(k[len(t)+1:])] = v
		}
	}
	return &fieldReader{m: out}, nil
}

func (s *uidStore) AddUID(_ context.Context, uid uint64, v protoreflect.Value, fds ...protoreflect.FieldDescriptor) error {
	s.m.Lock()
	defer s.m.Unlock()
	if len(fds) == 0 {
		return nil
	}
	n := fds[0].FullName()
	for _, v := range fds[1:] {
		n = n.Append(v.Name())
	}
	if _, ok := s.fields[n]; !ok {
		s.fields[n] = make([]*field, 0)
	}
	for _, fi := range s.fields[n] {
		if fi.value.Interface() == v.Interface() {
			fi.addUID(uid)
			return nil
		}
	}
	fi := newField(v, fds)
	fi.addUID(uid)
	s.fields[n] = append(s.fields[n], fi)
	return nil
}

func (s *uidStore) RemoveUID(_ context.Context, uid uint64, f protoreflect.FieldDescriptor, v protoreflect.Value) error {
	s.m.Lock()
	defer s.m.Unlock()
	if _, ok := s.fields[f.FullName()]; !ok {
		return nil
	}
	for _, fi := range s.fields[f.FullName()] {
		if fi.value.Interface() == v.Interface() {
			fi.removeUID(uid)
			return nil
		}
	}
	return nil
}

func (s *uidStore) ClearUID(_ context.Context, uid uint64) error {
	s.m.Lock()
	defer s.m.Unlock()
	for _, fis := range s.fields {
		for _, fi := range fis {
			fi.removeUID(uid)
		}
	}
	return nil
}
