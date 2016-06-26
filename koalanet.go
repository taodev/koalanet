package koalanet

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"
)

var (
	quit_chan chan int = make(chan int, 1)
)

type IActor interface {
	setContext(ctx *context)
	setHandle(h uint32)
	GetHandle() uint32
	OnMessage(funcName string, args interface{}, reply interface{}) error
}

func Send(handle uint32, fname string, args interface{}) error {
	ctx := get(handle)
	if ctx == nil {
		return fmt.Errorf("context:%d is nil.", handle)
	}

	return ctx.send(nil, fname, args)
}

func SendByName(name string, fname string, args interface{}) error {
	handle, err := getByName(name)
	if err != nil {
		return err
	}

	return Send(handle, fname, args)
}

func Call(handle uint32, fname string, args interface{}, reply interface{}) error {
	ctx := get(handle)
	if ctx == nil {
		return fmt.Errorf("context:%d is nil.", handle)
	}

	return ctx.call(nil, fname, args, reply)
}

func CallByName(name string, fname string, args interface{}, reply interface{}) error {
	handle, err := getByName(name)
	if err != nil {
		return err
	}

	return Call(handle, fname, args, reply)
}

func QueryName(name string) (uint32, error) {
	return getByName(name)
}

func RegName(handle uint32, name string) error {
	ctx := get(handle)
	if ctx == nil {
		return fmt.Errorf("Can't find Context. %s", name)
	}

	hs.lock.Lock()
	hs.name[name] = ctx
	ctx.name = name
	hs.lock.Unlock()

	return nil
}

func Run(mainActor string, maxThread int, debugEnable bool, debugPort string) {
	if maxThread <= 0 {
		maxThread = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(maxThread)

	if debugEnable {
		go func() {
			if debugPort != "" {
				log.Println(http.ListenAndServe(debugPort, nil))
				return
			}
			log.Println(http.ListenAndServe(":6067", nil))
		}()
	}

	log.Println("Welcome gonet.")

	NewActor(mainActor, "")

	for {
		select {
		case <-time.After(1 * time.Minute):
			debug.FreeOSMemory()
			runtime.GC()
		case <-quit_chan:
			return
		}
	}
}

func Exit() {
	quit_chan <- 1
}

func WaitAllQuit() {
	contextWG.Wait()
}

func WaitActorQuit(handle uint32) error {
	ctx := get(handle)
	if ctx == nil {
		return fmt.Errorf("Can't find Context. %d", handle)
	}

	ctx.wg.Wait()

	return nil
}

func init() {
	log.Println("KoalaNet init.")
}
