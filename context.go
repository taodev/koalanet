package koalanet

import (
	"fmt"
	"log"
	"sync"

	"time"
)

const (
	contextDefaultTimeout time.Duration = time.Second * 3
)

var (
	contextWG sync.WaitGroup
)

type IActor interface {
	setContext(ctx *context)
	setHandle(h uint32)
	getHandle() uint32
	OnMessage(funcName string, args interface{}, reply interface{}) error
}

type stackInfo struct {
}

type context struct {
	handle      uint32
	name        string
	actor       IActor
	messageChan chan *contextMessage
	timeout     int64
	stackInfo   *stackInfo
}

func (ctx *context) sendMessage(message *contextMessage) {
	ctx.messageChan <- message
}

func (ctx *context) send(src *context, fname string, args interface{}) error {
	// msg := contextMessageGet()
	msg := &contextMessage{src, fname, args, nil, nil, ctx_sysmsg_normal}

	ctx.messageChan <- msg

	//	select {
	//	case ctx.messageChan <- msg:
	//	case <-time.After(contextDefaultTimeout):
	//		return fmt.Errorf("time out")
	//	}

	return nil
}

func (ctx *context) setTimeout(timeout int64) {
	ctx.timeout = timeout
}

func (ctx *context) call(src *context, fname string, args interface{}, reply interface{}) error {
	msg := contextMessageGet()
	msg.src = src
	msg.fname = fname
	msg.args = args
	msg.reply = reply
	msg.replyChan = make(chan error, 1)

	// 发送到消息队列
	ctx.messageChan <- msg

	if ctx.timeout != 0 {
		select {
		case ret, _ := <-msg.replyChan:
			{
				return ret
			}
		case <-time.After(time.Second * time.Duration(ctx.timeout)):
			{
				// return fmt.Errorf("time out context current handle function:%s", ctx.Handle_func)
				return fmt.Errorf("time out context current handle function")
			}
		}
	} else {
		return <-msg.replyChan
	}

	return nil
}

func (ctx *context) kill(force bool) error {
	if force {
		log.Printf("context force killed, message(%d).", len(ctx.messageChan))
		close(ctx.messageChan)
		return nil
	}

	ctx.messageChan <- contextQuitMsg
	return nil
}

func context_thread(ctx *context) {
	countMsg := 0
	contextWG.Add(1)

	defer func() {
		//		if err := recover(); err != nil {
		//			log.Println(err)
		//		}

		//		log.Printf("context_thread[%d]: destroy(%d)", ctx.handle, countMsg)
		contextWG.Done()
		// close(ctx.messageChan)
	}()

	for {
		msg, ok := <-ctx.messageChan
		if !ok {
			break
		}

		countMsg++

		if msg.op == ctx_sysmsg_quit {
			break
		}

		err := ctx.actor.OnMessage(msg.fname, msg.args, msg.reply)
		if msg.replyChan != nil {
			msg.replyChan <- err
		}
	}
}
