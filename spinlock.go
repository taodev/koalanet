package koalanet

import "sync/atomic"

//type SpinLock struct {
//	lock *sync.Mutex
//}

//func (locker *SpinLock) Lock() {
//	locker.lock.Lock()
//}

//func (locker *SpinLock) Unlock() {
//	locker.lock.Unlock()
//}

//func NewSpinLock() *SpinLock {
//	locker := &SpinLock{}
//	locker.lock = &sync.Mutex{}
//	return locker
//}

type SpinLock struct {
	lock int32
	// lock *sync.Mutex
}

func (locker *SpinLock) Lock() {
	for atomic.CompareAndSwapInt32(&locker.lock, 0, 1) {
	}
}

func (locker *SpinLock) Unlock() {
	atomic.StoreInt32(&locker.lock, 0)
	//	for atomic.CompareAndSwapInt32(&locker.lock, 1, 0) {
	//	}
}

func NewSpinLock() *SpinLock {
	return &SpinLock{0}
}
