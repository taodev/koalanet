package knet

import (
	"fmt"
	"log"
)

type IActorMethod interface {
	Process() error
}

type ActorMethodCallback func(args interface{}, reply interface{}) error

type Actor struct {
	ctx       *context
	handle    uint32
	callbacks map[string]ActorMethodCallback
}

func (a *Actor) setContext(ctx *context) {
	a.ctx = ctx
}

func (a *Actor) setHandle(h uint32) {
	a.handle = h
}

func (a *Actor) GetHandle() uint32 {
	return a.handle
}

func (a *Actor) GetType() string {
	return "Actor"
}

func (a *Actor) InitActor() {
	a.callbacks = make(map[string]ActorMethodCallback)
}

func (a *Actor) OnMessage(funcName string, args interface{}, reply interface{}) error {
	callback, ok := a.callbacks[funcName]
	if !ok {
		return fmt.Errorf("context:%d not function(%s).", funcName)
	}

	return callback(args, reply)
}

func (a *Actor) RegMethod(name string, callback ActorMethodCallback) {
	_, ok := a.callbacks[name]
	if ok {
		return
	}

	a.callbacks[name] = callback
}

type ActorNewFunc func() IActor

type ActorFactory struct {
	actors map[string]ActorNewFunc
}

var (
	actorFactory *ActorFactory = &ActorFactory{make(map[string]ActorNewFunc)}
)

func RegActor(name string, fptr ActorNewFunc) {
	_, ok := actorFactory.actors[name]
	if ok {
		// panic(fmt.Sprintf("RegActor %s 已存在.", name))
		// log.Printf("Warning RegActor %s 已存在.", name)
	}

	actorFactory.actors[name] = fptr
}

func NewActor(actorName string, args interface{}) uint32 {
	factor, ok := actorFactory.actors[actorName]
	if !ok {
		log.Printf("can't find.")
		return 0
	}

	actor := factor()
	ctx := &context{
		actor:       actor,
		messageChan: make(chan *contextMessage, 2048),
	}
	ctx.actorType = actor.GetType()
	ctx.init()

	handle := register(ctx)
	ctx.handle = handle
	actor.setContext(ctx)
	actor.setHandle(handle)

	go context_thread(ctx)

	ctx.send(nil, "Init", args)

	return handle
}

func KillActor(h uint32, force bool) {
	ctx := get(h)
	if ctx == nil {
		log.Printf("Can't find Context. %d", h)
		return
	}

	ctx.kill(force)
}

func KillActorByName(name string, force bool) {
	handle, err := getByName(name)
	if err == nil {
		log.Printf("Can't find Context. %d", err.Error())
		return
	}

	KillActor(handle, force)
}

func init() {

}
