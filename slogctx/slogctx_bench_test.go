package slogctx

import (
	"context"
	"testing"
)

func BenchmarkWithValues_10(b *testing.B) {
	ctx := context.Background()
	args := []any{}
	for i := 0; i < 10; i++ {
		args = append(args, "k"+string(rune(i)), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WithValues(ctx, args...)
	}
}

func BenchmarkWithUniqueValues_10(b *testing.B) {
	ctx := context.Background()
	args := []any{}
	for i := 0; i < 10; i++ {
		args = append(args, "k"+string(rune(i)), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WithUniqueValues(ctx, args...)
	}
}

func BenchmarkWithoutKeys_10of100(b *testing.B) {
	ctx := context.Background()
	args := []any{}
	for i := 0; i < 100; i++ {
		args = append(args, "k"+string(rune(i)), i)
	}
	ctx = WithValues(ctx, args...)
	remove := []string{}
	for i := 0; i < 10; i++ {
		remove = append(remove, "k"+string(rune(i)))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WithoutKeys(ctx, remove...)
	}
}

func BenchmarkWithoutAllKeys_100(b *testing.B) {
	ctx := context.Background()
	args := []any{}
	for i := 0; i < 100; i++ {
		args = append(args, "k"+string(rune(i)), i)
	}
	ctx = WithValues(ctx, args...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WithoutAllKeys(ctx)
	}
}
