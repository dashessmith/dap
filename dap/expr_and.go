package dap

type ExprAnd struct {
	First  Express
	Second Express
}

func (exp *ExprAnd) Eval(s Scope) Object {
	return nil
}
