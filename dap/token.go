package dap

const (
	ttEOF = iota
	ttBlank
	ttSymbol
	ttImport
	ttConstString
	ttLeftParenthese
	ttRightParenthese
	ttLeftBracket
	ttRigthBracket
	ttLeftCurve
	ttRightCurve
	ttLineEnd
	ttDot
	ttComma
	ttVar
	ttAssign
	ttReturn
	ttIf
	ttElse
	ttSemi
	ttAdd
	ttSub
	ttMulti
	ttDiv
)

type Token struct {
	typ  int
	line int
	col  int
	val  string
}
