package domui

import (
	"fmt"
	"os"
	"reflect"
	"syscall/js"

	"github.com/reusee/sb"
	"github.com/reusee/store/schema"
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

type Reload func()

func (_ Def) Reload(
	scopedCall ScopedCall,
	call CallRemote,
) Reload {
	return func() {

		// save states
		scopedCall(func(
			scope Scope,
		) {
			states := make(map[string]any)
			ce(scope.RangePtrValues(func(t reflect.Type, vs []any) error {
				if len(vs) > 1 {
					// skip reducers
					return nil
				}
				name := sb.TypeName(t)
				if name == "" {
					return nil
				}
				pt("save state: %s\n", name)
				states[name] = vs[0]
				return nil
			}))
			call(schema.ShouldSaveStates{
				States: states,
			})
		})

		pt("about to reload, exit\n")
		Global.Set("reloading", true)
		os.Exit(0)
	}
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
