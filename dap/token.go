package dap

import (
	"dap/utils"
)

//go:generate stringer -type=TokenType
type TokenType int

const (
	ttEOF TokenType = iota
	ttAdd
	ttAddAdd
	ttAssign
	ttBitwiseAnd
	ttBitwiseAndAssign
	ttBitwiseOr
	ttBitwiseOrAssign
	ttBlank
	ttComma
	ttConstNumber
	ttConstString
	ttDiv
	ttDot
	ttElse
	ttEqual
	ttGT
	ttGTE
	ttIf
	ttImport
	ttLeftBracket
	ttLeftCurve
	ttLeftParenthese
	ttLineEnd
	ttLogicAnd
	ttLogicAndAssign
	ttLogicOr
	ttLogicOrAssign
	ttLT
	ttLTE
	ttMulti
	ttNot
	ttNotEqual
	ttReturn
	ttRightCurve
	ttRightParenthese
	ttRigthBracket
	ttSemi
	ttShiftRight
	ttSub
	ttSymbol
	ttThis
	ttVar
)

func (tt TokenType) MarshalJSON() (bs []byte, err error) {
	bs = []byte(`"` + tt.String() + `"`)
	return
}

type Token struct {
	Typ  TokenType `json:"typ"`
	Line int       `json:"line"`
	Col  int       `json:"col"`
	Val  string    `json:"val"`
}

func (tt *Token) String() string {
	return utils.JsonStr(tt)
}
