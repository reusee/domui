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
	ce, he = e4.Check, e4.Handle
	sp     = fmt.Sprintf
)
