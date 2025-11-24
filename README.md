# 🧙‍♂️ Go Memory Benchmarks — Stack vs Heap

## Bilbo & Smaug Study (Go 1.25.4)

This repository contains a set of microbenchmarks that explore how Go allocates and copies data on:

* **stack**
* **heap**
* across function-call boundaries
* with escape analysis disabled
* on different ARM64 CPUs

The goal is to deeply understand how **value copying**, **heap allocation**, and **memory bandwidth** behave in Go.

---

# ▶ How to Run

These benchmarks were executed using:

> **Go version:** `go1.25.4`

To run all benchmarks with allocation statistics:

```bash
go test -bench=. -benchmem
```

---

# 🧩 Overview of Benchmark Types

### **Bilbo**

A tiny struct (small, cheap to copy).
Used to demonstrate stack vs heap behavior for small objects.

### **SmaugScalingStrict**

A massive struct (1KB → 1MB).
Copying forced via:

* `//go:noinline`
* no inlining
* frame-boundary passing
* generics keeping type stability
* explicit heap escape

This produces **honest** stack memcpy and **real** heap allocation.

---

# 🧪 Results on Apple M1 Max (ARM64)

Machine:

> **Apple M1 Max**
> **goos: darwin**
> **goarch: arm64**

## Bilbo Benchmark

```
BenchmarkBilboStack-10     0.9928 ns/op    0 B/op    0 allocs/op
BenchmarkBilboHeap-10     12.16 ns/op     16 B/op   1 allocs/op
```

**Stack is ~12× faster.**

---

# 🐉 SmaugScalingStrict (1KB → 1MB)

## StackStrict (true memcpy across frames)

| Size  | Time (ns) |
| ----- | --------- |
| 1KB   | 78.36 ns  |
| 2KB   | 159.3 ns  |
| 4KB   | 379.5 ns  |
| 8KB   | 559.6 ns  |
| 16KB  | 1093 ns   |
| 32KB  | 2170 ns   |
| 64KB  | 4555 ns   |
| 128KB | 9562 ns   |
| 256KB | 25385 ns  |
| 512KB | 51753 ns  |
| 1MB   | 95518 ns  |

### Observations:

* Perfect linear scaling
* memcpy throughput ≈ **10–12 GB/s**
* Stack remains efficient even at 1MB

---

## HeapStrict (real heap allocation)

| Size  | Time (ns) | B/op        |
| ----- | --------- | ----------- |
| 1KB   | 207.6 ns  | 1024 B      |
| 2KB   | 378.7 ns  | 2048 B      |
| 4KB   | 783.2 ns  | 4096 B      |
| 8KB   | 1472 ns   | 8192 B      |
| 16KB  | 2650 ns   | 16384 B     |
| 32KB  | 4091 ns   | 32768 B     |
| 64KB  | 8818 ns   | 65536 B     |
| 128KB | 15704 ns  | 131072 B    |
| 256KB | 31312 ns  | 262152 B    |
| 512KB | 53725 ns  | 524311 B    |
| 1MB   | 248863 ns | 1,048,801 B |

## Final comparison — M1 Max

| Size  | StackStrict | HeapStrict | Faster        |
| ----- | ----------- | ---------- | ------------- |
| 1KB   | 78 ns       | 207 ns     | Stack (~2.6×) |
| 4KB   | 379 ns      | 783 ns     | Stack (~2×)   |
| 64KB  | 4555 ns     | 8818 ns    | Stack (~2×)   |
| 256KB | 25385 ns    | 31312 ns   | Stack (~1.2×) |
| 1MB   | 95518 ns    | 248863 ns  | Stack (~2.6×) |

> **Stack wins at all sizes on M1 Max.**

---

# 🐢 Results on Raspberry Pi 4 (ARM Cortex-A72)

Machine:

> **Raspberry Pi 4**
> **goos: linux**
> **goarch: arm64**

Raspberry Pi 4 has:

* low memory bandwidth (~3–5 GB/s)
* small L1/L2 caches
* slow heap allocator
* weak SIMD
* 1.5GHz ARM cores

This makes heap allocation dramatically slower.

---

# Bilbo Benchmark (RPi4)

