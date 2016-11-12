package knet

import (
	"errors"
	"sync"
)

type HandleStorage struct {
	lock        sync.RWMutex
	slotSize    uint32
	handleIndex uint32
	slot        []*context
	name        map[string]*context
}

var hs *HandleStorage

const (
	HANDLE_MASK       uint32 = 0xFFFFFF
	DEFAULT_SLOT_SIZE uint32 = 4
)

func register(ctx *context) uint32 {
	hs.lock.Lock()
	defer hs.lock.Unlock()

	var i, handle, hash uint32

	for {
		for i = 0; i < hs.slotSize; i++ {
			handle = (i + hs.handleIndex) & HANDLE_MASK
			hash = handle & (hs.slotSize - 1)
			if hs.slot[hash] == nil {
				hs.slot[hash] = ctx
				hs.handleIndex = handle + 1

				return handle
			}
		}

		new_slot := make([]*context, hs.slotSize*2)
		for i = 0; i < hs.slotSize; i++ {
			if hs.slot[i] != nil {
				hash = hs.slot[i].handle & (hs.slotSize*2 - 1)
				new_slot[hash] = hs.slot[i]
			}
		}
		hs.slot = new_slot
		hs.slotSize *= 2
	}
}

func retire(handle uint32) {
	hs.lock.Lock()
	defer hs.lock.Unlock()

	hash := handle & (hs.slotSize - 1)
	ctx := hs.slot[hash]

	if ctx != nil && ctx.handle == handle {
		hs.slot[hash] = nil
		if len(ctx.name) > 0 {
			delete(hs.name, ctx.name)
		}
	}
}

func retireAll() {
	hs.lock.Lock()
	defer hs.lock.Unlock()

	for i := uint32(0); i < hs.slotSize; i++ {
		ctx := hs.slot[i]

		if ctx != nil {
			hs.slot[i] = nil
			if len(ctx.name) > 0 {
				delete(hs.name, ctx.name)
			}
		}
	}
}

func get(handle uint32) *context {
	hs.lock.RLock()
	defer hs.lock.RUnlock()

	hash := handle & (hs.slotSize - 1)
	ctx := hs.slot[hash]

	if ctx != nil && ctx.handle == handle {
		return hs.slot[hash]
	}

	return nil
}

func getByName(name string) (uint32, error) {
	hs.lock.RLock()
	defer hs.lock.RUnlock()

	ctx := hs.name[name]
	if ctx == nil {
		return 0, errors.New("Can't find context.")
	}

	return ctx.handle, nil
}

func initHandleStorage() {
	hs = &HandleStorage{slotSize: DEFAULT_SLOT_SIZE, handleIndex: 1, slot: make([]*context, DEFAULT_SLOT_SIZE)}
	hs.name = make(map[string]*context)
}

func init() {
	initHandleStorage()
}
