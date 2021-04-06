package domui

import (
	"fmt"
	"syscall/js"
)

var (
	global      = js.Global()
	console     = global.Get("console")
	document    = global.Get("document")
	htmlElement = global.Get("HTMLElement")
	body        = document.Get("body")
)

func log(format string, args ...any) {
	console.Call("log", fmt.Sprintf(format, args...))
}

func logErr(format string, args ...any) {
	console.Call("error", fmt.Sprintf(format, args...))
}

func warn(format string, args ...any) {
	console.Call("warn", fmt.Sprintf(format, args...))
}

func pr(args ...any) {
	console.Call("log", args...)
}
