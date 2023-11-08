// Copyright (c) 2021-2023, Roman Atachiants
// Copyright (c) 2016, Brent Pedersen - Bioinformatics

package intmap

import (
	"fmt"
	"hash/crc32"
	"math"
	"testing"

	"github.com/kelindar/xxrand"
	"github.com/stretchr/testify/assert"
)

/*
cpu: 13th Gen Intel(R) Core(TM) i7-13700K
BenchmarkStore/intmap-24         	18568460	        61.14 ns/op	       0 B/op	       0 allocs/op
BenchmarkStore/sync-24           	16417874	        71.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkStore/stdmap-24         	14243356	        73.32 ns/op	       0 B/op	       0 allocs/op
*/
func BenchmarkStore(b *testing.B) {
	const count = 1000000
	our := New(count, .90)
	syn := NewSync(count, .90)
	std := make(map[uint32]uint32, count)

	b.Run("intmap", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			our.Store(xxrand.Uint32n(count), 1)
		}
	})

	b.Run("sync", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			syn.Store(xxrand.Uint32n(count), 1)
		}
	})

	b.Run("stdmap", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			std[xxrand.Uint32n(count)] = 1
		}
	})
}

/*
cpu: 13th Gen Intel(R) Core(TM) i7-13700K
BenchmarkLoad/intmap-0%-24         	13987644	        83.19 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-0%-24           	13353600	        88.65 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0%-24         	16724271	        66.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-10%-24        	14666626	        79.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-10%-24          	14067055	        84.99 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-10%-24        	17432941	        65.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-50%-24        	18002366	        66.59 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-50%-24          	16336176	        71.54 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-50%-24        	16711392	        66.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-90%-24        	23374497	        47.08 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-90%-24          	20600468	        51.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-90%-24        	17858258	        65.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-100%-24       	26168407	        45.27 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-100%-24         	20380226	        49.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-100%-24       	16759916	        66.16 ns/op	       0 B/op	       0 allocs/op
*/
func BenchmarkLoad(b *testing.B) {
	const count = 1000000
	our := sequentialMap(count)
	syn := sequentialSyncMap(count)
	std := make(map[uint32]uint32, count)
	for i := uint32(0); i < count; i++ {
		std[i] = i
	}

	for _, rate := range []float64{0, 10, 50, 90, 100} {
		rate := rate
		b.Run(fmt.Sprintf("intmap-%v%%", rate), func(b *testing.B) {
			shift := uint32(count - count*rate/100)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				our.Load(xxrand.Uint32n(count) + shift)
			}
		})

		b.Run(fmt.Sprintf("sync-%v%%", rate), func(b *testing.B) {
			shift := uint32(count - count*rate/100)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				syn.Load(xxrand.Uint32n(count) + shift)
			}
		})

		b.Run(fmt.Sprintf("stdmap-%v%%", rate), func(b *testing.B) {
			shift := uint32(count - count*rate/100)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = std[xxrand.Uint32n(count)+shift]
			}
		})
	}
}

func TestInvalidNew(t *testing.T) {
	assert.Panics(t, func() {
		New(10, 0)
	})

	assert.Panics(t, func() {
		New(0, .99)
	})
}

func TestMapSimple(t *testing.T) {
	m := New(10, 0.99)
	var i uint32
	var v uint32
	var ok bool

	// --------------------------------------------------------------------
	// Put() and Get()

	for i = 0; i < 20000; i += 2 {
		m.Store(i, i)
	}
	for i = 0; i < 20000; i += 2 {
		if v, ok = m.Load(i); !ok || v != i {
			t.Errorf("didn't get expected value")
		}
		if _, ok = m.Load(i + 1); ok {
			t.Errorf("didn't get expected 'not found' flag")
		}
	}

	if m.Count() != int(20000/2) {
		t.Errorf("size (%d) is not right, should be %d", m.Count(), int(20000/2))
	}
	// --------------------------------------------------------------------
	// Del()

	for i = 0; i < 20000; i += 2 {
		m.Delete(i)
	}
	for i = 0; i < 20000; i += 2 {
		if _, ok = m.Load(i); ok {
			t.Errorf("didn't get expected 'not found' flag")
		}
		if _, ok = m.Load(i + 1); ok {
			t.Errorf("didn't get expected 'not found' flag")
		}
	}

	// --------------------------------------------------------------------
	// Put() and Get()

	for i = 0; i < 20000; i += 2 {
		m.Store(i, i*2)
	}
	for i = 0; i < 20000; i += 2 {
		if v, ok = m.Load(i); !ok || v != i*2 {
			t.Errorf("didn't get expected value")
		}
		if _, ok = m.Load(i + 1); ok {
			t.Errorf("didn't get expected 'not found' flag")
		}
	}

}

func TestMap(t *testing.T) {
	m := New(10, 0.6)
	var ok bool
	var v uint32

	step := uint32(61)

	var i uint32
	m.Store(0, 12345)
	for i = 1; i < 100000000; i += step {
		m.Store(i, i+7)
		m.Store(-i, i-7)

		if v, ok = m.Load(i); !ok || v != i+7 {
			t.Errorf("expected %d as value for key %d, got %d", i+7, i, v)
		}
		if v, ok = m.Load(-i); !ok || v != i-7 {
			t.Errorf("expected %d as value for key %d, got %d", i-7, -i, v)
		}
	}
	for i = 1; i < 100000000; i += step {
		if v, ok = m.Load(i); !ok || v != i+7 {
			t.Errorf("expected %d as value for key %d, got %d", i+7, i, v)
		}
		if v, ok = m.Load(-i); !ok || v != i-7 {
			t.Errorf("expected %d as value for key %d, got %d", i-7, -i, v)
		}

		for j := i + 1; j < i+step; j++ {
			if v, ok = m.Load(j); ok {
				t.Errorf("expected 'not found' flag for %d, found %d", j, v)
			}
		}
	}

	if v, ok = m.Load(0); !ok || v != 12345 {
		t.Errorf("expected 12345 for key 0")
	}
}

