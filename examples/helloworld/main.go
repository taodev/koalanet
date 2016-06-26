//go:generate kactorgen actor

package main

import (
	"github.com/taodev/koalanet"
	_ "github.com/taodev/koalanet/examples/helloworld/actor"
)

func main() {
	koalanet.Run("HelloWorld", 0, false, "")
}
