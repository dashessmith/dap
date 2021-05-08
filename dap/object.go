package dap

type Object interface {
	Add(Object) Object
	Sub(Object) Object
}


