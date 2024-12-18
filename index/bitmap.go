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
	"encoding/binary"
)

type Bitmap struct {
	m map[uint64]struct{}
}

func NewBitmap() *Bitmap {
	return &Bitmap{
		m: make(map[uint64]struct{}),
	}
}

func NewBitmapWith(size int) *Bitmap {
	return &Bitmap{
		m: make(map[uint64]struct{}, size),
	}
}

func NewBitmapFrom(buf []byte) *Bitmap {
	m := make(map[uint64]struct{}, len(buf)/8)
	for i := 0; i < len(buf); i += 8 {
		m[binary.LittleEndian.Uint64(buf[i:])] = struct{}{}
	}
	return &Bitmap{m: m}
}

func (b *Bitmap) Set(k uint64) {
	b.m[k] = struct{}{}
}

func (b *Bitmap) Remove(k uint64) {
	delete(b.m, k)
}

func (b *Bitmap) And(o *Bitmap) {
	for k := range b.m {
		if _, exists := o.m[k]; !exists {
			delete(b.m, k)
		}
	}
}

func (b *Bitmap) Or(o *Bitmap) {
	for k := range o.m {
		b.m[k] = struct{}{}
	}
}

func (b *Bitmap) Bytes() []byte {
	buf := make([]byte, 8*len(b.m))
	i := 0
	for k := range b.m {
		binary.LittleEndian.PutUint64(buf[i:], k)
		i += 8
	}
	return buf
}

func (b *Bitmap) NewIterator() *BitmapIterator {
	keys := make([]uint64, 0, len(b.m))
	for k := range b.m {
		keys = append(keys, k)
	}
	return &BitmapIterator{v: keys}
}

type BitmapIterator struct {
	v []uint64
	i int
}

func (i *BitmapIterator) Next() uint64 {
	if i.i >= len(i.v) {
		return 0
	}
	v := i.v[i.i]
	i.i++
	return v
}
