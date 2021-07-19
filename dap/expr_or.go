package dap

type ExprOr struct {
	First  Express
	Second Express
}

func (exp *ExprOr) Eval(s Scope) Object {
	return nil
}
