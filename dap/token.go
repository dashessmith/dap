package dap

import (
	"dap/utils"
)

//go:generate stringer -type=TokenType
type TokenType int

const (
	ttEOF TokenType = iota
	ttBlank
	ttSymbol
	ttImport
	ttConstString
	ttConstNumber
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
	ttEqual
	ttGT
	ttShiftRight
	ttGTE
	ttNot
	ttNotEqual
	ttLT
	ttLTE
	ttReturn
	ttIf
	ttElse
	ttSemi
	ttAdd
	ttAddAdd
	ttSub
	ttMulti
	ttDiv
)

func (this TokenType) MarshalJSON() (bs []byte, err error) {
	bs = []byte(`"` + this.String() + `"`)
	return
}

type Token struct {
	Typ  TokenType `json:"typ"`
	Line int       `json:"line"`
	Col  int       `json:"col"`
	Val  string    `json:"val"`
}

func (this *Token) String() string {
	return utils.JsonStr(this)
}
