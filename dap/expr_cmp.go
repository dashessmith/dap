package dap

type ExprCmp struct {
	OP     TokenType
	First  Express
	Second Express
}

func (exp *ExprCmp) Eval(s Scope) Object {
	return exp.First.Eval(s).GT(exp.Second.Eval(s))
}
