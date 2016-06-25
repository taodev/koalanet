package koalanet

import (
	"sync/atomic"
)

type SpinLock struct {
	lock int32
}

func (locker *SpinLock) Lock() {
	for atomic.CompareAndSwapInt32(&locker.lock, 0, 1) {
	}
}

func (locker *SpinLock) Unlock() {
	atomic.StoreInt32(&locker.lock, 0)
}

func NewSpinLock() *SpinLock {
	return &SpinLock{0}
}
