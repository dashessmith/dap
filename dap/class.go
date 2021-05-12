package dap

import "dap/utils"

type Class struct {
	name    string
	fields  map[string]*Field
	methods map[string]*Method
}

func (this *Class) String() string {
	return utils.JsonStr(this)
}

type Method struct {
	class string
	Function
}

func (this *Method) String() string {
	return utils.JsonStr(this)
}

type Field struct {
	name  string
	class string
}

func (this *Field) String() string {
	return utils.JsonStr(this)
}

func (m *Method) Eval(s Scope) Object {
	return m.Function.Eval(s)
}
