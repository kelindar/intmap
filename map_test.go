// Copyright (c) 2021-2024, Roman Atachiants
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
BenchmarkStore/intmap-24         	103890154	        10.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkStore/sync-24           	40663836	        26.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkStore/stdmap-24         	32448744	        36.67 ns/op	       0 B/op	       0 allocs/op
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
BenchmarkLoad/intmap-0%-24         	20124705	        59.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-0%-24           	19219896	        62.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-0%-24         	66543930	        16.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-10%-24        	20681667	        56.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-10%-24          	19684521	        59.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-10%-24        	53632722	        21.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-50%-24        	31344277	        37.56 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-50%-24          	28165971	        41.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-50%-24        	43359806	        26.15 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-90%-24        	93809362	        11.31 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-90%-24          	61292348	        17.59 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-90%-24        	38287116	        30.59 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/intmap-100%-24       	156096963	         7.426 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/sync-100%-24         	63665881	        16.59 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoad/stdmap-100%-24       	35932876	        32.91 ns/op	       0 B/op	       0 allocs/op
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
				our.Load(rand.Uint32N(count) + shift)
			}
		})

		b.Run(fmt.Sprintf("sync-%v%%", rate), func(b *testing.B) {
			shift := uint32(count - count*rate/100)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				syn.Load(rand.Uint32N(count) + shift)
			}
		})

		b.Run(fmt.Sprintf("stdmap-%v%%", rate), func(b *testing.B) {
			shift := uint32(count - count*rate/100)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = std[rand.Uint32N(count)+shift]
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
				if v, ok = m.Load(j); ok {
					t.Errorf("expected 'not found' flag for %d, found %d", j, v)
				}
			}
		}

		if v, ok = m.Load(0); !ok || v != 12345 {
			t.Errorf("expected 12345 for key 0")
		}

	}
}

func TestCapacity(t *testing.T) {
	m := New(10, 0.6)
	assert.Equal(t, 32, m.Capacity())
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

func TestArraySize(t *testing.T) {
	assert.Equal(t, 16, arraySize(10, .99))
	assert.Equal(t, 8, arraySize(0, .99))
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
		m.Store(rand.Uint32(), uint32(i))
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

func TestRangeEach(t *testing.T) {
	m := New(10, 0.6)
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
	m := New(10, 0.6)
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
	m := New(10, 0.6)
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
	m := New(10, 0.6)
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
	m := New(10, 0.6)
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
