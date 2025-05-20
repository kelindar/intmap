// Copyright (c) 2021-2025, Roman Atachiants
// Copyright (c) 2016, Brent Pedersen - Bioinformatics

package intmap

import (
	"fmt"
	"hash/crc32"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
cpu: 13th Gen Intel(R) Core(TM) i7-13700K
BenchmarkStore/intmap-24         	122662514	         9.357 ns/op	       0 B/op	       0 allocs/op
BenchmarkStore/sync-24           	45237855	        26.67 ns/op	       0 B/op	       0 allocs/op
BenchmarkStore/stdmap-24         	58779109	        22.53 ns/op	       0 B/op	       0 allocs/op
*/
func BenchmarkStore(b *testing.B) {
	const count = 1000000
	our := New(count)
	syn := NewSync(count)
	std := make(map[uint32]uint32, count)

	b.Run("intmap", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			our.Store(rand.Uint32N(count), 1)
		}
	})

	b.Run("sync", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			syn.Store(rand.Uint32N(count), 1)
		}
	})

	b.Run("stdmap", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			std[rand.Uint32N(count)] = 1
		}
	})
}

/*
cpu: 13th Gen Intel(R) Core(TM) i7-13700K
BenchmarkLoad/intmap-0.75-0%-24         	89028771	        12.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.75-0%-24         	49369306	        21.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.75-10%-24        	84331623	        14.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.75-10%-24        	47706888	        23.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.75-50%-24        	54766512	        19.19 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.75-50%-24        	36239652	        30.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.75-90%-24        	73981824	        14.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.75-90%-24        	43927723	        24.20 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.75-100%-24       	84131428	        14.02 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.75-100%-24       	45653761	        22.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.80-0%-24         	86299892	        14.02 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.80-0%-24         	54096052	        23.06 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.80-10%-24        	62889126	        16.08 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.80-10%-24        	44653488	        23.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.80-50%-24        	52184122	        19.35 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.80-50%-24        	36200295	        31.27 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.80-90%-24        	70501558	        15.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.80-90%-24        	46411608	        24.63 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.80-100%-24       	83554544	        14.07 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.80-100%-24       	53360131	        22.91 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.95-0%-24         	84525913	        14.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.95-0%-24         	48704263	        21.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.95-10%-24        	77747949	        15.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.95-10%-24        	43498275	        23.20 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.95-50%-24        	54157244	        20.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.95-50%-24        	34816112	        33.49 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.95-90%-24        	61174550	        17.20 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.95-90%-24        	46979418	        25.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.95-100%-24       	85645840	        14.39 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.95-100%-24       	47679972	        22.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.99-0%-24         	123137582	         9.767 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.99-0%-24         	54901183	        21.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.99-10%-24        	110715706	        10.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.99-10%-24        	45014460	        23.37 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.99-50%-24        	71354244	        14.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.99-50%-24        	37305776	        30.94 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.99-90%-24        	138160173	         8.867 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.99-90%-24        	40865253	        25.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-0.99-100%-24       	171428276	         7.208 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0.99-100%-24       	46678076	        22.50 ns/op	       0 B/op	       0 allocs/op
*/
func BenchmarkLoad(b *testing.B) {
	const count = 1000000

	for _, fill := range []float64{0.75, 0.80, 0.95, 0.99} {
		our := sequentialMap(count, fill)
		std := make(map[uint32]uint32, count)
		for i := uint32(0); i < count; i++ {
			std[i] = i
		}

		for _, rate := range []float64{0, 10, 50, 90, 100} {
			rate := rate
			b.Run(fmt.Sprintf("intmap-%.2f-%v%%", fill, rate), func(b *testing.B) {
				shift := uint32(count - count*rate/100)

				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					our.Load(rand.Uint32N(count) + shift)
				}
			})

			b.Run(fmt.Sprintf("stdmap-%.2f-%v%%", fill, rate), func(b *testing.B) {
				shift := uint32(count - count*rate/100)

				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = std[rand.Uint32N(count)+shift]
				}
			})
		}
	}
}

func TestInvalidNew(t *testing.T) {
	assert.Panics(t, func() {
		newMap(10, 1)
	})

	assert.Panics(t, func() {
		newMap(0, 0)
	})
}

func TestMapSimple(t *testing.T) {
	m := New(10)
	var i uint32
	var v uint32
	var ok bool

	// --------------------------------------------------------------------
	// Put() and Get()

	for i = 0; i < 20000; i += 2 {
		m.Store(i, i)
	}
	for i = 0; i < 20000; i += 2 {
		v, ok = m.Load(i)
		assert.True(t, ok, "expected key %d to be present", i)
		assert.Equal(t, i, v, "expected value for key %d to be %d, got %d", i, i, v)
		if _, ok = m.Load(i + 1); ok {
			assert.False(t, ok, "expected key %d to be absent", i+1)
		}
	}

	assert.Equal(t, int(20000/2), m.Count(), "size is not right, should be %d", int(20000/2))
	// --------------------------------------------------------------------
	// Del()

	for i = 0; i < 20000; i += 2 {
		m.Delete(i)
	}
	for i = 0; i < 20000; i += 2 {
		_, ok = m.Load(i)
		assert.False(t, ok, "expected key %d to be absent", i)
		if _, ok = m.Load(i + 1); ok {
			assert.False(t, ok, "expected key %d to be absent", i+1)
		}
	}

	// --------------------------------------------------------------------
	// Put() and Get()

	for i = 0; i < 20000; i += 2 {
		m.Store(i, i*2)
	}
	for i = 0; i < 20000; i += 2 {
		v, ok = m.Load(i)
		assert.True(t, ok, "expected key %d to be present", i)
		assert.Equal(t, i*2, v, "expected value for key %d to be %d, got %d", i, i*2, v)
		if _, ok = m.Load(i + 1); ok {
			t.Errorf("didn't get expected 'not found' flag")
		}
	}

}

