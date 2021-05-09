package dap

type Ast struct {
	tmpObject Object
	pkgs      map[string]*Package
	main      Function
}

func (a *Ast) pushTempObject(obj Object) {
	a.tmpObject = obj
}

func (a *Ast) Eval() Object {
	return a.main.Eval(a)
}
