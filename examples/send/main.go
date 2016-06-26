//go:generate kactorgen

package main

import (
	"log"

	"github.com/taodev/koalanet"
)

type MainActor struct {
	koalanet.Actor
}

func (act *MainActor) Init() error {
	log.Printf("MainActor::Init")

	wrap := MainActorWrap{act.GetHandle()}
	wrap.MethodA(false, ArgsSend{200, "Hello"})

	koalanet.KillActor(act.GetHandle(), false)
	koalanet.Exit()

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
