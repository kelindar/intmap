// Copyright (c) 2021-2025 Roman Atachiants
// Copyright (c) 2016, Brent Pedersen - Bioinformatics

package intmap

import (
	"math"
)

const (
	isFree            = 0    // sentinel key marking an empty slot
	defaultFillFactor = 0.80 // default fill factor
)

// Map is a contiguous hash table with interleaved key/value slots.
type Map struct {
	data       []uint32  // [k0,v0,k1,v1,…]
	fillFactor float32   // max load factor before resize (e.g. 0.75)
	threshold  int32     // resize threshold (capacity*fillFactor)
	count      int32     // number of live entries (excl. free key)
	mask       [2]uint32 // mask[0] = bucket mask, mask[1] = slice index mask
	freeVal    uint32    // value for key == isFree
	hasFreeKey bool      // whether the free-key entry is present
}

// New allocates a map sized for at least `size` entries.
func New(size int) *Map {
	return newMap(size, defaultFillFactor)
}

// New allocates a map sized for at least `size` entries.
func NewWithFill(size int, fillFactor float64) *Map {
	return newMap(size, fillFactor)
}

// newMap allocates a map sized for at least `size` entries.
func newMap(size int, fillFactor float64) *Map {
	if fillFactor <= 0 || fillFactor >= 1 {
		panic("intmap: fill factor must be in (0,1)")
	}
	if size <= 0 {
		panic("intmap: size must be positive")
	}

	capSlots := arraySize(size, fillFactor)
	return &Map{
		data:       make([]uint32, capSlots*2),
		fillFactor: float32(fillFactor),
		threshold:  int32(math.Floor(float64(capSlots) * fillFactor)),
		mask:       [2]uint32{uint32(capSlots - 1), uint32(capSlots*2 - 1)},
	}
}

//go:nosplit
//go:inline
func bucketOf(key, mask uint32) uint32 {
	const phi32 = 0x9E3779B9         // 2^32 / golden-ratio
	return (key * phi32) & mask << 1 // 1 MUL, no XOR
}

// Capacity returns the maximum number of entries before resize.
func (m *Map) Capacity() int { return len(m.data) / 2 }

// Load returns the value stored in the map for a key, or nil if no value is
// present. The ok result indicates whether value was found in the map.
func (m *Map) Load(key uint32) (uint32, bool) {
	if key == isFree {
		if m.hasFreeKey {
			return m.freeVal, true
		}
		return 0, false
	}

	data := m.data
	mask := m.mask[0]
	mask1 := m.mask[1]
	ptr := bucketOf(key, mask) // starting slot
	dist := uint32(0)          // probe distance of seeker

	for {
		k := data[ptr]
		if k == key {
			return data[ptr+1], true // found
		}
		if k == isFree {
			return 0, false // hit gap – key absent
		}

		// displacement of occupant at ptr
		occDist := ((ptr - bucketOf(k, mask)) & mask1) >> 1
		if occDist < dist { // early exit – RH property
			return 0, false
		}

		ptr = (ptr + 2) & mask1
		dist++
	}
}

// Store sets the value for a key.
func (m *Map) Store(key, val uint32) {
	if key == isFree {
		if !m.hasFreeKey {
			m.count++
		}
		m.hasFreeKey = true
		m.freeVal = val
		return
	}

	data := m.data
	mask := m.mask[0]
	mask1 := m.mask[1]
	ptr := bucketOf(key, mask)
	dist := uint32(0)

	for {
		k := data[ptr]
		switch k {
		case isFree: // empty slot → place key here
			data[ptr] = key
			data[ptr+1] = val
			m.count++
			if m.count >= m.threshold {
				m.rehash()
			}
			return
		case key: // overwrite existing value
			data[ptr+1] = val
			return
		default:
			occDist := ((ptr - bucketOf(k, mask)) & mask1) >> 1
			if occDist < dist { // steal slot
				// swap (key,val,dist) ↔ occupant
				key, data[ptr] = data[ptr], key
				val, data[ptr+1] = data[ptr+1], val
				dist = occDist // continue insertion with displaced key
			}
			ptr = (ptr + 2) & mask1
			dist++
		}
	}
}

