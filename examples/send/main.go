//go:generate kactorgen

package main

import (
	"log"

	"github.com/taodev/koalanet"
)

type ActorA struct {
	koalanet.Actor
}

func (act *ActorA) Init() error {
	log.Printf("ActorA::Init")
	koalanet.RegName(act.GetHandle(), "act1")

	return nil
}

func (act *ActorA) MethodA(reply *int) error {
	*reply = 32
	return nil
}

func (act *ActorA) MethodB() error {
	log.Printf("Actor::MethodB")
	return nil
}

type MainActor struct {
	koalanet.Actor
}

func (act *MainActor) Init() error {
	log.Printf("MainActor::Init")

	hAct1 := koalanet.NewActor("ActorA", nil)

	act1 := &ActorAWrap{}
	act1.Handle = hAct1

	// sync call
	reply := 0
	act1.MethodA(true, &reply)
	log.Printf("Call ActorA::MethodA %d", reply)

	// asyn send
	act1.MethodB(false)

	// koalanet.KillActor(hAct1, false)

	// koalanet.KillActor(act.GetHandle(), false)

	// koalanet.Exit()

	return nil
}

type ArgsSend struct {
	M1 int
	M2 string
}

func (act *MainActor) MethodA(args ArgsSend) error {
	log.Println(args)
	return nil
}

func main() {
	koalanet.Run("MainActor", 0, false, "")
	koalanet.WaitAllQuit()
}
