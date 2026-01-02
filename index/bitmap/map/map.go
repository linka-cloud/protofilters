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

package simplemap

import (
	"encoding/binary"
	"iter"

	bitmap2 "go.linka.cloud/protofilters/index/bitmap"
)

var (
	_ bitmap2.Provider = (*prov)(nil)
	_ bitmap2.Bitmap   = (*bitmap)(nil)
)

type prov struct{}

func (prov) New() bitmap2.Bitmap {
	return &bitmap{
		m: make(map[uint64]struct{}),
	}
}

func (prov) NewWith(size int) bitmap2.Bitmap {
	return &bitmap{
		m: make(map[uint64]struct{}, size),
	}
}

func (prov) NewFrom(buf []byte) bitmap2.Bitmap {
	m := make(map[uint64]struct{}, len(buf)/8)
	for i := 0; i < len(buf); i += 8 {
		m[binary.LittleEndian.Uint64(buf[i:])] = struct{}{}
	}
	return &bitmap{m: m}
}

type bitmap struct {
	m map[uint64]struct{}
}

func (b *bitmap) Set(k uint64) {
	b.m[k] = struct{}{}
}

func (b *bitmap) Remove(k uint64) {
	delete(b.m, k)
}

func (b *bitmap) And(other bitmap2.Bitmap) {
	o := other.(*bitmap)
	for k := range b.m {
		if _, exists := o.m[k]; !exists {
			delete(b.m, k)
		}
	}
}

func (b *bitmap) Or(other bitmap2.Bitmap) {
	o := other.(*bitmap)
	for k := range o.m {
		b.m[k] = struct{}{}
	}
}

func (b *bitmap) Bytes() []byte {
	buf := make([]byte, 8*len(b.m))
	i := 0
	for k := range b.m {
		binary.LittleEndian.PutUint64(buf[i:], k)
		i += 8
	}
	return buf
}

func (b *bitmap) Iter() iter.Seq[uint64] {
	keys := make([]uint64, 0, len(b.m))
	for k := range b.m {
		keys = append(keys, k)
	}
	return func(yield func(uint64) bool) {
		for _, k := range keys {
			if !yield(k) {
				return
			}
		}
	}
}

func init() {
	bitmap2.SetProvider(prov{})
}
