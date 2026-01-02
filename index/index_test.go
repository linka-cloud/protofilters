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
	"slices"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"

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
			MessageField: &test.Test{
				RepeatedMessageField: []*test.Test{
					{
						StringField: "whatever",
					},
				},
			},
		},
		{
			RepeatedStringField: []string{"one"},
		},
		{
			TimeValueField: timestamppb.New(time.Now().Add(time.Hour)),
		},
		{},
	}
	f1 := filters.Where("string_field").StringEquals("whatever").
		Or(filters.Where("bool_field").True()).
		AndWhere("number_field").NumberIN(42, 43).
		OrWhere("optional_bool_field").False().
		OrWhere("message_field.string_field").StringEquals("whatever").
		OrWhere("repeated_string_field").StringIN("one", "two").
		OrWhere("message_field.repeated_message_field.string_field").StringEquals("whatever").
		OrWhere("time_value_field").TimeAfter(time.Now())

	i := New(nil, func(ctx context.Context, name protoreflect.FullName, fds ...protoreflect.FieldDescriptor) (bool, error) {
		if name != "linka.cloud.test.Test" {
			return false, nil
		}
		var f protoreflect.FullName
		for _, v := range fds {
			f = f.Append(v.Name())
		}
		switch f {
		case "string_field",
			"bool_field",
			"number_field",
			"optional_bool_field",
			"message_field.string_field",
			"repeated_string_field",
			"message_field.repeated_message_field.string_field",
			"time_value_field":
			return true, nil
		default:
			return false, nil
		}
	})
	var matches []string
	msgs := make(map[string]*test.Test)
	var d time.Duration
	var di time.Duration
	count := len(ms) * 100_000
	t.Logf("Inserting %d records", count)
	for j := range count {
		v := ms[j%len(ms)]
		id := uuid.New().String()
		msgs[id] = v
		n := time.Now()
		require.NoError(t, i.Insert(ctx, id, v))
		di += time.Since(n)
		n = time.Now()
		ok, err := protofilters.Match(v, f1)
		d += time.Since(n)
		require.NoError(t, err)
		if ok {
			matches = append(matches, id)
		}
	}
	t.Logf("Insert took %s", di)
	t.Logf("Match took %s", d)
	sort.Strings(matches)
	n := time.Now()
	keys, collisions, err := i.Find(ctx, ms[0].ProtoReflect().Descriptor().FullName(), f1)
	t.Logf("Found %d collisions", len(collisions))
	for _, v := range collisions {
		ok, err := protofilters.Match(msgs[v], f1)
		require.NoError(t, err)
		if ok {
			keys = append(keys, v)
		}
	}
	di = time.Since(n)
	t.Logf("Find took %s", di)
	t.Logf("Ratio: %.2fx", float64(d)/float64(di))
	require.NoError(t, err)
	sort.Strings(keys)
	assert.Equal(t, matches, keys)
	assert.Len(t, keys, 8*count/len(ms))
	if !slices.Equal(keys, matches) {
		s := slices.DeleteFunc(keys, func(s string) bool {
			return !assert.Contains(t, matches, s)
		})
		for _, v := range s {
			t.Logf("Missing: %s", msgs[v])
		}
	}
}
