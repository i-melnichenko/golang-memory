package main

import (
	"fmt"
	"testing"
)

// --- Smaug type generator (same as before) ---

func makeSmaugType(size int) interface{} {
	switch size {
	case 1 * 1024:
		return struct{ Fire [1 * 1024]byte }{}
	case 2 * 1024:
		return struct{ Fire [2 * 1024]byte }{}
	case 4 * 1024:
		return struct{ Fire [4 * 1024]byte }{}
	case 8 * 1024:
		return struct{ Fire [8 * 1024]byte }{}
	case 16 * 1024:
		return struct{ Fire [16 * 1024]byte }{}
	case 32 * 1024:
		return struct{ Fire [32 * 1024]byte }{}
	case 64 * 1024:
		return struct{ Fire [64 * 1024]byte }{}
	case 128 * 1024:
		return struct{ Fire [128 * 1024]byte }{}
	case 256 * 1024:
		return struct{ Fire [256 * 1024]byte }{}
	case 512 * 1024:
		return struct{ Fire [512 * 1024]byte }{}
	case 1024 * 1024:
		return struct{ Fire [1024 * 1024]byte }{}
	default:
		panic("not implemented")
	}
}

// ===============================================================
//  STRICT STACK COPY (value return across frame boundary)
// ===============================================================

//go:noinline
func stackCopy[T any](x T) T {
	return x // forces real, full-size copy on return
}

// ===============================================================
//  STRICT HEAP ALLOCATION
// ===============================================================

//go:noinline
func heapAlloc[T any](x T) *T {
	y := x    // one stack copy (real) but cheap
	return &y // escape → allocate full copy in heap
}

// ===============================================================
//  BENCHMARK GENERATOR
// ===============================================================

func BenchmarkSmaugScalingStrict(b *testing.B) {
	sizes := []int{
		1 * 1024,
		2 * 1024,
		4 * 1024,
		8 * 1024,
		16 * 1024,
		32 * 1024,
		64 * 1024,
		128 * 1024,
		256 * 1024,
		512 * 1024,
		1024 * 1024,
	}

	for _, size := range sizes {
		s := makeSmaugType(size)

		// ---------- STACK (value copy) ----------
		b.Run("StackStrict_"+formatSize(size), func(b *testing.B) {
			switch v := s.(type) {
			case struct{ Fire [1 * 1024]byte }:
				runStack(b, v)
			case struct{ Fire [2 * 1024]byte }:
				runStack(b, v)
			case struct{ Fire [4 * 1024]byte }:
				runStack(b, v)
			case struct{ Fire [8 * 1024]byte }:
				runStack(b, v)
			case struct{ Fire [16 * 1024]byte }:
				runStack(b, v)
			case struct{ Fire [32 * 1024]byte }:
				runStack(b, v)
			case struct{ Fire [64 * 1024]byte }:
				runStack(b, v)
			case struct{ Fire [128 * 1024]byte }:
				runStack(b, v)
			case struct{ Fire [256 * 1024]byte }:
				runStack(b, v)
			case struct{ Fire [512 * 1024]byte }:
				runStack(b, v)
			case struct{ Fire [1024 * 1024]byte }:
				runStack(b, v)
			}
		})

		// ---------- HEAP (pointer return) ----------
		b.Run("HeapStrict_"+formatSize(size), func(b *testing.B) {
			switch v := s.(type) {
			case struct{ Fire [1 * 1024]byte }:
				runHeap(b, v)
			case struct{ Fire [2 * 1024]byte }:
				runHeap(b, v)
			case struct{ Fire [4 * 1024]byte }:
				runHeap(b, v)
			case struct{ Fire [8 * 1024]byte }:
				runHeap(b, v)
			case struct{ Fire [16 * 1024]byte }:
				runHeap(b, v)
			case struct{ Fire [32 * 1024]byte }:
				runHeap(b, v)
			case struct{ Fire [64 * 1024]byte }:
				runHeap(b, v)
			case struct{ Fire [128 * 1024]byte }:
				runHeap(b, v)
			case struct{ Fire [256 * 1024]byte }:
				runHeap(b, v)
			case struct{ Fire [512 * 1024]byte }:
				runHeap(b, v)
			case struct{ Fire [1024 * 1024]byte }:
				runHeap(b, v)
			}
		})
	}
}

func runStack[T any](b *testing.B, zero T) {
	var sink T
	for i := 0; i < b.N; i++ {
		sink = stackCopy(zero)
	}
	_ = sink
}

func runHeap[T any](b *testing.B, zero T) {
	var sink *T
	for i := 0; i < b.N; i++ {
		sink = heapAlloc(zero)
	}
	_ = sink
}

func formatSize(n int) string {
	if n >= 1024*1024 {
		return "1MB"
	}
	return fmt.Sprintf("%dKB", n/1024)
}
