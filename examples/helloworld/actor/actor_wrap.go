package actor

import "github.com/taodev/koalanet"

type HelloWorldWrap struct {
	Handle uint32
}

func (actor *HelloWorldWrap) Init(isSync bool) error {
	if isSync {
		return koalanet.Call(actor.Handle, "Init", nil, nil)
	}
	
	return koalanet.Send(actor.Handle, "Init", nil)
}


type HelloWorldImpl struct {
	HelloWorld
}

func (actor *HelloWorldImpl) InitWrap(args interface{}, reply interface{}) error {
	return actor.Init()
}

func init() {
	koalanet.RegActor("HelloWorld", func() koalanet.IActor {
		actor := &HelloWorldImpl{}
		actor.InitActor()
		actor.RegMethod("Init", actor.InitWrap)
		return actor
	})
}