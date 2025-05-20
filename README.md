<p align="center">
<img width="330" height="110" src=".github/logo.png" border="0" alt="kelindar/intmap">
<br>
<img src="https://img.shields.io/github/go-mod/go-version/kelindar/intmap" alt="Go Version">
<a href="https://pkg.go.dev/github.com/kelindar/intmap"><img src="https://pkg.go.dev/badge/github.com/kelindar/intmap" alt="PkgGoDev"></a>
<a href="https://goreportcard.com/report/github.com/kelindar/intmap"><img src="https://goreportcard.com/badge/github.com/kelindar/intmap" alt="Go Report Card"></a>
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License"></a>
<a href="https://coveralls.io/github/kelindar/intmap"><img src="https://coveralls.io/repos/github/kelindar/intmap/badge.svg" alt="Coverage"></a>
</p>

# Uint32-to-Uint32 Map

This repository contains an implementation of `uint32-to-uint32` map which is **~20-50%** faster than Go standard map for the same types (see benchmarks below). The code was based of [Brent Pedersen's intintmap](https://github.com/brentp/intintmap) and the [main logic](http://java-performance.info/implementing-world-fastest-java-int-to-int-hash-map/) remains intact, with some bug fixes and improvements of the API itself. The map is backed by a single array which interleaves keys and values to improve data locality.

## Usage

```go
// Create a new map with capacity of 1024 (resizeable) and 90% desired fill rate
m := intmap.New(1024, 0.90)

// Store a few key/value pairs
m.Store(1, 100)
m.Store(2, 200)

// Load them
v, ok := m.Load(1)
v, ok := m.Load(2)

// Delete keys
m.Delete(1)
m.Delete(2)
```

## Benchmarks

Looking at the benchmarks agains the standard Go map, this map should perform roughly 20-50% better depending on the conditions.

```
cpu: 13th Gen Intel(R) Core(TM) i7-13700K
BenchmarkStore/intmap-24         	144682518	        8.203 ns/op	       0 B/op	       0 allocs/op
BenchmarkStore/sync-24           	47746893	        24.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkStore/stdmap-24         	65009140	        17.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0%-24         	40558078	        28.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-0%-24           	36516896	        31.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0%-24         	76847426	        15.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-10%-24        	40846196	        28.31 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-10%-24          	37622625	        31.21 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-10%-24        	69917145	        16.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-50%-24        	52405636	        21.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-50%-24          	44567251	        25.75 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-50%-24        	50478930	        23.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-90%-24        	140594277	        8.486 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-90%-24          	83630574	        14.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-90%-24        	69522844	        17.17 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-100%-24       	189147504	        6.271 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-100%-24         	88044195	        13.58 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-100%-24       	78736231	        15.12 ns/op	       0 B/op	       0 allocs/op
```
