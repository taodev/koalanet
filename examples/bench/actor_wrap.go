package main

import "github.com/taodev/koalanet"

type WriterWrap struct {
	Handle uint32
}

func (actor *WriterWrap) Init(isSync bool) error {
	if isSync {
		return koalanet.Call(actor.Handle, "Init", nil, nil)
	}
	
	return koalanet.Send(actor.Handle, "Init", nil)
}

func (actor *WriterWrap) StartCall(isSync bool) error {
	if isSync {
		return koalanet.Call(actor.Handle, "StartCall", nil, nil)
	}
	
	return koalanet.Send(actor.Handle, "StartCall", nil)
}

func (actor *WriterWrap) StartSend(isSync bool) error {
	if isSync {
		return koalanet.Call(actor.Handle, "StartSend", nil, nil)
	}
	
	return koalanet.Send(actor.Handle, "StartSend", nil)
}

type ReaderWrap struct {
	Handle uint32
}

func (actor *ReaderWrap) Init(isSync bool) error {
	if isSync {
		return koalanet.Call(actor.Handle, "Init", nil, nil)
	}
	
	return koalanet.Send(actor.Handle, "Init", nil)
}

func (actor *ReaderWrap) MethodCall(isSync bool) error {
	if isSync {
		return koalanet.Call(actor.Handle, "MethodCall", nil, nil)
	}
	
	return koalanet.Send(actor.Handle, "MethodCall", nil)
}

func (actor *ReaderWrap) MethodSend(isSync bool) error {
	if isSync {
		return koalanet.Call(actor.Handle, "MethodSend", nil, nil)
	}
	
	return koalanet.Send(actor.Handle, "MethodSend", nil)
}

type BenchActorWrap struct {
	Handle uint32
}

func (actor *BenchActorWrap) Init(isSync bool) error {
	if isSync {
		return koalanet.Call(actor.Handle, "Init", nil, nil)
	}
	
	return koalanet.Send(actor.Handle, "Init", nil)
}


type WriterImpl struct {
	Writer
}

func (actor *WriterImpl) GetType() string {
	return "Writer"
}

func (actor *WriterImpl) InitWrap(args interface{}, reply interface{}) error {
	return actor.Init()
}

func (actor *WriterImpl) StartCallWrap(args interface{}, reply interface{}) error {
	return actor.StartCall()
}

func (actor *WriterImpl) StartSendWrap(args interface{}, reply interface{}) error {
	return actor.StartSend()
}

type ReaderImpl struct {
	Reader
}

func (actor *ReaderImpl) GetType() string {
	return "Reader"
}

func (actor *ReaderImpl) InitWrap(args interface{}, reply interface{}) error {
	return actor.Init()
}

func (actor *ReaderImpl) MethodCallWrap(args interface{}, reply interface{}) error {
	return actor.MethodCall()
}

func (actor *ReaderImpl) MethodSendWrap(args interface{}, reply interface{}) error {
	return actor.MethodSend()
}

type BenchActorImpl struct {
	BenchActor
}

func (actor *BenchActorImpl) GetType() string {
	return "BenchActor"
}

func (actor *BenchActorImpl) InitWrap(args interface{}, reply interface{}) error {
	return actor.Init()
}

func init() {
	koalanet.RegActor("Writer", func() koalanet.IActor {
		actor := &WriterImpl{}
		actor.InitActor()
		actor.RegMethod("Init", actor.InitWrap)
		actor.RegMethod("StartCall", actor.StartCallWrap)
		actor.RegMethod("StartSend", actor.StartSendWrap)
		return actor
	})

	koalanet.RegActor("Reader", func() koalanet.IActor {
		actor := &ReaderImpl{}
		actor.InitActor()
		actor.RegMethod("Init", actor.InitWrap)
		actor.RegMethod("MethodCall", actor.MethodCallWrap)
		actor.RegMethod("MethodSend", actor.MethodSendWrap)
		return actor
	})

	koalanet.RegActor("BenchActor", func() koalanet.IActor {
		actor := &BenchActorImpl{}
		actor.InitActor()
		actor.RegMethod("Init", actor.InitWrap)
		return actor
	})
}