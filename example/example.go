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

package main

import (
	"log"

	"google.golang.org/protobuf/types/known/wrapperspb"

	pf "go.linka.cloud/protofilters"
	"go.linka.cloud/protofilters/matcher"
	test "go.linka.cloud/protofilters/tests/pb"
)

func main() {
	m := &test.Test{
		BoolField:      true,
		BoolValueField: wrapperspb.Bool(false),
	}
	ok, err := matcher.MatchFilters(m, &pf.FieldFilter{Field: "bool_field", Filter: pf.True()})
	if err != nil {
		log.Fatalln(err)
	}
	if !ok {
		log.Fatalln("should be true")
	}
	ok, err = matcher.MatchFilters(m, &pf.FieldFilter{Field: "bool_value_field", Filter: pf.False()})
	if err != nil {
		log.Fatalln(err)
	}
	if !ok {
		log.Fatalln("should be true")
	}
}
