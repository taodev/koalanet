package knet

import "testing"

var t1 *testing.T

var tkm_arg1 int

type TestKoalanetMain struct {
	Actor
}

func (tkm *TestKoalanetMain) Init(args interface{}, reply interface{}) error {
	tkm_arg1 = 2
	// *reply.(*int) = 64
	return nil
}

func (tkm *TestKoalanetMain) SendTest(args interface{}, reply interface{}) error {
	// log.Printf("TestKoalanetMain::SendTest")

	// *reply.(*int) = 32
	return nil
}

//func Test_koalanet(t *testing.T) {
//	t1 = t
//	RegActor("TestKoalanetMain", func() IActor {
//		actor := &TestKoalanetMain{}
//		actor.init()
//		actor.RegMethod("Init", actor.Init)
//		actor.RegMethod("SendTest", actor.SendTest)
//		return actor
//	})

//	hTKM := NewActor("TestKoalanetMain", nil)

//	//	var a1 int = 0
//	//	Call(hTKM, "SendTest", 120, &a1)
//	//	if a1 != 32 {
//	//		t.Errorf("a1 != 32")
//	//	}

//	KillActor(hTKM, false)

//	WaitActorQuit(hTKM)

//	//	if tkm_arg1 != 2 {
//	//		t.Errorf("tkm_arg1 != 2")
//	//	}
//}

func Benchmark_koalanet_send(b *testing.B) {
	RegActor("TestKoalanetMain", func() IActor {
		actor := &TestKoalanetMain{}
		actor.InitActor()
		actor.RegMethod("Init", actor.Init)
		actor.RegMethod("SendTest", actor.SendTest)
		return actor
	})

	hTKM := NewActor("TestKoalanetMain", nil)

	for i := 0; i < b.N; i++ {
		Send(hTKM, "SendTest", nil)
	}

	KillActor(hTKM, false)

	WaitActorQuit(hTKM)
}
