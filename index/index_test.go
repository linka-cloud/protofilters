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
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"go.linka.cloud/protofilters"
	"go.linka.cloud/protofilters/filters"
	test "go.linka.cloud/protofilters/tests/pb"
)

func TestIndex(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ms := []*test.Test{
		{
			StringField: "other",
		},
		{
			StringField: "whatever",
			NumberField: 42,
		},
		{
			StringField: "whatever",
			NumberField: 43,
		},
		{
			BoolField: true,
		},
		{
			OptionalBoolField: proto.Bool(false),
		},
		{
			MessageField: &test.Test{
				StringField: "whatever",
			},
		},
		{
			RepeatedStringField: []string{"one"},
		},
		{},
	}
	f1 := filters.Where("string_field").StringEquals("whatever").
		Or(filters.Where("bool_field").True()).
		And("number_field").NumberIN(42, 43).
		Or(filters.Where("optional_bool_field").False()).
		Or(filters.Where("message_field.string_field").StringEquals("whatever")).
		Or(filters.Where("repeated_string_field").StringIN("one", "two"))
	i := New(nil, nil)
	var matches []string
	var d time.Duration
	for _, v := range ms {
		id := uuid.New().String()
		require.NoError(t, i.Index(ctx, id, v))
		n := time.Now()
		ok, err := protofilters.Match(v, f1)
		d += time.Since(n)
		require.NoError(t, err)
		if ok {
			matches = append(matches, id)
		}
	}
	t.Logf("Match took %s", d)
	sort.Strings(matches)
	d = 0
	n := time.Now()
	indexes, err := i.Find(ctx, ms[0].ProtoReflect().Descriptor().FullName(), f1)
	d += time.Since(n)
	t.Logf("Find took %s", d)
	require.NoError(t, err)
	sort.Strings(indexes)
	assert.Equal(t, matches, indexes)
	assert.Len(t, indexes, 6)
}
