package knet

import (
	"sync"
)

const (
	contextMessageCount uint32 = 2048
)

const (
	ctx_sysmsg_normal = 0
	ctx_sysmsg_quit   = 1
)

type contextMessage struct {
	src       *context
	fname     string
	args      interface{}
	reply     interface{}
	replyChan chan error
	op        int
}

var (
	contextQuitMsg *contextMessage = &contextMessage{nil, "", nil, nil, nil, ctx_sysmsg_quit}
)

var contextMessagePool *sync.Pool

func init() {
	contextMessagePool = &sync.Pool{
		New: func() interface{} {
			return &contextMessage{nil, "", nil, nil, nil, ctx_sysmsg_normal}
		},
	}
}

func contextMessageGet() *contextMessage {
	// return contextMessagePool.Get().(*contextMessage)
	return &contextMessage{nil, "", nil, nil, nil, ctx_sysmsg_normal}
}
