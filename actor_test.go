package koalanet

import "log"

type TestActor1 struct {
	Actor
}

func (a1 *TestActor1) Init(args interface{}) {

}

func (a *TestActor1) OnMessage(funcName string, args interface{}, reply interface{}) error {
	log.Printf("TestActor1::OnMessage:%s", funcName)
	return nil
}

//func Test_actor_send(t *testing.T) {
//	RegActor("TestActor1", func() IActor { return &TestActor1{} })

//	hA1 := NewActor("TestActor1", nil)

//	if len(hs.slot) == 0 {
//		t.Errorf("Test_actor_send NewActor failed.")
//	}

//	ctx := get(hA1)
//	ctx.send(nil, "TestFunc", nil)
//	ctx.kill(false)

//	// contextWG.Wait()
//	ctx.wg.Wait()
//}
