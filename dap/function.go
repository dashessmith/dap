package dap

type Function struct {
	exprs []Express
}

func (f *Function) Eval() (res Object) {
	for _, expr := range f.exprs {
		res = expr.Eval()
	}
	return
}
