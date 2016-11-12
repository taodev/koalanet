package examples

import "github.com/taodev/koalanet"

type TestActor struct {
	koalanet.Actor
}

type TestActor_InitArgs struct {
}

func testActorFunc() {

}

func (t *TestActor) Init(args TestActor_InitArgs) error {
	return nil
}

func (t *TestActor) testA(args TestActor_InitArgs) {
}

type TestActor_Call1Args struct {
}

type TestActor_Call1Reply struct {
}

func (t *TestActor) Call1(args TestActor_Call1Args, reply *TestActor_Call1Reply) error {
	return nil
}

type TestActorB struct {
	koalanet.Actor
}

func (t *TestActorB) Init() error {
	return nil
}
