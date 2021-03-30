package domui

type IsTesting bool

func (_ Def) IsTesting() (is IsTesting) {
	return isTesting
}

var isTesting = func() (is IsTesting) {
	if _, ok := ((any)(is)).(interface {
		Testing()
	}); ok {
		is = true
	}
	return
}()
