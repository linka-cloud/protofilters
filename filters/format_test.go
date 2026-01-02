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

package filters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    FieldFilterer
		expected string
	}{
		{"NullField", Where("test").Null(), "test is null"},
		{"StringEquals", Where("name").StringEquals("John"), "name eq 'John'"},
		{"StringNotEquals", Where("name").StringNotEquals("John"), "name not eq 'John'"},
		{"StringIEquals", Where("name").StringIEquals("john"), "name ieq 'john'"},
		{"StringHasPrefix", Where("name").StringHasPrefix("Jo"), "name has_prefix 'Jo'"},
		{"StringNotHasPrefix", Where("name").StringNotHasPrefix("Jo"), "name not has_prefix 'Jo'"},
		{"StringIHasPrefix", Where("name").StringIHasPrefix("jo"), "name ihas_prefix 'jo'"},
		{"StringNotIHasPrefix", Where("name").StringNotIHasPrefix("jo"), "name not ihas_prefix 'jo'"},
		{"StringHasSuffix", Where("name").StringHasSuffix("hn"), "name has_suffix 'hn'"},
		{"StringNotHasSuffix", Where("name").StringNotHasSuffix("hn"), "name not has_suffix 'hn'"},
		{"StringIHasSuffix", Where("name").StringIHasSuffix("HN"), "name ihas_suffix 'HN'"},
		{"StringNotIHasSuffix", Where("name").StringNotIHasSuffix("HN"), "name not ihas_suffix 'HN'"},
		{"StringRegex", Where("name").StringRegex("Jo.*"), "name matches 'Jo.*'"},
		{"StringNotRegex", Where("name").StringNotRegex("Jo.*"), "name not matches 'Jo.*'"},
		{"StringIN", Where("name").StringIN("John", "Doe"), "name in ('John', 'Doe')"},
		{"StringNotIN", Where("name").StringNotIN("John", "Doe"), "name not in ('John', 'Doe')"},
		{"NumberEquals", Where("age").NumberEquals(30), "age eq 30"},
		{"NumberNotEquals", Where("age").NumberNotEquals(30), "age not eq 30"},
		{"NumberInf", Where("age").NumberInf(30), "age inf 30"},
		{"NumberSup", Where("age").NumberSup(30), "age sup 30"},
		{"NumberIN", Where("age").NumberIN(25, 30), "age in (25, 30)"},
		{"NumberNotIN", Where("age").NumberNotIN(25, 30), "age not in (25, 30)"},
		{"True", Where("active").True(), "active is true"},
		{"False", Where("active").False(), "active is false"},
		{"Null", Where("data").Null(), "data is null"},
		{"NotNull", Where("data").NotNull(), "data not is null"},
		{"DurationEquals", Where("timeout").DurationEquals(5 * 60 * 1000), "timeout eq 300µs"},
		{"DurationNotEquals", Where("timeout").DurationNotEquals(5 * 60 * 1000), "timeout not eq 300µs"},
		{"DurationInf", Where("timeout").DurationInf(5 * 60 * 1000), "timeout inf 300µs"},
		{"DurationSup", Where("timeout").DurationSup(5 * 60 * 1000), "timeout sup 300µs"},
		{"StringInf", Where("name").StringInf("A"), "name inf 'A'"},
		{"StringSup", Where("name").StringSup("Z"), "name sup 'Z'"},
		{"StringIInf", Where("name").StringIInf("a"), "name iinf 'a'"},
		{"StringISup", Where("name").StringISup("z"), "name isup 'z'"},
		{"TimeEquals", Where("created").TimeEquals(time.Unix(0, 0)), "created eq 1970-01-01T00:00:00Z"},
		{"TimeNotEquals", Where("created").TimeNotEquals(time.Unix(0, 0)), "created not eq 1970-01-01T00:00:00Z"},
		{"TimeAfter", Where("created").TimeAfter(time.Unix(0, 0)), "created after 1970-01-01T00:00:00Z"},
		{"TimeBefore", Where("created").TimeBefore(time.Unix(0, 0)), "created before 1970-01-01T00:00:00Z"},
		{
			"And simple",
			Where("name").StringEquals("John").AndWhere("age").NumberSup(18),
			"name eq 'John' and age sup 18",
		},
		{
			"Or simple",
			Where("name").StringEquals("John").OrWhere("age").StringEquals("Doe"),
			"name eq 'John' or age eq 'Doe'",
		},
		{
			"And nested simple",
			Where("name").StringEquals("John").And(Where("age").NumberSup(18)),
			"name eq 'John' and age sup 18",
		},
		{
			"Or nested simple",
			Where("name").StringEquals("John").Or(Where("age").NumberSup(18)),
			"name eq 'John' or age sup 18",
		},
		{
			"Complex nested",
			Where("name").StringEquals("John").And(Where("age").NumberSup(18).OrWhere("active").True()),
			"name eq 'John' and (age sup 18 or active is true)",
		},
		{
			"Deeply nested",
			Where("a").StringEquals("x").
				And(Where("b").StringEquals("y").OrWhere("c").StringEquals("z")).
				Or(Where("d").True().
					And(Where("e").StringEquals("w").OrWhere("f").StringEquals("v")),
				),
			"a eq 'x' and (b eq 'y' or c eq 'z') or (d is true and (e eq 'w' or f eq 'v'))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.Expr().Format())
		})
	}
}
