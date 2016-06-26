package main

import "github.com/taodev/koalanet"

type ActorAWrap struct {
	Handle uint32
}

func (actor *ActorAWrap) Init(isSync bool) error {
	if isSync {
		return koalanet.Call(actor.Handle, "Init", nil, nil)
	}
	
	return koalanet.Send(actor.Handle, "Init", nil)
}

func (actor *ActorAWrap) MethodA(isSync bool, reply *int) error {
	if isSync {
		return koalanet.Call(actor.Handle, "MethodA", nil, reply)
	}
	
	return koalanet.Send(actor.Handle, "MethodA", nil)
}

func (actor *ActorAWrap) MethodB(isSync bool) error {
	if isSync {
		return koalanet.Call(actor.Handle, "MethodB", nil, nil)
	}
	
	return koalanet.Send(actor.Handle, "MethodB", nil)
}

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


type ActorAImpl struct {
	ActorA
}

func (actor *ActorAImpl) InitWrap(args interface{}, reply interface{}) error {
	return actor.Init()
}

func (actor *ActorAImpl) MethodAWrap(args interface{}, reply interface{}) error {
	return actor.MethodA(reply.(*int))
}

func (actor *ActorAImpl) MethodBWrap(args interface{}, reply interface{}) error {
	return actor.MethodB()
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
	koalanet.RegActor("ActorA", func() koalanet.IActor {
		actor := &ActorAImpl{}
		actor.InitActor()
		actor.RegMethod("Init", actor.InitWrap)
		actor.RegMethod("MethodA", actor.MethodAWrap)
		actor.RegMethod("MethodB", actor.MethodBWrap)
		return actor
	})

	koalanet.RegActor("MainActor", func() koalanet.IActor {
		actor := &MainActorImpl{}
		actor.InitActor()
		actor.RegMethod("Init", actor.InitWrap)
		actor.RegMethod("MethodA", actor.MethodAWrap)
		return actor
	})
}