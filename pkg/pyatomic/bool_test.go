package pyatomic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomicBool(t *testing.T) {
	val := NewAtomicBoolWithValue(true)
	assert.True(t, val.IsTrue())
	val.Set(false)
	assert.False(t, val.IsTrue())
	val.Set(true)
	assert.True(t, val.IsTrue())
	val.Set(false)
	assert.False(t, val.IsTrue())
	ok := val.CompareAndSwap(false, true)
	assert.True(t, ok)
	assert.True(t, val.IsTrue())
	ok = val.CompareAndSwap(true, false)
	assert.True(t, ok)
	assert.False(t, val.IsTrue())
	//
	ok = val.CompareAndSwap(true, false)
	assert.False(t, ok)
	assert.False(t, val.IsTrue())
}
