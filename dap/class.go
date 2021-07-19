package dap

import "dap/utils"

type Class struct {
	Name   string
	Fields map[string]*Field
	// Methods map[string]*Method
}

func (cls *Class) String() string {
	return utils.JsonStr(cls)
}

type Method struct {
	Class string
	Function
}

func (cls *Method) String() string {
	return utils.JsonStr(cls)
}

type Field struct {
	Name  string
	Class *ClassRef
}

func (cls *Field) String() string {
	return utils.JsonStr(cls)
}

func (m *Method) Eval(s Scope) Object {
	return m.Function.Eval(s)
}
