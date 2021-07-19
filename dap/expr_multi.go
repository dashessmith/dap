package dap

type ExprMulti struct {
	First  Express
	Second Express
}

func (exp *ExprMulti) Eval(s Scope) Object {
	return exp.First.Eval(s).Multi(exp.Second.Eval(s))
}
