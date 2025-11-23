package main

import "testing"

// Bilbo — a tiny hobbit-sized struct.
type Bilbo struct {
	Age   int
	Rings int
}

//go:noinline
func stackBilbo() Bilbo {
	// Bilbo lives on the stack (no heap escape).
	b := Bilbo{Age: 50, Rings: 0}
	return b
}

//go:noinline
func heapBilbo() *Bilbo {
	// Returning a pointer forces Bilbo to escape to the heap.
	b := &Bilbo{Age: 50, Rings: 0}
	return b
}

func BenchmarkBilboStack(b *testing.B) {
	var sink Bilbo
	for i := 0; i < b.N; i++ {
		sink = stackBilbo()
	}
	_ = sink // prevent optimization
}

func BenchmarkBilboHeap(b *testing.B) {
	var sink *Bilbo
	for i := 0; i < b.N; i++ {
		sink = heapBilbo()
	}
	_ = sink // prevent optimization
}
