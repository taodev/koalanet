package examples

import "github.com/taodev/koalanet"

type TestActorWrap struct {
	Handle uint32
}

func (actor *TestActorWrap) Init(isSync bool, args TestActor_InitArgs) error {
	if isSync {
		return koalanet.Call(actor.Handle, "Init", args, nil)
	}
	
	return koalanet.Send(actor.Handle, "Init", args)
}

func (actor *TestActorWrap) Call1(isSync bool, args TestActor_Call1Args, reply *TestActor_Call1Reply) error {
	if isSync {
		return koalanet.Call(actor.Handle, "Call1", args, reply)
	}
	
	return koalanet.Send(actor.Handle, "Call1", args)
}

type TestActorBWrap struct {
	Handle uint32
}

func (actor *TestActorBWrap) Init(isSync bool) error {
	if isSync {
		return koalanet.Call(actor.Handle, "Init", nil, nil)
	}
	
	return koalanet.Send(actor.Handle, "Init", nil)
}


type TestActorImpl struct {
	TestActor
}

func (actor *TestActorImpl) InitWrap(args interface{}, reply interface{}) error {
	return actor.Init(args.(TestActor_InitArgs))
}

func (actor *TestActorImpl) Call1Wrap(args interface{}, reply interface{}) error {
	return actor.Call1(args.(TestActor_Call1Args), reply.(*TestActor_Call1Reply))
}

type TestActorBImpl struct {
	TestActorB
}

func (actor *TestActorBImpl) InitWrap(args interface{}, reply interface{}) error {
	return actor.Init()
}

func init() {
	koalanet.RegActor("TestActor", func() koalanet.IActor {
		actor := &TestActorImpl{}
		actor.InitActor()
		actor.RegMethod("Init", actor.InitWrap)
		actor.RegMethod("Call1", actor.Call1Wrap)
		return actor
	})

	koalanet.RegActor("TestActorB", func() koalanet.IActor {
		actor := &TestActorBImpl{}
		actor.InitActor()
		actor.RegMethod("Init", actor.InitWrap)
		return actor
	})
}