```
BenchmarkBilboStack-4     4.329 ns/op
BenchmarkBilboHeap-4    128.0 ns/op
```

> **Stack is ~30× faster.**

---

# SmaugScalingStrict (RPi4)

## StackStrict (real memcpy)

| Size  | Time (ns)  |
| ----- | ---------- |
| 1KB   | 312.9 ns   |
| 2KB   | 1004 ns    |
| 4KB   | 2056 ns    |
| 8KB   | 4010 ns    |
| 16KB  | 8356 ns    |
| 32KB  | 16979 ns   |
| 64KB  | 33483 ns   |
| 128KB | 73383 ns   |
| 256KB | 259084 ns  |
| 512KB | 1087095 ns |
| 1MB   | 2860634 ns |

### Notes:

* Linear scaling
* MUCH slower than M1 Max (≈20–30×)
* memcpy throughput ≈ **350–400 MB/s**

---

## HeapStrict (real heap)

| Size  | Time (ns)   | B/op        |
| ----- | ----------- | ----------- |
| 1KB   | 3624 ns     | 1024 B      |
| 2KB   | 5141 ns     | 2048 B      |
| 4KB   | 7662 ns     | 4096 B      |
| 8KB   | 16647 ns    | 8192 B      |
| 16KB  | 32821 ns    | 16384 B     |
| 32KB  | 52303 ns    | 32768 B     |
| 64KB  | 104008 ns   | 65536 B     |
| 128KB | 237656 ns   | 131072 B    |
| 256KB | 4244385 ns  | 263227 B    |
| 512KB | 22560989 ns | 529532 B    |
| 1MB   | 59123750 ns | 1,059,063 B |

## Final comparison — Raspberry Pi 4

| Size  | StackStrict | HeapStrict  | Faster       |
| ----- | ----------- | ----------- | ------------ |
| 1KB   | 313 ns      | 3624 ns     | Stack (~11×) |
| 4KB   | 2056 ns     | 7662 ns     | Stack (~4×)  |
| 64KB  | 33483 ns    | 104008 ns   | Stack (~3×)  |
| 256KB | 259084 ns   | 4244385 ns  | Stack (~16×) |
| 1MB   | 2860634 ns  | 59123750 ns | Stack (~20×) |

> **Heap is catastrophically slow on Raspberry Pi 4** due to slow zeroing and allocator overhead.

---

# 🧠 Final Engineering Conclusions

### ✔ Small structs (Bilbo)

Always use stack/value return — fastest on all architectures.

### ✔ Medium structs (1–64KB)

Stack is consistently 2×–10× faster.

### ✔ Large structs (128KB–1MB)

Stack memcpy remains predictable and faster even on weak hardware.

### ✔ Heap is slower everywhere

But **on Raspberry Pi 4 it's 20× slower** due to:

* slow zeroing
* small caches
* weak SIMD
* slow runtime allocator

### ✔ Apple M1 Max is a memory monster

* 10–12 GB/s memcpy
* heap performance much better
* sometimes heap can compete when stack does *multiple* copies

---

## ⚠ Important Note: Single-Threaded Benchmarks

All benchmarks in this repository are **single-threaded** by design.
Each test runs in a single goroutine and measures:

* pure stack value copying
* pure heap allocation
* raw memory bandwidth
* function call boundaries
* compiler behavior (`//go:noinline`, escape analysis, ABI returns)

This setup is intentional, because it isolates the effects we want to observe.

However, real Go applications often create hundreds or thousands of goroutines, and under such conditions:

* the heap allocator behaves differently
* per-P caches (`mcache`) become more active
* contention on shared allocator structures (`mcentral`, `mheap`) may appear
* GC write barriers fire more frequently
* memory locality changes dramatically
* stack growth/shrink operations may occur

In multi-goroutine or highly parallel programs, performance characteristics can shift.
Heap allocations that appear expensive in microbenchmarks may be partially amortized by concurrent allocators, while stack copying may interact differently with CPU caches under load.

That said, **let’s be honest**:
in the vast majority of real-world Go workloads — easily **90%+** — the critical execution paths still run *logically single-threaded*, even if wrapped in goroutines. Most tasks process data sequentially, or with minimal parallelism, and therefore the single-threaded behavior measured here is highly relevant for everyday engineering work.
