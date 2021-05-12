package dap

type Function struct {
	name  string
	args  Args
	exprs []Express
}

func (f *Function) Eval(s Scope) (res Object) {
	for _, expr := range f.exprs {
		res = expr.Eval(s)
		s.pushTempObject(res)
	}
	return
}

type Arg struct {
	name  string
	class string
}

type Args []*Arg

type Lambda struct {
	args  Args
	ret   Args
	exprs []Express
}

func (f *Lambda) Eval(s Scope) (res Object) {
	for _, expr := range f.exprs {
		res = expr.Eval(s)
		s.pushTempObject(res)
	}
	return
}
