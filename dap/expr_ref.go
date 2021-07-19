package dap

type ExprRef struct {
	Names []string
}

func (exp *ExprRef) Eval(s Scope) Object {
	return nil
}
