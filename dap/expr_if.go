package dap

type ExprIf struct {
	Prepare Express
	Cond    Express
	Exprs   []Express
	Else    Express
}

func (exp *ExprIf) Eval(s Scope) Object {
	return nil
}
