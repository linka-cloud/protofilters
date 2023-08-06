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
	"fmt"
)

// Iterator is an interface for iterating over a collection of values
// It is used by the Index to iterate over the values of a field that can be stored in a remote database
type Iterator[T any] interface {
	Next() bool
	Value() (T, error)
}

type sliceIterator[T any] struct {
	slice []T
	index int
}

func (s *sliceIterator[T]) Next() bool {
	if s.index >= len(s.slice) {
		return false
	}
	s.index++
	return true
}

func (s *sliceIterator[T]) Value() (t T, err error) {
	if s.index > len(s.slice) {
		return t, fmt.Errorf("iterator out of bounds")
	}
	return s.slice[s.index-1], nil
}
