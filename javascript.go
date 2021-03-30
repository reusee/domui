package domui

import (
	"fmt"
	"syscall/js"
)

var (
	Global      = js.Global()
	IPCRenderer = Global.Get("ipcRenderer")
	Console     = Global.Get("console")
	Document    = Global.Get("document")
	HTMLElement = Global.Get("HTMLElement")
	Body        = Document.Get("body")
)

func init() {
	Global.Set("frontendOK", true)
}

func log(format string, args ...any) {
	Console.Call("log", fmt.Sprintf(format, args...))
}

func logErr(format string, args ...any) {
	Console.Call("error", fmt.Sprintf(format, args...))
}

func warn(format string, args ...any) {
	Console.Call("warn", fmt.Sprintf(format, args...))
}

func pr(args ...any) {
	Console.Call("log", args...)
}

func jsInvoke(method string, args []any, cb func(ret js.Value)) {
	var fn js.Func
	fn = js.FuncOf(func(_this js.Value, args []js.Value) any {
		cb(args[0])
		fn.Release()
		return nil
	})
	Global.Call("invoke", method, args, fn)
}
