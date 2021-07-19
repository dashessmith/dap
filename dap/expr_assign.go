package dap

type DefineOrRefOrReturn interface{}

type ExprAssignTarget []DefineOrRefOrReturn

func (exp ExprAssignTarget) Assign(Object) Object {
	return nil
}

type ExprAssign struct {
	Target ExprAssignTarget
	Src    Express
}

func (exp *ExprAssign) Eval(s Scope) Object {
	return exp.Target.Assign(exp.Src.Eval(s))
}
