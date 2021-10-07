package pyatomic

//提供atomic-bool的实现

import "sync/atomic"

type AtomicBool uint32

func NewAtomicBool() *AtomicBool {
	return new(AtomicBool)
}

func NewAtomicBoolWithVal(val bool) *AtomicBool {
	b := NewAtomicBool()
	b.Set(val)
	return b
}

func (b *AtomicBool) CompareAndSwap(oldv, val bool) bool {
	var oldvu, newvu uint32
	if oldv {
		oldvu = 1
	}
	if val {
		newvu = 1
	}
	return atomic.CompareAndSwapUint32((*uint32)(b), oldvu, newvu)
}

func (b *AtomicBool) Set(v bool) {
	if v {
		atomic.StoreUint32((*uint32)(b), 1) //true
	} else {
		atomic.StoreUint32((*uint32)(b), 0) //false
	}
}

func (b *AtomicBool) IsTrue() bool {
	return atomic.LoadUint32((*uint32)(b)) == 1
}
