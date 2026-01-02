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

package set

import (
	"encoding/binary"
	"iter"

	"github.com/tidwall/btree"

	bitmap2 "go.linka.cloud/protofilters/index/bitmap"
)

var (
	_ bitmap2.Provider = (*prov)(nil)
	_ bitmap2.Bitmap   = (*bitmap)(nil)
)

type prov struct{}

func (prov) New() bitmap2.Bitmap {
	return &bitmap{
		s: &btree.Set[uint64]{},
	}
}

func (prov) NewWith(_ int) bitmap2.Bitmap {
	return &bitmap{
		s: &btree.Set[uint64]{},
	}
}

func (prov) NewFrom(buf []byte) bitmap2.Bitmap {
	s := &btree.Set[uint64]{}
	for i := 0; i < len(buf); i += 8 {
		s.Insert(binary.LittleEndian.Uint64(buf[i:]))
	}
	return &bitmap{
		s: s,
	}
}

type bitmap struct {
	s *btree.Set[uint64]
}

func (b *bitmap) Set(k uint64) {
	b.s.Insert(k)
}

func (b *bitmap) Remove(k uint64) {
	b.s.Delete(k)
}

func (b *bitmap) And(other bitmap2.Bitmap) {
	o := other.(*bitmap)
	it := b.s.Iter()
	for it.Next() {
		if !o.s.Contains(it.Key()) {
			b.s.Delete(it.Key())
		}
	}
}

func (b *bitmap) Or(other bitmap2.Bitmap) {
	o := other.(*bitmap)
	for it := o.s.Iter(); it.Next(); {
		b.s.Insert(it.Key())
	}
}

func (b *bitmap) Bytes() []byte {
	buf := make([]byte, 8*b.s.Len())
	i := 0
	for it := b.s.Iter(); it.Next(); {
		binary.LittleEndian.PutUint64(buf[i:], it.Key())
		i += 8
	}
	return buf
}

func (b *bitmap) Iter() iter.Seq[uint64] {
	return func(yield func(uint64) bool) {
		it := b.s.Iter()
		for it.Next() {
			if !yield(it.Key()) {
				return
			}
		}
	}
}

func init() {
	bitmap2.SetProvider(prov{})
}
