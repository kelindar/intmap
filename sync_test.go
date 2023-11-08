// Copyright (c) 2021-2023, Roman Atachiants

package intmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangeRandomSync(t *testing.T) {
	for _, size := range []int{100, 10000, 1000000} {
		count := 0
		m := sequentialSyncMap(size)
		m.Range(func(key, value uint32) bool {
			count++
			return true
		})
		assert.Equal(t, m.Count(), count)
	}
}

func TestLoadOrStoreLoaded(t *testing.T) {
	m := sequentialSyncMap(10)
	v, loaded := m.LoadOrStore(1, func() uint32 {
		return 1
	})
	assert.Equal(t, uint32(1), v)
	assert.True(t, loaded)
}

func TestLoadOrStoreMissed(t *testing.T) {
	m := sequentialSyncMap(10)
	v, loaded := m.LoadOrStore(20, func() uint32 {
		return 20
	})
	assert.Equal(t, uint32(20), v)
	assert.False(t, loaded)
}

func TestSyncDelete(t *testing.T) {
	m := sequentialSyncMap(10)
	m.Delete(1)

	_, ok := m.Load(1)
	assert.False(t, ok)
}
