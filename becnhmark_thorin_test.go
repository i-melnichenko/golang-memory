package main

import (
	"fmt"
	"testing"
)

// newThorin creates a Thorin struct with a specific payload size.
// This allows us to benchmark stack vs heap behavior across multiple sizes.
type ThorinDynamic struct {
	Payload []byte
}

// -----------------------------------------------------------
// STACK STRICT: VALUE RETURN (forces full memcpy)
// -----------------------------------------------------------

//go:noinline
func repoThorinStack(size int) ThorinDynamic {
	t := ThorinDynamic{Payload: make([]byte, size)}
	t.Payload[0] = 1
	return t // FULL COPY
}

//go:noinline
func serviceThorinStack(size int) ThorinDynamic {
	t := repoThorinStack(size) // COPY #2
	t.Payload[1] = 2
	return t // COPY #3
}

//go:noinline
func transportThorinStack(size int) ThorinDynamic {
	t := serviceThorinStack(size) // COPY #4
	t.Payload[2] = 3
	return t // COPY #5
}

// -----------------------------------------------------------
// HEAP STRICT: POINTER RETURN (single allocation, no copies)
// -----------------------------------------------------------

//go:noinline
func repoThorinHeap(size int) *ThorinDynamic {
	t := &ThorinDynamic{Payload: make([]byte, size)} // ALLOC #1
	t.Payload[0] = 1
	return t
}

//go:noinline
func serviceThorinHeap(size int) *ThorinDynamic {
	t := repoThorinHeap(size) // NO COPY
	t.Payload[1] = 2
	return t
}

//go:noinline
func transportThorinHeap(size int) *ThorinDynamic {
	t := serviceThorinHeap(size) // NO COPY
	t.Payload[2] = 3
	return t
}

// -----------------------------------------------------------
// BENCHMARKS
// -----------------------------------------------------------

var thorinSizes = []int{
	1 << 10,   // 1 KB
	2 << 10,   // 2 KB
	4 << 10,   // 4 KB
	8 << 10,   // 8 KB
	16 << 10,  // 16 KB
	32 << 10,  // 32 KB
	64 << 10,  // 64 KB
	128 << 10, // 128 KB
	256 << 10, // 256 KB
	512 << 10, // 512 KB
	1 << 20,   // 1 MB
}

func BenchmarkThorinScalingStrict(b *testing.B) {
	for _, size := range thorinSizes {

		// ----- STACK STRICT -----
		b.Run(
			"StackStrict_"+formatSizeThorin(size),
			func(b *testing.B) {
				var sink ThorinDynamic
				for i := 0; i < b.N; i++ {
					sink = transportThorinStack(size)
				}
				_ = sink
			},
		)

		// ----- HEAP STRICT -----
		b.Run(
			"HeapStrict_"+formatSizeThorin(size),
			func(b *testing.B) {
				var sink *ThorinDynamic
				for i := 0; i < b.N; i++ {
					sink = transportThorinHeap(size)
				}
				_ = sink
			},
		)
	}
}

// Helper to format “1KB”, “2KB”, …
func formatSizeThorin(sz int) string {
	switch {
	case sz >= 1<<20:
		return fmt.Sprintf("%dMB", sz>>20)
	case sz >= 1<<10:
		return fmt.Sprintf("%dKB", sz>>10)
	default:
		return fmt.Sprintf("%dB", sz)
	}
}
