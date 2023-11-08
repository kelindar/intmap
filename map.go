// Copyright (c) 2021-2023, Roman Atachiants
// Copyright (c) 2016, Brent Pedersen - Bioinformatics

package intmap

import (
	"math"
)

// isFree is the 'free' key
const isFree = 0

// Map is a map-like data-structure for int64s
type Map struct {
	data       []uint32  // Keys and values, interleaved keys
	fillFactor float32   // Desired fill factor
	threshold  int32     // Threshold for resize
	count      int32     // Number of elements in the map
	mask       [2]uint32 // Mask to calculate the original bucket and collisions
	freeVal    uint32    // Value of 'free' key
	hasFreeKey bool      // Whether 'free' key exists
}

// New returns a map initialized with n spaces and uses the stated fillFactor.
// The map will grow as needed.
func New(size int, fillFactor float64) *Map {
	if fillFactor <= 0 || fillFactor >= 1 {
		panic("intmap: fill factor must be in (0, 1)")
	}
	if size <= 0 {
		panic("intmap: size must be positive")
	}

	capacity := arraySize(size, fillFactor)
	return &Map{
		data:       make([]uint32, 2*capacity),
		fillFactor: float32(fillFactor),
		threshold:  int32(math.Floor(float64(capacity) * fillFactor)),
		mask:       [2]uint32{uint32(capacity - 1), uint32(2*capacity - 1)},
	}
}

// Load returns the value stored in the map for a key, or nil if no value is
// present. The ok result indicates whether value was found in the map.
func (m *Map) Load(key uint32) (uint32, bool) {
	if key == isFree {
		if m.hasFreeKey {
			return m.freeVal, true
		}
		return 0, false
	}

	ptr := bucketOf(key, m.mask[0])
	if ptr < 0 || ptr >= uint32(len(m.data)) { // Check to help to compiler to eliminate a bounds check below.
		return 0, false
	}

	switch m.data[ptr] {
	case isFree: // end of chain already
		return 0, false
	case key: // we check FREE prior to this call
		return m.data[ptr+1], true
	default:
		for {
			ptr = (ptr + 2) & m.mask[1]
			switch m.data[ptr] {
			case isFree:
				return 0, false
			case key:
				return m.data[ptr+1], true
			}
		}
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

	ptr := bucketOf(key, m.mask[0])
	switch m.data[ptr] {
	case isFree: // end of chain already
		m.data[ptr] = key
		m.data[ptr+1] = val
		if m.count >= m.threshold {
			m.rehash()
		} else {
			m.count++
		}
		return
	case key: // overwrite existed value
		m.data[ptr+1] = val
		return
	default:
		for {
			ptr = (ptr + 2) & m.mask[1]
			switch m.data[ptr] {
			case isFree:
				m.data[ptr] = key
				m.data[ptr+1] = val
				if m.count >= m.threshold {
					m.rehash()
				} else {
					m.count++
				}
				return
			case key:
				m.data[ptr+1] = val
				return
			}
		}
	}
}

// Delete deletes the value for a key.
func (m *Map) Delete(key uint32) {
	if m.hasFreeKey && key == isFree {
		m.hasFreeKey = false
		m.count--
		return
	}

	ptr := bucketOf(key, m.mask[0])
	switch m.data[ptr] {
	case isFree: // end of chain already
		return
	case key:
		m.shiftKeys(ptr)
		m.count--
		return
	default:
		for {
			ptr = (ptr + 2) & m.mask[1]
			switch m.data[ptr] {
			case isFree:
				return
			case key:
				m.shiftKeys(ptr)
				m.count--
				return
			}
		}
	}
}

// Count returns number of key/value pairs in the map.
func (m *Map) Count() int {
	return int(m.count)
}

// Range calls f sequentially for each key and value present in the map. If f
// returns false, range stops the iteration.
func (m *Map) Range(f func(key, value uint32) bool) {
	if m.hasFreeKey && !f(isFree, m.freeVal) {
		return
	}

	for i := 0; i < len(m.data); i += 2 {
		if k := m.data[i]; k != isFree {
			if !f(k, m.data[i+1]) {
				return
			}
		}
	}
}

// Clone returns a copy of the map.
func (m *Map) Clone() *Map {
	clone := New(len(m.data)/2, float64(m.fillFactor))
	clone.count = m.count
	clone.mask[0] = m.mask[0]
	clone.mask[1] = m.mask[1]
	clone.hasFreeKey = m.hasFreeKey
	clone.freeVal = m.freeVal
	copy(clone.data, m.data)
	return clone
}

// shiftKeys shifts entries with the same hash.
func (m *Map) shiftKeys(pos uint32) {
	var last, slot uint32
	var k uint32
	var data = m.data
	for {
		last = pos
		pos = (last + 2) & m.mask[1]
		for {
			k = data[pos]
			if k == isFree {
				data[last] = isFree
				return
			}

			slot = bucketOf(k, m.mask[0])
			if last <= pos {
				if last >= slot || slot > pos {
					break
				}
			} else {
				if last >= slot && slot > pos {
					break
				}
			}
			pos = (pos + 2) & m.mask[1]
		}
		data[last] = k
		data[last+1] = data[pos+1]
	}
}

// rehash rehashes the key space and resizes the map
func (m *Map) rehash() {
	newCapacity := len(m.data) * 2
	m.threshold = int32(math.Floor(float64(newCapacity/2) * float64(m.fillFactor)))
	m.mask = [2]uint32{uint32(newCapacity/2 - 1), uint32(newCapacity - 1)}

	// copy of original data
	data := make([]uint32, len(m.data))
	copy(data, m.data)

	m.data = make([]uint32, newCapacity)
	if m.hasFreeKey { // reset size
		m.count = 1
	} else {
		m.count = 0
	}

	var o uint32
	for i := 0; i < len(data); i += 2 {
		o = data[i]
		if o != isFree {
			m.Store(o, data[i+1])
		}
	}
}

// bucketOf calcultes the hash bucket for the integer key
func bucketOf(key, mask uint32) uint32 {
	h := key*0xdeece66d + 0xb
	return (h & mask) << 1
}

func capacityFor(x uint32) uint32 {
	if x == math.MaxUint32 {
		return x
	}

	if x == 0 {
		return 1
	}

	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return x + 1
}

func arraySize(size int, fill float64) int {
	s := capacityFor(uint32(math.Ceil(float64(size) / fill)))
	if s < 2 {
		s = 2
	}

	return int(s)
}
