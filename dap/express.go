package dap

type Express interface {
	Eval(s Scope) Object
}

type ExprAdd struct {
	First  Express
	Second Express
}

func (this *ExprAdd) Eval(s Scope) Object {
	return this.First.Eval(s).Add(this.Second.Eval(s))
}

type ExprSub struct {
	First  Express
	Second Express
}

func (this *ExprSub) Eval(s Scope) Object {
	return this.First.Eval(s).Sub(this.Second.Eval(s))
}

type ExprMulti struct {
	First  Express
	Second Express
}

func (this *ExprMulti) Eval(s Scope) Object {
	return this.First.Eval(s).Multi(this.Second.Eval(s))
}

type ExprDiv struct {
	First  Express
	Second Express
}

func (this *ExprDiv) Eval(s Scope) Object {
	return this.First.Eval(s).Div(this.Second.Eval(s))
}

type ExprAssign struct {
	Target ExprAssignTarget
	Src    Express
}

type DefineOrRefOrReturn interface{}

type ExprAssignTarget []DefineOrRefOrReturn

func (this ExprAssignTarget) Assign(Object) Object {
	return nil
}

func (this *ExprAssign) Eval(s Scope) Object {
	return this.Target.Assign(this.Src.Eval(s))
}

type ExprIf struct {
	Prepare Express
	Cond    Express
	Exprs   []Express
	Else    Express
}

func (this *ExprIf) Eval(s Scope) Object {
	return nil
}

type ExprRef struct {
	Names []string
}

func (this *ExprRef) Eval(s Scope) Object {
	return nil
}

type ClassRef struct {
	Pkg  string
	Name string
}

type ExprDefine struct {
	Name  string
	Class *ClassRef
}

func (this *ExprDefine) Eval(s Scope) Object {
	return nil
}

type ExpressLiteralNumber struct {
	Val string
}

func (this *ExpressLiteralNumber) Eval(s Scope) Object {
	return nil
}
