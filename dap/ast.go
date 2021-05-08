package dap

type Ast struct {
	main Function
}

func (a *Ast) Eval() Object {
	return a.main.Eval()
}
