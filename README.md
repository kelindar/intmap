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
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
BenchmarkStore/intmap-8                 14596622                83.30 ns/op            0 B/op          0 allocs/op
BenchmarkStore/stdmap-8                 10043176               116.1 ns/op             0 B/op          0 allocs/op
BenchmarkLoad/intmap-0%-8               10720995               108.1 ns/op             0 B/op          0 allocs/op
BenchmarkLoad/stdmap-0%-8               10443445               110.4 ns/op             0 B/op          0 allocs/op
BenchmarkLoad/intmap-10%-8              10510770               105.3 ns/op             0 B/op          0 allocs/op
BenchmarkLoad/stdmap-10%-8              10319277               109.3 ns/op             0 B/op          0 allocs/op
BenchmarkLoad/intmap-50%-8              12594564                89.33 ns/op            0 B/op          0 allocs/op
BenchmarkLoad/stdmap-50%-8              10215164               112.6 ns/op             0 B/op          0 allocs/op
BenchmarkLoad/intmap-90%-8              19359448                64.39 ns/op            0 B/op          0 allocs/op
BenchmarkLoad/stdmap-90%-8              10160139               111.7 ns/op             0 B/op          0 allocs/op
BenchmarkLoad/intmap-100%-8             20707099                62.17 ns/op            0 B/op          0 allocs/op
BenchmarkLoad/stdmap-100%-8              9713601               110.4 ns/op             0 B/op          0 allocs/op
```
