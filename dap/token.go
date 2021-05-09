package dap

const (
	ttSymbol = iota + 1
	ttImport
	ttConstString
	ttLeftBrace
	ttRightBrace
	ttLineEnd
)

type Token struct {
	typ  int
	line int
	col  int
	val  string
}
