package domui

import (
	"reflect"
	"sync"

	"github.com/reusee/sb"
	"github.com/reusee/store/schema"
)

var nameToType sync.Map

type GetState func(obj any)

func (_ Def) StateFuncs(
	call CallRemote,
) (
	get GetState,
) {

	get = func(obj any) {
		t := reflect.TypeOf(obj)
		name := sb.TypeName(t)
		nameToType.LoadOrStore(name, t)
		call(schema.ShouldSendState{
			Type: name,
		})
	}

	return
}

func (_ Def) StateRPCs(
	update Update,
) RPCs {
	return RPCs{

		func(set schema.SendState) {
			v, ok := nameToType.Load(set.Type)
			if !ok {
				return
			}
			obj := reflect.New(v.(reflect.Type))
			ce(sb.Copy(
				set.Tokens.Iter(),
				sb.Unmarshal(obj.Interface()),
			))
			update(obj.Interface())
			pt("load state: %s\n", set.Type)
		},

		//
	}
}
