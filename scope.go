package domui

type Def struct{}

type Update func(decls ...any) Scope
