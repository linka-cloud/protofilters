package protofilters

import (
	"testing"

	"go.linka.cloud/protofilters/filters"
	test "go.linka.cloud/protofilters/tests/pb"
)

var benchMatchSink int

func BenchmarkMatchScan100K(b *testing.B) {
	const total = 100_000
	const matchEvery = 10

	msgs := benchmarkBuildMatchMessages(total, matchEvery)
	simple := filters.Where("string_field").StringEquals("match")
	complex := filters.Where("string_field").StringEquals("match").
		AndWhere("number_field").NumberSup(10).
		OrWhere("bool_field").True()

	b.Run("default/simple", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchMatchSink = benchmarkMatchScanDefault(b, msgs, simple)
		}
	})

	b.Run("matcher/simple", func(b *testing.B) {
		m := NewMatcher()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchMatchSink = benchmarkMatchScanWithMatcher(b, m, msgs, simple)
		}
	})

	b.Run("matcher/complex", func(b *testing.B) {
		m := NewMatcher()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchMatchSink = benchmarkMatchScanWithMatcher(b, m, msgs, complex)
		}
	})
}

func BenchmarkMatchScan1M(b *testing.B) {
	const total = 1_000_000
	const matchEvery = 10

	msgs := benchmarkBuildMatchMessages(total, matchEvery)
	simple := filters.Where("string_field").StringEquals("match")
	m := NewMatcher()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchMatchSink = benchmarkMatchScanWithMatcher(b, m, msgs, simple)
	}
}

func benchmarkBuildMatchMessages(total, matchEvery int) []*test.Test {
	msgs := make([]*test.Test, 0, total)
	for n := 0; n < total; n++ {
		msg := &test.Test{
			StringField: "other",
			NumberField: int64(n % 100),
			BoolField:   n%7 == 0,
		}
		if n%matchEvery == 0 {
			msg.StringField = "match"
		}
		msgs = append(msgs, msg)
	}
	return msgs
}

func benchmarkMatchScanDefault(b *testing.B, msgs []*test.Test, f filters.FieldFilterer) int {
	b.Helper()
	count := 0
	for _, msg := range msgs {
		ok, err := Match(msg, f)
		if err != nil {
			b.Fatal(err)
		}
		if ok {
			count++
		}
	}
	return count
}

func benchmarkMatchScanWithMatcher(b *testing.B, m Matcher, msgs []*test.Test, f filters.FieldFilterer) int {
	b.Helper()
	count := 0
	for _, msg := range msgs {
		ok, err := m.Match(msg, f)
		if err != nil {
			b.Fatal(err)
		}
		if ok {
			count++
		}
	}
	return count
}
