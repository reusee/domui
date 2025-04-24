package domui

import (
	"fmt"

	"github.com/reusee/dscope"
	"github.com/reusee/e5"
)

type (
	any = interface{}

	Scope = dscope.Scope
)

var (
	ce = e5.Check.With(e5.DropFrame(func(frame e5.Frame) bool {
		switch frame.Pkg {
		case "dscope", "runtime", "reflect":
			return true
		}
		return false
	}))
	he = e5.Handle
	sp = fmt.Sprintf

	Methods = dscope.Methods
)
