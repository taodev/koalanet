package koalanet

import (
	"fmt"
	"log"
	"sync"

	"github.com/taodev/koalanet/toolbox"

	"time"
)

const (
	contextDefaultTimeout time.Duration = time.Second * 3
)

var (
	contextWG sync.WaitGroup
)

type stackInfo struct {
}

type context struct {
	handle      uint32
	name        string
	actor       IActor
	actorType   string
	messageChan chan *contextMessage
	timeout     int64
	stackInfo   *stackInfo
	wg          sync.WaitGroup
}

func (ctx *context) init() {
}

func (ctx *context) send(src *context, fname string, args interface{}) error {
	msg := &contextMessage{src, fname, args, nil, nil, ctx_sysmsg_normal}

	ctx.messageChan <- msg

	return nil
}

func (ctx *context) setTimeout(timeout int64) {
	ctx.timeout = timeout
}

func (ctx *context) call(src *context, fname string, args interface{}, reply interface{}) error {
	msg := &contextMessage{src, fname, args, reply, make(chan error, 1), ctx_sysmsg_normal}

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
				close(msg.replyChan)
				msg.replyChan = nil
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
	ctx.wg.Add(1)

	maxMQCount := 0

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}

		log.Printf("context_thread[%d]: destroy(%d:%d)", ctx.handle, maxMQCount, countMsg)
		contextWG.Done()
		ctx.wg.Done()
	}()

	for {
		msg, ok := <-ctx.messageChan
		if !ok {
			break
		}

		startTime := time.Now()

		countMsg++
		if maxMQCount < len(ctx.messageChan) {
			maxMQCount = len(ctx.messageChan) + 1
		}

		if msg.op == ctx_sysmsg_quit {
			close(ctx.messageChan)
			break
		}

		err := ctx.actor.OnMessage(msg.fname, msg.args, msg.reply)
		if msg.replyChan != nil {
			msg.replyChan <- err
		}

		timeDur := time.Since(startTime)
		toolbox.StatisticsMap.AddStatistics(ctx.actorType, msg.fname, timeDur)
	}
}
