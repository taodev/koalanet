package koalanet

import (
	"fmt"
	"log"
	"os"
)

type Actor struct {
	ctx    *context
	handle uint32
}

func (a *Actor) setContext(ctx *context) {
	a.ctx = ctx
}

func (a *Actor) setHandle(h uint32) {
	a.handle = h
}

func (a *Actor) getHandle() uint32 {
	return a.handle
}

func (a *Actor) init(ctx *context, handle uint32) {
	a.ctx = ctx
	a.handle = handle
}

func (a *Actor) OnMessage(funcName string, args interface{}, reply interface{}) error {
	return nil
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
		panic(fmt.Sprintf("RegActor %s 已存在.", name))
	}

	actorFactory.actors[name] = fptr
}

func NewActor(actorName string, args interface{}) uint32 {
	defer func() {
		if err := recover(); err != nil {
			errorInfo := fmt.Sprint(err)
			log.Printf("NewActor: panic:%s", errorInfo)
			os.Exit(1)
			return
		}
	}()

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
