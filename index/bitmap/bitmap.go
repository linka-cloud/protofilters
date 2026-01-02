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

package bitmap

import (
	"iter"
)

var p Provider

func New() Bitmap {
	if p == nil {
		panic("no bitmap implementation imported")
	}
	return p.New()
}

func NewWith(n int) Bitmap {
	if p == nil {
		panic("no bitmap implementation imported")
	}
	return p.NewWith(n)
}

func NewFrom(buf []byte) Bitmap {
	if p == nil {
		panic("no bitmap implementation imported")
	}
	return p.NewFrom(buf)
}

type Provider interface {
	New() Bitmap
	NewWith(n int) Bitmap
	NewFrom(buf []byte) Bitmap
}

type Bitmap interface {
	Set(k uint64)
	Remove(k uint64)
	And(o Bitmap)
	Or(o Bitmap)
	Bytes() []byte
	Iter() iter.Seq[uint64]
}

type BitmapIterator interface {
	Next() uint64
}

func SetProvider(prv Provider) {
	p = prv
}
