package actor

import (
	"log"

	"github.com/taodev/koalanet"
)

type HelloWorld struct {
	koalanet.Actor
}

func (h *HelloWorld) Init() error {
	log.Printf("Hello World!")
	koalanet.Exit()
	return nil
}
