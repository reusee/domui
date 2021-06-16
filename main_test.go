package domui

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Div      = Tag("div")
	P        = Tag("p")
	FontSize = Style("font-size")
	OnClick  = On("click")

	pt = fmt.Printf
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
