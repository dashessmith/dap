package dap

import "dap/utils"

type Class struct {
	Name   string
	Fields map[string]*Field
	// Methods map[string]*Method
}

func (this *Class) String() string {
	return utils.JsonStr(this)
}

type Method struct {
	Class string
	Function
}

func (this *Method) String() string {
	return utils.JsonStr(this)
}

type Field struct {
	Name  string
	Class *ClassRef
}

func (this *Field) String() string {
	return utils.JsonStr(this)
}

func (m *Method) Eval(s Scope) Object {
	return m.Function.Eval(s)
}
