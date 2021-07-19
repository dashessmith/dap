package dap

type ExprDiv struct {
	First  Express
	Second Express
}

func (exp *ExprDiv) Eval(s Scope) Object {
	return exp.First.Eval(s).Div(exp.Second.Eval(s))
}
