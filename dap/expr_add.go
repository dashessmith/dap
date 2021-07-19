package dap

type ExprAdd struct {
	First  Express
	Second Express
}

func (exp *ExprAdd) Eval(s Scope) Object {
	return exp.First.Eval(s).Add(exp.Second.Eval(s))
}
