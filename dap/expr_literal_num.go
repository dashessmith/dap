package dap

type ExpressLiteralNumber struct {
	Val string
}

func (exp *ExpressLiteralNumber) Eval(s Scope) Object {
	return nil
}
