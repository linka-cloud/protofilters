package index

import (
	"context"
	"strconv"
	"testing"

	"go.linka.cloud/protofilters/filters"
	test "go.linka.cloud/protofilters/tests/pb"
)

var (
	benchKeysSink int
	benchUIDSink  uint64
)

func BenchmarkIndexFind(b *testing.B) {
	ctx := context.Background()
	const total = 100_000
	const matchEvery = 10

	filter := filters.Where("string_field").StringEquals("match")
	keyIndex := benchmarkBuildKeyIndex(b, ctx, total, matchEvery)
	uidIndex := benchmarkBuildUIDIndex(b, ctx, total, matchEvery)

	b.Run("key/find", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			keys, collisions, err := keyIndex.Find(ctx, "linka.cloud.test.Test", filter)
			if err != nil {
				b.Fatal(err)
			}
			benchKeysSink = len(keys) + len(collisions)
		}
	})

	b.Run("uid/find", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			count := 0
			for uid, err := range uidIndex.Find(ctx, "linka.cloud.test.Test", filter, FindOptions{}) {
				if err != nil {
					b.Fatal(err)
				}
				benchUIDSink = uid
				count++
			}
			benchKeysSink = count
		}
	})

	b.Run("uid/find-limit-100", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			count := 0
			for uid, err := range uidIndex.Find(ctx, "linka.cloud.test.Test", filter, FindOptions{Limit: 100}) {
				if err != nil {
					b.Fatal(err)
				}
				benchUIDSink = uid
				count++
			}
			benchKeysSink = count
		}
	})

	b.Run("uid/find-limit-100-reverse", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			count := 0
			for uid, err := range uidIndex.Find(ctx, "linka.cloud.test.Test", filter, FindOptions{Limit: 100, Reverse: true}) {
				if err != nil {
					b.Fatal(err)
				}
				benchUIDSink = uid
				count++
			}
			benchKeysSink = count
		}
	})
}

func BenchmarkUIDIndexFind1M(b *testing.B) {
	benchmarkUIDIndexFindWithOptions1M(b, FindOptions{})
}

func BenchmarkUIDIndexFindLimit1001M(b *testing.B) {
	benchmarkUIDIndexFindWithOptions1M(b, FindOptions{Limit: 100})
}

func BenchmarkUIDIndexFindReverse1M(b *testing.B) {
	benchmarkUIDIndexFindWithOptions1M(b, FindOptions{Reverse: true})
}

func BenchmarkUIDIndexFindOffset10000Limit1001M(b *testing.B) {
	benchmarkUIDIndexFindWithOptions1M(b, FindOptions{Offset: 10_000, Limit: 100})
}

func benchmarkUIDIndexFindWithOptions1M(b *testing.B, opts FindOptions) {
	ctx := context.Background()
	const total = 1_000_000
	const matchEvery = 10

	filter := filters.Where("string_field").StringEquals("match")
	uidIndex := benchmarkBuildUIDIndex(b, ctx, total, matchEvery)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		for uid, err := range uidIndex.Find(ctx, "linka.cloud.test.Test", filter, opts) {
			if err != nil {
				b.Fatal(err)
			}
			benchUIDSink = uid
			count++
		}
		benchKeysSink = count
	}
}

func benchmarkBuildKeyIndex(b *testing.B, ctx context.Context, total, matchEvery int) Index {
	b.Helper()
	i := New(nil, All)
	for n := 0; n < total; n++ {
		msg := &test.Test{StringField: "other"}
		if n%matchEvery == 0 {
			msg.StringField = "match"
		}
		if err := i.Insert(ctx, strconv.Itoa(n+1), msg); err != nil {
			b.Fatal(err)
		}
	}
	return i
}

func benchmarkBuildUIDIndex(b *testing.B, ctx context.Context, total, matchEvery int) UIDIndex {
	b.Helper()
	i := NewUID(nil, All)
	for n := 0; n < total; n++ {
		msg := &test.Test{StringField: "other"}
		if n%matchEvery == 0 {
			msg.StringField = "match"
		}
		if err := i.Insert(ctx, uint64(n+1), msg); err != nil {
			b.Fatal(err)
		}
	}
	return i
}
