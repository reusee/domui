package domui

import (
	"fmt"

	"github.com/reusee/dscope"
	"github.com/reusee/e4"
)

type (
	any = interface{}

	Scope = dscope.Scope
)

var (
	ce = e4.Check.With(e4.DropFrame(func(frame e4.Frame) bool {
		switch frame.Pkg {
		case "dscope", "runtime", "reflect":
			return true
		}
		return false
	}))
	he = e4.Handle
	sp = fmt.Sprintf

	Methods = dscope.Methods
)
