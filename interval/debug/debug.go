package debug

import (
	"fmt"
	"runtime/debug"

	"github.com/zooyer/gvm/interval/conf"
)

func Println(v ...interface{}) {
	if !conf.Debug {
		return
	}

	info, ok := debug.ReadBuildInfo()
	if ok {
		fmt.Println(info)
	}
	fmt.Println(v...)
}
