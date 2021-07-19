package dap

type ExprDefine struct {
	Name  string
	Class *ClassRef
}

func (exp *ExprDefine) Eval(s Scope) Object {
	return nil
}