func TestMap(t *testing.T) {
	m := New(10)
	var ok bool
	var v uint32
	var i uint32
	step := uint32(61)

	for retry := 0; retry < 3; retry++ {
		m.Clear()
		assert.Equal(t, 0, m.Count())

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
				v, ok = m.Load(j)
				assert.False(t, ok, "expected key %d to be absent", j)
			}
		}

		v, ok = m.Load(0)
		assert.True(t, ok, "expected key 0 to be present")
		assert.Equal(t, uint32(12345), v, "expected value for key 0 to be 12345, got %v", v)

	}
}

func TestCapacity(t *testing.T) {
	m := New(10)
	assert.Equal(t, 16, m.Capacity())
}

func TestDeleteSequential(t *testing.T) {
	const size = 100
	m := sequentialMap(size, defaultFillFactor)

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
		m := New(size)
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

func TestArraySize(t *testing.T) {
	assert.Equal(t, 16, arraySize(10, .99))
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
			return rand.Uint32()
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
	mask := arraySize(count, 1) - 1
	for i := 0; i < count; i++ {
		offset := bucketOf(next(uint32(i)), uint32(mask))
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
func sequentialMap(size int, fill float64) *Map {
	m := NewWithFill(size, fill)
	for i := 0; i < size; i++ {
		m.Store(uint32(i), uint32(i))
	}
	return m
}

// randomMap creates a new map with random keys
func randomMap(size int) *Map {
	m := New(size)
	for i := 0; i < size; i++ {
		m.Store(rand.Uint32(), uint32(i))
	}
	return m
}

// sequentialSyncMap creates a new map with sequential keys
func sequentialSyncMap(size int) *Sync {
	m := NewSync(size)
	for i := 0; i < size; i++ {
		m.Store(uint32(i), uint32(i))
	}
	return m
}

func TestMapClone(t *testing.T) {
	original := New(10)
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

func TestRangeEach(t *testing.T) {
	m := New(10)
	m.Store(isFree, 10)
	m.Store(2, 20)
	m.Store(3, 30)
	m.Store(4, 40)
	m.Store(5, 50)

	keys, values := []uint32{}, []uint32{}
	m.RangeEach(func(key, value uint32) {
		keys = append(keys, key)
		values = append(values, value)
	})

	assert.ElementsMatch(t, []uint32{0, 2, 3, 4, 5}, keys)
	assert.ElementsMatch(t, []uint32{10, 20, 30, 40, 50}, values)
}

func TestRangeErr(t *testing.T) {
	m := New(10)
	m.Store(1, 10)
	m.Store(2, 20)
	m.Store(3, 30)
	m.Store(4, 40)
	m.Store(5, 50)

	keys, values := []uint32{}, []uint32{}
	assert.NoError(t, m.RangeErr(func(key, value uint32) error {
		keys = append(keys, key)
		values = append(values, value)
		return nil
	}))

	assert.ElementsMatch(t, []uint32{1, 2, 3, 4, 5}, keys)
	assert.ElementsMatch(t, []uint32{10, 20, 30, 40, 50}, values)
}

func TestRangeErrStop(t *testing.T) {
	m := New(10)
	m.Store(1, 10)
	m.Store(2, 20)
	m.Store(3, 30)
	m.Store(4, 40)
	m.Store(5, 50)

	keys, values := []uint32{}, []uint32{}
	assert.EqualError(t, m.RangeErr(func(key, value uint32) error {
		keys = append(keys, key)
		values = append(values, value)
		return fmt.Errorf("stop")
	}), "stop")

	assert.Len(t, keys, 1)
	assert.Len(t, values, 1)
}

func TestRangeErrFreeKey(t *testing.T) {
	m := New(10)
	m.Store(isFree, 10)

	keys, values := []uint32{}, []uint32{}
	assert.NoError(t, m.RangeErr(func(key, value uint32) error {
		keys = append(keys, key)
		values = append(values, value)
		return nil
	}))

	assert.Len(t, keys, 1)
	assert.Len(t, values, 1)
}

func TestRangeStop(t *testing.T) {
	m := New(10)
	m.Store(0, 0)
	m.Store(1, 10)

	keys, values := []uint32{}, []uint32{}
	m.Range(func(key, value uint32) bool {
		keys = append(keys, key)
		values = append(values, value)
		return false
	})

	assert.Len(t, keys, 1)
	assert.Len(t, values, 1)
}
