package dap

type Express interface {
	Eval(s Scope) Object
}

type ExprAdd struct {
	first  Express
	second Express
}

func (this *ExprAdd) Eval(s Scope) Object {
	return this.first.Eval(s).Add(this.second.Eval(s))
}

type ExprSub struct {
	first  Express
	second Express
}

func (this *ExprSub) Eval(s Scope) Object {
	return this.first.Eval(s).Sub(this.second.Eval(s))
}

type ExprMulti struct {
	first  Express
	second Express
}

func (this *ExprMulti) Eval(s Scope) Object {
	return this.first.Eval(s).Multi(this.second.Eval(s))
}

type ExprDiv struct {
	first  Express
	second Express
}

func (this *ExprDiv) Eval(s Scope) Object {
	return this.first.Eval(s).Div(this.second.Eval(s))
}

type ExprAssign struct {
	target ExprAssignTarget
	src    Express
}

type DefineOrRefOrReturn interface{}

type ExprAssignTarget []DefineOrRefOrReturn

func (this ExprAssignTarget) Assign(Object) Object {
	return nil
}

func (this *ExprAssign) Eval(s Scope) Object {
	return this.target.Assign(this.src.Eval(s))
}

type ExprIf struct {
	prepare   Express
	condition Express
	exprs     []Express
	el        Express
}

func (this *ExprIf) Eval(s Scope) Object {
	return nil
}

type ExprRef struct {
	names []string
}

func (this *ExprRef) Eval(s Scope) Object {
	return nil
}

type ExprDefine struct {
	name  string
	class string
}

func (this *ExprDefine) Eval(s Scope) Object {
	return nil
}
