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
	"slices"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"go.linka.cloud/protofilters"
	"go.linka.cloud/protofilters/filters"
	_ "go.linka.cloud/protofilters/index/bitmap/sroar"
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
		id := strconv.Itoa(j + 1) // use simple integer ids for performance
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
		if (j+1)%100_000 == 0 {
			t.Logf("Inserted %d records", j+1)
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

func TestUpdateClearsRepeatedFields(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	msg1 := &test.Test{
		RepeatedStringField: []string{"one", "two"},
		RepeatedMessageField: []*test.Test{
			{StringField: "a"},
			{StringField: "b"},
		},
	}
	require.NoError(t, i.Insert(ctx, "1", msg1))

	keys, collisions, err := i.Find(ctx, msg1.ProtoReflect().Descriptor().FullName(),
		filters.Where("repeated_string_field").StringEquals("one"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg1.ProtoReflect().Descriptor().FullName(),
		filters.Where("repeated_message_field.string_field").StringEquals("a"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")

	msg2 := &test.Test{
		RepeatedStringField: []string{"two"},
		RepeatedMessageField: []*test.Test{
			{StringField: "b"},
			{StringField: "c"},
		},
	}
	require.NoError(t, i.Update(ctx, "1", msg1, msg2))

	keys, collisions, err = i.Find(ctx, msg2.ProtoReflect().Descriptor().FullName(),
		filters.Where("repeated_string_field").StringEquals("one"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg2.ProtoReflect().Descriptor().FullName(),
		filters.Where("repeated_message_field.string_field").StringEquals("a"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg2.ProtoReflect().Descriptor().FullName(),
		filters.Where("repeated_string_field").StringEquals("two"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg2.ProtoReflect().Descriptor().FullName(),
		filters.Where("repeated_message_field.string_field").StringEquals("c"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")
}

func TestUpdateClearsEmptyListsAndNestedMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	msg1 := &test.Test{
		MessageField:         &test.Test{StringField: "nested"},
		RepeatedStringField:  []string{"one"},
		RepeatedMessageField: []*test.Test{{StringField: "a"}},
	}
	require.NoError(t, i.Insert(ctx, "1", msg1))

	msg2 := &test.Test{}
	require.NoError(t, i.Update(ctx, "1", msg1, msg2))

	keys, collisions, err := i.Find(ctx, msg1.ProtoReflect().Descriptor().FullName(),
		filters.Where("repeated_string_field").StringEquals("one"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg1.ProtoReflect().Descriptor().FullName(),
		filters.Where("repeated_message_field.string_field").StringEquals("a"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg1.ProtoReflect().Descriptor().FullName(),
		filters.Where("message_field.string_field").StringEquals("nested"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")
}

func TestFindReportsCollisions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := collisionStore{Store: newStore()}
	i := New(cs, All)
	require.NoError(t, i.Insert(ctx, "1", &test.Test{StringField: "value"}))

	keys, collisions, err := i.Find(ctx, "linka.cloud.test.Test",
		filters.Where("string_field").StringEquals("value"),
	)
	require.NoError(t, err)
	assert.Empty(t, keys)
	assert.Contains(t, collisions, "1")
	assert.Contains(t, collisions, "collision")
}

func TestIndexLookupMixedTypes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := time.Now()
	msg := &test.Test{
		StringField:        "hello",
		NumberField:        7,
		BoolField:          true,
		EnumField:          test.Test_ONE,
		NumberValueField:   wrapperspb.Int64(10),
		StringValueField:   wrapperspb.String("s"),
		BoolValueField:     wrapperspb.Bool(true),
		TimeValueField:     timestamppb.New(now.Add(time.Minute)),
		DurationValueField: durationpb.New(2 * time.Second),
		OptionalBoolField:  proto.Bool(false),
	}

	i := New(nil, All)
	require.NoError(t, i.Insert(ctx, "1", msg))

	filter := filters.Where("enum_field").StringEquals("ONE").
		AndWhere("number_value_field").NumberEquals(10).
		AndWhere("string_value_field").StringEquals("s").
		AndWhere("bool_value_field").True().
		AndWhere("time_value_field").TimeAfter(now).
		AndWhere("duration_value_field").DurationEquals(2 * time.Second).
		AndWhere("optional_bool_field").False()

	keys, collisions, err := i.Find(ctx, msg.ProtoReflect().Descriptor().FullName(), filter)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")
}

func TestUpdateOptionalUnsetAndNullFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	msg1 := &test.Test{OptionalBoolField: proto.Bool(true)}
	require.NoError(t, i.Insert(ctx, "1", msg1))

	msg2 := &test.Test{}
	require.NoError(t, i.Update(ctx, "1", msg1, msg2))

	keys, collisions, err := i.Find(ctx, msg1.ProtoReflect().Descriptor().FullName(),
		filters.Where("optional_bool_field").True(),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg1.ProtoReflect().Descriptor().FullName(),
		filters.Where("optional_bool_field").Null(),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")
}

func TestUpdateOptionalZeroVsUnset(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	msg1 := &test.Test{OptionalNumberField: proto.Int64(0)}
	require.NoError(t, i.Insert(ctx, "1", msg1))

	keys, collisions, err := i.Find(ctx, msg1.ProtoReflect().Descriptor().FullName(),
		filters.Where("optional_number_field").NumberEquals(0),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")

	msg2 := &test.Test{}
	require.NoError(t, i.Update(ctx, "1", msg1, msg2))

	keys, collisions, err = i.Find(ctx, msg1.ProtoReflect().Descriptor().FullName(),
		filters.Where("optional_number_field").NumberEquals(0),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg1.ProtoReflect().Descriptor().FullName(),
		filters.Where("optional_number_field").Null(),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")
}

func TestRepeatedDuplicateValuesUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	msg1 := &test.Test{RepeatedStringField: []string{"dup", "dup"}}
	require.NoError(t, i.Insert(ctx, "1", msg1))

	msg2 := &test.Test{RepeatedStringField: []string{"dup"}}
	require.NoError(t, i.Update(ctx, "1", msg1, msg2))

	keys, collisions, err := i.Find(ctx, msg2.ProtoReflect().Descriptor().FullName(),
		filters.Where("repeated_string_field").StringEquals("dup"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")

	msg3 := &test.Test{RepeatedStringField: []string{}}
	require.NoError(t, i.Update(ctx, "1", msg2, msg3))

	keys, collisions, err = i.Find(ctx, msg3.ProtoReflect().Descriptor().FullName(),
		filters.Where("repeated_string_field").StringEquals("dup"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")
}

func TestUpdateOneofSwitch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	msg1 := &test.Test{Choice: &test.Test_OneofStringField{OneofStringField: "a"}}
	require.NoError(t, i.Insert(ctx, "1", msg1))

	msg2 := &test.Test{Choice: &test.Test_OneofNumberField{OneofNumberField: 2}}
	require.NoError(t, i.Update(ctx, "1", msg1, msg2))

	keys, collisions, err := i.Find(ctx, msg2.ProtoReflect().Descriptor().FullName(),
		filters.Where("oneof_string_field").StringEquals("a"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg2.ProtoReflect().Descriptor().FullName(),
		filters.Where("oneof_number_field").NumberEquals(2),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")
}

func TestOneofUnsetDoesNotMatchDefaultValue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	require.NoError(t, i.Insert(ctx, "1", &test.Test{Choice: &test.Test_OneofStringField{OneofStringField: "a"}}))
	require.NoError(t, i.Insert(ctx, "2", &test.Test{Choice: &test.Test_OneofNumberField{OneofNumberField: 0}}))
	require.NoError(t, i.Insert(ctx, "3", &test.Test{}))
	require.NoError(t, i.Insert(ctx, "4", &test.Test{Choice: &test.Test_OneofStringField{OneofStringField: ""}}))

	keys, collisions, err := i.Find(ctx, "linka.cloud.test.Test",
		filters.Where("oneof_number_field").NumberEquals(0),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "2")
	assert.NotContains(t, keys, "1")
	assert.NotContains(t, keys, "3")
	assert.NotContains(t, keys, "4")

	keys, collisions, err = i.Find(ctx, "linka.cloud.test.Test",
		filters.Where("oneof_string_field").StringEquals(""),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "4")
	assert.NotContains(t, keys, "1")
	assert.NotContains(t, keys, "2")
	assert.NotContains(t, keys, "3")
}

func TestOneofNotEqualsRequiresSelectedField(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	require.NoError(t, i.Insert(ctx, "1", &test.Test{Choice: &test.Test_OneofStringField{OneofStringField: "a"}}))
	require.NoError(t, i.Insert(ctx, "2", &test.Test{Choice: &test.Test_OneofNumberField{OneofNumberField: 0}}))
	require.NoError(t, i.Insert(ctx, "3", &test.Test{}))
	require.NoError(t, i.Insert(ctx, "4", &test.Test{Choice: &test.Test_OneofNumberField{OneofNumberField: 1}}))

	keys, collisions, err := i.Find(ctx, "linka.cloud.test.Test",
		filters.Where("oneof_number_field").NumberNotEquals(0),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "4")
	assert.NotContains(t, keys, "1")
	assert.NotContains(t, keys, "2")
	assert.NotContains(t, keys, "3")
}

func TestRepeatedEmptyStringAndClear(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	msg1 := &test.Test{RepeatedStringField: []string{"", "value"}}
	require.NoError(t, i.Insert(ctx, "1", msg1))

	keys, collisions, err := i.Find(ctx, msg1.ProtoReflect().Descriptor().FullName(),
		filters.Where("repeated_string_field").StringEquals(""),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")

	msg2 := &test.Test{RepeatedStringField: []string{"value"}}
	require.NoError(t, i.Update(ctx, "1", msg1, msg2))

	keys, collisions, err = i.Find(ctx, msg2.ProtoReflect().Descriptor().FullName(),
		filters.Where("repeated_string_field").StringEquals(""),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")
}

func TestWKTypeUnsetDoesNotMatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	msg := &test.Test{}
	require.NoError(t, i.Insert(ctx, "1", msg))

	keys, collisions, err := i.Find(ctx, msg.ProtoReflect().Descriptor().FullName(),
		filters.Where("time_value_field").TimeAfter(time.Now()),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg.ProtoReflect().Descriptor().FullName(),
		filters.Where("duration_value_field").DurationSup(time.Second),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg.ProtoReflect().Descriptor().FullName(),
		filters.Where("string_value_field").StringEquals("value"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")
}

func TestEnumUnknownNumber(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	msg := &test.Test{EnumField: test.Test_Type(99)}
	require.NoError(t, i.Insert(ctx, "1", msg))

	keys, collisions, err := i.Find(ctx, msg.ProtoReflect().Descriptor().FullName(),
		filters.Where("enum_field").StringEquals("ONE"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg.ProtoReflect().Descriptor().FullName(),
		filters.Where("enum_field").NumberEquals(99),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")
}

func TestNumericAndHashedKeys(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	require.NoError(t, i.Insert(ctx, "1", &test.Test{StringField: "value"}))
	require.NoError(t, i.Insert(ctx, "abc", &test.Test{StringField: "value"}))

	keys, collisions, err := i.Find(ctx, "linka.cloud.test.Test",
		filters.Where("string_field").StringEquals("value"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")
	assert.Contains(t, keys, "abc")

	require.NoError(t, i.Update(ctx, "1", &test.Test{StringField: "value"}, &test.Test{StringField: "other"}))
	keys, collisions, err = i.Find(ctx, "linka.cloud.test.Test",
		filters.Where("string_field").StringEquals("value"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")
	assert.Contains(t, keys, "abc")
}

func TestCreateUpdateDeleteRemovesStaleEntries(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, id := range []string{"1", "abc"} {
		t.Run(id, func(t *testing.T) {
			s := newStore()
			i := New(s, All)
			oldMsg := &test.Test{StringField: "created", NumberField: 1}
			newMsg := &test.Test{StringField: "updated", NumberField: 2}

			require.NoError(t, i.Insert(ctx, id, oldMsg))
			assertStoreHasValueForKey(t, ctx, s, "linka.cloud.test.Test", "string_field", protoreflect.ValueOfString("created"), id, true)

			keys, collisions, err := i.Find(ctx, "linka.cloud.test.Test",
				filters.Where("string_field").StringEquals("created"),
			)
			require.NoError(t, err)
			require.Empty(t, collisions)
			assert.Contains(t, keys, id)

			require.NoError(t, i.Update(ctx, id, oldMsg, newMsg))

			keys, collisions, err = i.Find(ctx, "linka.cloud.test.Test",
				filters.Where("string_field").StringEquals("created"),
			)
			require.NoError(t, err)
			require.Empty(t, collisions)
			assert.NotContains(t, keys, id)
			assertStoreHasValueForKey(t, ctx, s, "linka.cloud.test.Test", "string_field", protoreflect.ValueOfString("created"), id, false)

			keys, collisions, err = i.Find(ctx, "linka.cloud.test.Test",
				filters.Where("string_field").StringEquals("updated"),
			)
			require.NoError(t, err)
			require.Empty(t, collisions)
			assert.Contains(t, keys, id)
			assertStoreHasValueForKey(t, ctx, s, "linka.cloud.test.Test", "string_field", protoreflect.ValueOfString("updated"), id, true)

			require.NoError(t, i.Remove(ctx, id))

			keys, collisions, err = i.Find(ctx, "linka.cloud.test.Test",
				filters.Where("string_field").StringEquals("updated"),
			)
			require.NoError(t, err)
			require.Empty(t, collisions)
			assert.NotContains(t, keys, id)
			assertStoreHasValueForKey(t, ctx, s, "linka.cloud.test.Test", "string_field", protoreflect.ValueOfString("updated"), id, false)
		})
	}
}

func assertStoreHasValueForKey(t *testing.T, ctx context.Context, s Store, typ protoreflect.FullName, field protoreflect.Name, val protoreflect.Value, key string, want bool) {
	t.Helper()

	fr, err := s.For(ctx, typ)
	require.NoError(t, err)

	h := keyHash(key)
	has := false
	for f, err := range fr.Get(ctx, field) {
		require.NoError(t, err)
		if !valueEqual(f.Value(), val) {
			continue
		}
		b, err := f.Bitmap(ctx)
		require.NoError(t, err)
		has = b.Contains(h)
		break
	}
	assert.Equal(t, want, has)
}

func TestDeepNestedUpdateClears(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	msg1 := &test.Test{
		MessageField: &test.Test{
			RepeatedMessageField: []*test.Test{{StringField: "deep"}},
		},
	}
	require.NoError(t, i.Insert(ctx, "1", msg1))

	msg2 := &test.Test{MessageField: &test.Test{}}
	require.NoError(t, i.Update(ctx, "1", msg1, msg2))

	keys, collisions, err := i.Find(ctx, msg1.ProtoReflect().Descriptor().FullName(),
		filters.Where("message_field.repeated_message_field.string_field").StringEquals("deep"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")
}

func TestDeepNestedUpdateMultipleBranches(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := New(nil, All)
	msg1 := &test.Test{
		MessageField: &test.Test{
			MessageField: &test.Test{StringField: "deep"},
			RepeatedMessageField: []*test.Test{
				{StringField: "a"},
				{StringField: "b"},
			},
		},
	}
	require.NoError(t, i.Insert(ctx, "1", msg1))

	msg2 := &test.Test{
		MessageField: &test.Test{
			MessageField: &test.Test{StringField: "deep2"},
			RepeatedMessageField: []*test.Test{
				{StringField: "a"},
				{StringField: "b"},
			},
		},
	}
	require.NoError(t, i.Update(ctx, "1", msg1, msg2))

	keys, collisions, err := i.Find(ctx, msg2.ProtoReflect().Descriptor().FullName(),
		filters.Where("message_field.repeated_message_field.string_field").StringEquals("a"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg2.ProtoReflect().Descriptor().FullName(),
		filters.Where("message_field.message_field.string_field").StringEquals("deep"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, msg2.ProtoReflect().Descriptor().FullName(),
		filters.Where("message_field.message_field.string_field").StringEquals("deep2"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Contains(t, keys, "1")
}

func TestCollisionMultipleKeys(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := collisionStore{Store: newStore()}
	i := New(cs, All)
	require.NoError(t, i.Insert(ctx, "1", &test.Test{StringField: "value"}))
	require.NoError(t, i.Insert(ctx, "2", &test.Test{StringField: "value"}))

	keys, collisions, err := i.Find(ctx, "linka.cloud.test.Test",
		filters.Where("string_field").StringEquals("value"),
	)
	require.NoError(t, err)
	assert.Empty(t, keys)
	assert.Contains(t, collisions, "1")
	assert.Contains(t, collisions, "2")
	assert.Contains(t, collisions, "collision")
}

func TestCollisionResolutionFiltering(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := collisionStore{Store: newStore()}
	i := New(cs, All)
	msg1 := &test.Test{StringField: "match"}
	msg2 := &test.Test{StringField: "other"}
	msgs := map[string]*test.Test{
		"1": msg1,
		"2": msg2,
	}
	require.NoError(t, i.Insert(ctx, "1", msg1))
	require.NoError(t, i.Insert(ctx, "2", msg2))

	filter := filters.Where("string_field").StringEquals("match")
	keys, collisions, err := i.Find(ctx, "linka.cloud.test.Test", filter)
	require.NoError(t, err)
	assert.Empty(t, keys)
	assert.Contains(t, collisions, "1")

	var resolved []string
	for _, v := range collisions {
		m, ok := msgs[v]
		if !ok {
			continue
		}
		ok, err := protofilters.Match(m, filter)
		require.NoError(t, err)
		if ok {
			resolved = append(resolved, v)
		}
	}
	assert.Equal(t, []string{"1"}, resolved)
}

func TestFuncFilterSkipsFields(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fn := func(_ context.Context, _ protoreflect.FullName, fds ...protoreflect.FieldDescriptor) (bool, error) {
		var f protoreflect.FullName
		for _, v := range fds {
			f = f.Append(v.Name())
		}
		return f == "string_field", nil
	}
	i := New(nil, fn)
	require.NoError(t, i.Insert(ctx, "1", &test.Test{StringField: "a", NumberField: 1}))
	require.NoError(t, i.Update(ctx, "1", &test.Test{StringField: "a", NumberField: 1}, &test.Test{StringField: "b", NumberField: 2}))

	keys, collisions, err := i.Find(ctx, "linka.cloud.test.Test",
		filters.Where("string_field").StringEquals("a"),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.NotContains(t, keys, "1")

	keys, collisions, err = i.Find(ctx, "linka.cloud.test.Test",
		filters.Where("number_field").NumberEquals(2),
	)
	require.NoError(t, err)
	require.Empty(t, collisions)
	assert.Empty(t, keys)
}

func TestUIDIndexFindOptions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ui := NewUID(nil, All)
	for _, id := range []uint64{10, 2, 30, 4, 5} {
		require.NoError(t, ui.Insert(ctx, id, &test.Test{StringField: "value"}))
	}

	filter := filters.Where("string_field").StringEquals("value")

	var asc []uint64
	for uid, err := range ui.Find(ctx, "linka.cloud.test.Test", filter, FindOptions{}) {
		require.NoError(t, err)
		asc = append(asc, uid)
	}
	assert.Equal(t, []uint64{2, 4, 5, 10, 30}, asc)

	var paged []uint64
	for uid, err := range ui.Find(ctx, "linka.cloud.test.Test", filter, FindOptions{Offset: 1, Limit: 2}) {
		require.NoError(t, err)
		paged = append(paged, uid)
	}
	assert.Equal(t, []uint64{4, 5}, paged)

	var reversed []uint64
	for uid, err := range ui.Find(ctx, "linka.cloud.test.Test", filter, FindOptions{Reverse: true, Offset: 1, Limit: 2}) {
		require.NoError(t, err)
		reversed = append(reversed, uid)
	}
	assert.Equal(t, []uint64{10, 5}, reversed)
}

func TestUIDIndexUpdateAndRemove(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ui := NewUID(nil, All)
	require.NoError(t, ui.Insert(ctx, 42, &test.Test{StringField: "one"}))

	uids, err := collectUIDs(ui.Find(ctx, "linka.cloud.test.Test", filters.Where("string_field").StringEquals("one"), FindOptions{}))
	require.NoError(t, err)
	assert.Equal(t, []uint64{42}, uids)

	require.NoError(t, ui.Update(ctx, 42, &test.Test{StringField: "one"}, &test.Test{StringField: "two"}))

	uids, err = collectUIDs(ui.Find(ctx, "linka.cloud.test.Test", filters.Where("string_field").StringEquals("one"), FindOptions{}))
	require.NoError(t, err)
	assert.Empty(t, uids)

	uids, err = collectUIDs(ui.Find(ctx, "linka.cloud.test.Test", filters.Where("string_field").StringEquals("two"), FindOptions{}))
	require.NoError(t, err)
	assert.Equal(t, []uint64{42}, uids)

	require.NoError(t, ui.Remove(ctx, 42))

	uids, err = collectUIDs(ui.Find(ctx, "linka.cloud.test.Test", filters.Where("string_field").StringEquals("two"), FindOptions{}))
	require.NoError(t, err)
	assert.Empty(t, uids)
}

func TestUIDIndexFindEmptyFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ui := NewUID(nil, All)
	require.NoError(t, ui.Insert(ctx, 1, &test.Test{StringField: "value"}))

	uids, err := collectUIDs(ui.Find(ctx, "linka.cloud.test.Test", nil, FindOptions{}))
	require.NoError(t, err)
	assert.Empty(t, uids)
}

func collectUIDs(seq iter.Seq2[uint64, error]) ([]uint64, error) {
	var out []uint64
	for uid, err := range seq {
		if err != nil {
			return nil, err
		}
		out = append(out, uid)
	}
	return out, nil
}

type collisionStore struct {
	Store
}

func (c collisionStore) Keys(ctx context.Context, i uint64) ([]string, error) {
	ks, err := c.Store.Keys(ctx, i)
	if err != nil || len(ks) == 0 {
		return ks, err
	}
	if len(ks) == 1 {
		return []string{ks[0], "collision"}, nil
	}
	return ks, nil
}
