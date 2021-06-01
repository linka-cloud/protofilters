package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pf "go.linka.cloud/protofilters"
	test "go.linka.cloud/protofilters/tests/pb"
)

func TestFieldFilter(t *testing.T) {
	assert := assert.New(t)
	m := &test.Test{StringField: "ok"}
	ok, err := Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"noop": pf.StringEquals("ok"),
	}})
	assert.Error(err)
	assert.False(ok)
	ok, err = Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"messageField": pf.Null(),
	}})
	assert.Error(err)
	assert.False(ok)
	assert.True(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"string_field": pf.StringEquals("ok"),
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"string_field": pf.StringIN("other", "ok"),
	}}))
	assert.False(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"string_field": pf.StringIN("other", "noop"),
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"string_field": pf.StringNotIN(),
		"enum_field":   pf.StringIN("NONE"),
	}}))
	assert.False(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"string_field": pf.StringNotRegex(`[a-z](.+)`),
	}}))
	m.EnumField = test.Test_Type(42)
	assert.False(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"enum_field": pf.StringIN("OTHER"),
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"message_field": pf.Null(),
	}}))
	m.MessageField = m
	assert.False(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"message_field": pf.Null(),
	}}))
	ok, err = Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"message_field.string_field.message_field": pf.Null(),
	}})
	assert.Error(err)
	assert.False(ok)
	assert.False(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"message_field.message_field.message_field": pf.Null(),
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"message_field.message_field.message_field.string_field": pf.StringIN("ok"),
	}}))

	assert.False(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"enum_field": pf.StringIN("OTHER"),
	}}))

	m.NumberField = 42
	assert.False(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"number_field": pf.NumberEquals(0),
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"number_field": pf.NumberEquals(42),
	}}))
	assert.False(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"number_field": pf.NumberIN(0, 22),
	}}))
	assert.True(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"enum_field": pf.NumberIN(0, 42),
	}}))
	assert.False(Match(m, &pf.FieldsFilter{Filters: map[string]*pf.Filter{
		"enum_field": pf.NumberNotIN(0, 42),
	}}))

}
