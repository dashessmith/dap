package dap

type Class struct{}

type Method struct {
	This Object
	Function
}

func (m *Method) Eval() Object {
	return m.Function.Eval()
}
