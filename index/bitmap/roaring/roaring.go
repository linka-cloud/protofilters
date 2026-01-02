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

package roaring

import (
	"encoding/binary"
	"iter"

	"github.com/RoaringBitmap/roaring/v2/roaring64"

	bitmap2 "go.linka.cloud/protofilters/index/bitmap"
)

var (
	_ bitmap2.Provider = (*prov)(nil)
	_ bitmap2.Bitmap   = (*bitmap)(nil)
)

type prov struct{}

func (prov) New() bitmap2.Bitmap {
	return &bitmap{
		m: roaring64.New(),
	}
}

func (prov) NewWith(_ int) bitmap2.Bitmap {
	return &bitmap{
		m: roaring64.New(),
	}
}

func (prov) NewFrom(buf []byte) bitmap2.Bitmap {
	m := roaring64.New()
	for i := 0; i < len(buf); i += 8 {
		m.Add(binary.LittleEndian.Uint64(buf[i:]))
	}
	return &bitmap{m: m}
}

type bitmap struct {
	m *roaring64.Bitmap
}

func (r *bitmap) Set(k uint64) {
	r.m.Add(k)
}

func (r *bitmap) Remove(k uint64) {
	r.m.Remove(k)
}

func (r *bitmap) And(o bitmap2.Bitmap) {
	other := o.(*bitmap)
	r.m.And(other.m)
}

func (r *bitmap) Or(o bitmap2.Bitmap) {
	other := o.(*bitmap)
	r.m.Or(other.m)
}

func (r *bitmap) Bytes() []byte {
	buf := make([]byte, r.m.GetCardinality()*8)
	it := r.m.Iterator()
	i := 0
	for it.HasNext() {
		binary.LittleEndian.PutUint64(buf[i:], it.Next())
		i += 8
	}
	return buf
}

func (r *bitmap) Iter() iter.Seq[uint64] {
	it := r.m.Iterator()
	return func(yield func(uint64) bool) {
		for it.HasNext() {
			if !yield(it.Next()) {
				return
			}
		}
	}
}

func init() {
	bitmap2.SetProvider(prov{})
}
