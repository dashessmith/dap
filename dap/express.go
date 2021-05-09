package dap

type Express interface {
	Eval(s Scope) Object
}

type ExprAdd struct {
	first  Object
	second Object
}

func (eadd *ExprAdd) Eval(s Scope) Object {
	return eadd.first.Add(eadd.second)
}

type ExprSub struct {
	first  Object
	second Object
}

func (esub *ExprSub) Eval(s Scope) Object {
	return esub.first.Sub(esub.second)
}
