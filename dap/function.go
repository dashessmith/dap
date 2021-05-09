package dap

type Function struct {
	exprs []Express
}

func (f *Function) Eval(s Scope) (res Object) {
	for _, expr := range f.exprs {
		res = expr.Eval(s)
		s.pushTempObject(res)
	}
	return
}
