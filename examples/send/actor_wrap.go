package main

import "github.com/taodev/koalanet"

type MainActorWrap struct {
	Handle uint32
}

func (actor *MainActorWrap) Init(isSync bool) error {
	if isSync {
		return koalanet.Call(actor.Handle, "Init", nil, nil)
	}
	
	return koalanet.Send(actor.Handle, "Init", nil)
}

func (actor *MainActorWrap) MethodA(isSync bool, args ArgsSend) error {
	if isSync {
		return koalanet.Call(actor.Handle, "MethodA", args, nil)
	}
	
	return koalanet.Send(actor.Handle, "MethodA", args)
}


type MainActorImpl struct {
	MainActor
}

func (actor *MainActorImpl) InitWrap(args interface{}, reply interface{}) error {
	return actor.Init()
}

func (actor *MainActorImpl) MethodAWrap(args interface{}, reply interface{}) error {
	return actor.MethodA(args.(ArgsSend))
}

func init() {
	koalanet.RegActor("MainActor", func() koalanet.IActor {
		actor := &MainActorImpl{}
		actor.InitActor()
		actor.RegMethod("Init", actor.InitWrap)
		actor.RegMethod("MethodA", actor.MethodAWrap)
		return actor
	})
}