func TestDeleteSequential(t *testing.T) {
	const size = 100
	m := sequentialMap(size)

	// Try to delete multiple times
	for retry := 0; retry < 3; retry++ {
		for i := 0; i < size; i += 2 {
			m.Delete(uint32(i))
		}
		assert.Equal(t, size/2, m.Count())
	}

	// Now delete the rest
	for i := 1; i < size; i += 2 {
		m.Delete(uint32(i))
	}
	assert.Equal(t, 0, m.Count())
}

func TestDeleteRandom(t *testing.T) {
	const size = 1000000
	m := randomMap(size)

	// Try to delete multiple times
	for retry := 0; retry < 3; retry++ {
		i := 0
		m.Range(func(key, value uint32) bool {
			if i++; i%2 == 0 {
				m.Delete(key)
			}
			return true
		})
	}

	// Delete the rest
	i := 0
	m.Range(func(key, value uint32) bool {
		if i++; i%2 == 1 {
			m.Delete(key)
		}
		return true
	})
}

func TestRangeSequential(t *testing.T) {
	for _, size := range []int{100, 10000, 1000000} {
		m := New(size, 0.99)
		expect := 0
		for i := 0; i < size; i++ {
			m.Store(uint32(i), uint32(i))
			expect += i
		}

		// Range and check if sum is the same
		sum := 0
		m.Range(func(key, value uint32) bool {
			sum += int(key)
			return true
		})
		assert.Equal(t, expect, sum)
	}
}

func TestRangeRandom(t *testing.T) {
	for _, size := range []int{100, 10000, 1000000} {
		count := 0
		m := randomMap(size)
		m.Range(func(key, value uint32) bool {
			count++
			return true
		})
		assert.Equal(t, m.Count(), count)
	}
}

func TestCapacityFor(t *testing.T) {
	assert.Equal(t, uint32(0x1), capacityFor(0))
	assert.Equal(t, uint32(0xffffffff), capacityFor(math.MaxUint32))
	assert.Equal(t, uint32(0x10), capacityFor(10))
}

func TestArraySize(t *testing.T) {
	assert.Equal(t, 16, arraySize(10, .99))
	assert.Equal(t, 2, arraySize(0, .99))
}

func TestSequentialCollisions(t *testing.T) {
	for _, size := range []int{1e4, 1e5, 1e6} {
		avg, max := collisionRate(size, func(i uint32) uint32 {
			return i
		})
		assert.LessOrEqual(t, avg, 2.0)
		assert.LessOrEqual(t, max, 10)
	}
}

func TestRandomCollisions(t *testing.T) {
	for _, size := range []int{100, 10000, 1000000} {
		avg, max := collisionRate(size, func(i uint32) uint32 {
			return xxrand.Uint32()
		})
		assert.LessOrEqual(t, avg, 2.0)
		assert.LessOrEqual(t, max, 10)
	}
}

func TestStringCollisions(t *testing.T) {
	for _, size := range []int{100, 10000, 1000000} {
		avg, max := collisionRate(size, func(i uint32) uint32 {
			return crc32.ChecksumIEEE([]byte(fmt.Sprintf("value of %x", i)))
		})
		assert.LessOrEqual(t, avg, 2.0)
		assert.LessOrEqual(t, max, 10)
	}
}

func collisionRate(count int, next func(i uint32) uint32) (avg float64, max int) {
	counts := make(map[uint32]int, count)
	mask := capacityFor(uint32(count)) - 1
	for i := 0; i < count; i++ {
		offset := bucketOf(next(uint32(i)), mask)
		counts[offset] += 1
	}

	sum, n := .0, .0
	for _, v := range counts {
		sum += float64(v)
		n++
	}

	return sum / n, max
}

// sequentialMap creates a new map with sequential keys
func sequentialMap(size int) *Map {
	m := New(size, 0.99)
	for i := 0; i < size; i++ {
		m.Store(uint32(i), uint32(i))
	}
	return m
}

// randomMap creates a new map with random keys
func randomMap(size int) *Map {
	m := New(size, 0.99)
	for i := 0; i < size; i++ {
		m.Store(xxrand.Uint32(), uint32(i))
	}
	return m
}

// sequentialSyncMap creates a new map with sequential keys
func sequentialSyncMap(size int) *Sync {
	m := NewSync(size, 0.99)
	for i := 0; i < size; i++ {
		m.Store(uint32(i), uint32(i))
	}
	return m
}

func TestMapClone(t *testing.T) {
	original := New(10, 0.6)
	original.Store(1, 10)
	original.Store(2, 20)
	original.Store(3, 30)

	clone := original.Clone()

	// Check that the clone is not the same object as the original
	assert.NotEqual(t, clone, original, "clone and original are the same object")

	// Check that the clone has the same count
	assert.Equal(t, original.Count(), clone.Count(), "clone count does not match original count")

	// Check that the clone has the same data
	for i := uint32(1); i <= 3; i++ {
		v1, ok1 := original.Load(i)
		v2, ok2 := clone.Load(i)
		assert.True(t, ok1, "original does not have key %d", i)
		assert.True(t, ok2, "clone does not have key %d", i)
		assert.Equal(t, v1, v2, "clone does not have the same data as the original")
	}

	// Check that modifying the clone does not modify the original
	clone.Store(4, 40)
	_, ok := original.Load(4)
	assert.False(t, ok, "modifying clone modified the original")
}
