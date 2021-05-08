package dap

type Express interface {
	Eval() Object
}

type ExprAdd struct {
	first  Object
	second Object
}

func (eadd *ExprAdd) Eval() Object {
	return eadd.first.Add(eadd.second)
}

type ExprSub struct {
	first  Object
	second Object
}

func (esub *ExprSub) Eval() Object {
	return esub.first.Sub(esub.second)
}
