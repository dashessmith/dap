package dap

type Express interface {
	Eval(s Scope) Object
}

type ClassRef struct {
	Pkg  string
	Name string
}