// Delete removes the value for a key.
func (m *Map) Delete(key uint32) {
	if key == isFree {
		if m.hasFreeKey {
			m.hasFreeKey = false
			m.count--
		}
		return
	}

	data := m.data
	mask := m.mask[0]
	mask1 := m.mask[1]
	ptr := bucketOf(key, mask)

	// find the key
	for {
		k := data[ptr]
		if k == key {
			break // found
		}
		if k == isFree {
			return // absent
		}
		ptr = (ptr + 2) & mask1
	}

	// back-shift deletion loop
	next := (ptr + 2) & mask1
	for {
		k := data[next]
		if k == isFree {
			data[ptr] = isFree
			m.count--
			return
		}
		home := bucketOf(k, mask)
		// distance the entry would have if we move it back one slot
		if ((next - home) & mask1) == 0 {
			data[ptr] = isFree
			m.count--
			return
		}
		// shift next back into ptr
		data[ptr] = k
		data[ptr+1] = data[next+1]
		ptr = next
		next = (next + 2) & mask1
	}
}

// Count returns number of key/value pairs in the map.
func (m *Map) Count() int { return int(m.count) }

// Range visits every key/value pair in the map.
func (m *Map) Range(fn func(key, val uint32) bool) {
	if m.hasFreeKey && !fn(isFree, m.freeVal) {
		return
	}
	for i := 0; i < len(m.data); i += 2 {
		if k := m.data[i]; k != isFree {
			if !fn(k, m.data[i+1]) {
				return
			}
		}
	}
}

// RangeEach visits every key/value pair without early‑exit capability.
func (m *Map) RangeEach(fn func(key, val uint32)) {
	if m.hasFreeKey {
		fn(isFree, m.freeVal)
	}
	for i := 0; i < len(m.data); i += 2 {
		if k := m.data[i]; k != isFree {
			fn(k, m.data[i+1])
		}
	}
}

// RangeErr stops on the first error returned by fn and propagates it.
func (m *Map) RangeErr(fn func(key, val uint32) error) error {
	if m.hasFreeKey {
		if err := fn(isFree, m.freeVal); err != nil {
			return err
		}
	}
	for i := 0; i < len(m.data); i += 2 {
		if k := m.data[i]; k != isFree {
			if err := fn(k, m.data[i+1]); err != nil {
				return err
			}
		}
	}
	return nil
}

// Clear removes all key/value pairs from the map.
func (m *Map) Clear() {
	clear(m.data)
	m.count = 0
	m.hasFreeKey = false
	m.freeVal = 0
}

// Clone returns a copy of the map.
func (m *Map) Clone() *Map {
	clone := New(len(m.data) / 2)
	clone.fillFactor = m.fillFactor
	clone.count = m.count
	clone.mask = m.mask
	clone.hasFreeKey = m.hasFreeKey
	clone.freeVal = m.freeVal
	copy(clone.data, m.data)
	return clone
}

// rehash doubles table size and reinserts all keys.
func (m *Map) rehash() {
	old := m.data
	newCap := len(old)
	if newCap >= math.MaxInt32/2 {
		panic("intmap: maximum size reached")
	}
	newCap *= 2

	m.data = make([]uint32, newCap)
	m.mask = [2]uint32{uint32(newCap/2 - 1), uint32(newCap - 1)}
	m.threshold = int32(float64(newCap/2) * float64(m.fillFactor))

	// reinsertion – Robin Hood store handles collisions.
	oldCount := m.count
	if m.hasFreeKey {
		m.count = 1
	} else {
		m.count = 0
	}
	for i := 0; i < len(old); i += 2 {
		if k := old[i]; k != isFree {
			m.Store(k, old[i+1])
		}
	}
	// after rehash Store increments m.count, so we assert equality
	_ = oldCount // (could sanity-check here in debug build)
}

// arraySize returns the smallest power-of-two ≥ size / fill.
// Panics if the result would overflow int.
/*func arraySize(size int, fill float64) int {
	if size <= 0 {
		panic("intmap: size must be positive")
	}
	if fill <= 0 || fill >= 1 {
		panic("intmap: fill factor must be in (0,1)")
	}

	need := uint64(math.Ceil(float64(size) / fill)) // exact ceiling
	if need < 8 {
		return 8
	}
	if need > uint64(^uint(0)) { // overflow check
		panic("intmap: requested capacity overflows int")
	}
	// next power-of-two: 1 << (bits.Len64(need-1))
	capacity := uint64(1) << (64 - bits.LeadingZeros64(need-1))
	return int(capacity)
}
*/

// arraySize returns the next power-of-two ≥ size/fill.
func arraySize(size int, fill float64) int {
	x := uint32(math.Ceil(float64(size) / fill))
	if x < 8 {
		return 8
	}
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return int(x + 1)
}
