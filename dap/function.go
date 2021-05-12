package dap

import "dap/utils"

type Function struct {
	Name  string
	Args  Args
	Exprs []Express
}

func (this *Function) String() string {
	return utils.JsonStr(this)
}

func (f *Function) Eval(s Scope) (res Object) {
	for _, expr := range f.Exprs {
		res = expr.Eval(s)
		s.pushTempObject(res)
	}
	return
}

type Arg struct {
	Name  string
	Class *ClassRef
}

type Args []*Arg

type Lambda struct {
	Args  Args
	Rets  Args
	Exprs []Express
}

func (f *Lambda) Eval(s Scope) (res Object) {
	for _, expr := range f.Exprs {
		res = expr.Eval(s)
		s.pushTempObject(res)
	}
	return
}
