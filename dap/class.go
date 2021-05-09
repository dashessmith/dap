package dap

type Class struct {
	name    string
	fields  map[string]*Field
	methods map[string]*Method
}

type Method struct {
	This Object
	Function
}

type Field struct {
	name  string
	class string
}

func (m *Method) Eval(s Scope) Object {
	return m.Function.Eval(s)
}
