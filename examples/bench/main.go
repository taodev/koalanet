//go:generate kactorgen

package main

import (
	"log"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/taodev/koalanet"
)

var (
	callTime     uint64 = 0
	lastCallTime uint64 = 0

	sendTime     uint64 = 0
	lastSendTime uint64 = 0
)

type Writer struct {
	koalanet.Actor
}

func (w *Writer) Init() error {
	return nil
}

func (w *Writer) StartCall() error {
	hReader := koalanet.NewActor("Reader", nil)
	reader := &ReaderWrap{hReader}

	for {
		reader.MethodCall(true)
	}
	return nil
}

func (w *Writer) StartSend() error {
	hReader := koalanet.NewActor("Reader", nil)
	reader := &ReaderWrap{hReader}

	for {
		reader.MethodSend(false)
	}
	return nil
}

type Reader struct {
	koalanet.Actor
}

func (r *Reader) Init() error {
	return nil
}

func (r *Reader) MethodCall() error {
	atomic.AddUint64(&callTime, 1)
	return nil
}

func (r *Reader) MethodSend() error {
	atomic.AddUint64(&sendTime, 1)
	return nil
}

type BenchActor struct {
	koalanet.Actor
}

func (b *BenchActor) Init() error {
	//	for i := 0; i < 1000; i++ {
	//		hWriter := koalanet.NewActor("Writer", nil)
	//		koalanet.Send(hWriter, "StartCall", nil)
	//	}

	for i := 0; i < 1000; i++ {
		hWriter := koalanet.NewActor("Writer", nil)
		koalanet.Send(hWriter, "StartSend", nil)
	}

	return nil
}

func main() {
	go func() {
		lastTime := time.Now()

		for {
			<-time.After(time.Second * 1)
			atomic.LoadUint64(&callTime)
			atomic.LoadUint64(&sendTime)

			now := time.Now()
			t1 := now.UnixNano() - lastTime.UnixNano()
			lastTime = now

			t2 := float64(t1) / float64(time.Second.Nanoseconds())

			callCount := float64(callTime - lastCallTime)
			sendCount := float64(sendTime - lastSendTime)

			log.Printf("Time %d:%d", uint64(callCount/t2), uint64(sendCount/t2))

			lastCallTime = callTime
			lastSendTime = sendTime
		}
	}()

	log.Println("CPU: ", runtime.NumCPU())
	koalanet.Run("BenchActor", runtime.NumCPU(), false, "")
}
