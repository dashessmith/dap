package dap

type ExprSub struct {
	First  Express
	Second Express
}

func (exp *ExprSub) Eval(s Scope) Object {
	return exp.First.Eval(s).Sub(exp.Second.Eval(s))
}
