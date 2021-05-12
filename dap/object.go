package dap

type Object interface {
	Add(Object) Object
	Sub(Object) Object
	Multi(Object) Object
	Div(Object) Object
	Assign(Object) Object
	GT(Object) Object
}
