package koalanet

import "testing"

type TestActor struct {
	Actor
	fname string
	args  interface{}
	reply interface{}
}

func (svr *TestActor) setContext(ctx *context) {

}

func (svr *TestActor) setHandle(handle uint32) {

}

func (svr *TestActor) GetHandle() uint32 {
	return 0
}

func (svr *TestActor) OnMessage(funcName string, args interface{}, reply interface{}) error {
	svr.fname = funcName
	svr.args = args
	svr.reply = reply
	return nil
}

func Test_context_call(t *testing.T) {
	actor := &TestActor{}
	ctx := &context{}
	ctx.handle = 0
	ctx.name = "TestContext"
	ctx.actor = actor
	ctx.messageChan = make(chan *contextMessage, 2048)

	go context_thread(ctx)

	ctx.call(nil, "TestMethod", 32, nil)

	ctx.kill(false)

	if actor.fname != "TestMethod" {
		t.Errorf("call(TestMethod) failed.")
	}

	if actor.args.(int) != 32 {
		t.Errorf("call(32) failed.")
	}

	if actor.reply != nil {
		t.Fail()
	}

	contextWG.Wait()
}

func Benchmark_context(b *testing.B) {
	ctx := &context{}
	ctx.handle = 0
	ctx.name = "TestContext"
	ctx.actor = &TestActor{}
	ctx.messageChan = make(chan *contextMessage, 2048)

	go context_thread(ctx)

	for i := 0; i < b.N; i++ {
		ctx.send(nil, "TestMethod", 32)
	}

	ctx.kill(true)

	contextWG.Wait()
}
