// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package dap

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ttEOF-0]
	_ = x[ttAdd-1]
	_ = x[ttAddAdd-2]
	_ = x[ttAssign-3]
	_ = x[ttBitwiseAnd-4]
	_ = x[ttBitwiseAndAssign-5]
	_ = x[ttBitwiseOr-6]
	_ = x[ttBitwiseOrAssign-7]
	_ = x[ttBlank-8]
	_ = x[ttComma-9]
	_ = x[ttConstNumber-10]
	_ = x[ttConstString-11]
	_ = x[ttDiv-12]
	_ = x[ttDot-13]
	_ = x[ttElse-14]
	_ = x[ttEqual-15]
	_ = x[ttGT-16]
	_ = x[ttGTE-17]
	_ = x[ttIf-18]
	_ = x[ttImport-19]
	_ = x[ttLeftBracket-20]
	_ = x[ttLeftCurve-21]
	_ = x[ttLeftParenthese-22]
	_ = x[ttLineEnd-23]
	_ = x[ttLogicAnd-24]
	_ = x[ttLogicAndAssign-25]
	_ = x[ttLogicOr-26]
	_ = x[ttLogicOrAssign-27]
	_ = x[ttLT-28]
	_ = x[ttLTE-29]
	_ = x[ttMulti-30]
	_ = x[ttNot-31]
	_ = x[ttNotEqual-32]
	_ = x[ttReturn-33]
	_ = x[ttRightCurve-34]
	_ = x[ttRightParenthese-35]
	_ = x[ttRigthBracket-36]
	_ = x[ttSemi-37]
	_ = x[ttShiftRight-38]
	_ = x[ttSub-39]
	_ = x[ttSymbol-40]
	_ = x[ttThis-41]
	_ = x[ttVar-42]
}

const _TokenType_name = "ttEOFttAddttAddAddttAssignttBitwiseAndttBitwiseAndAssignttBitwiseOrttBitwiseOrAssignttBlankttCommattConstNumberttConstStringttDivttDotttElsettEqualttGTttGTEttIfttImportttLeftBracketttLeftCurvettLeftParenthesettLineEndttLogicAndttLogicAndAssignttLogicOrttLogicOrAssignttLTttLTEttMultittNotttNotEqualttReturnttRightCurvettRightParenthesettRigthBracketttSemittShiftRightttSubttSymbolttThisttVar"

var _TokenType_index = [...]uint16{0, 5, 10, 18, 26, 38, 56, 67, 84, 91, 98, 111, 124, 129, 134, 140, 147, 151, 156, 160, 168, 181, 192, 208, 217, 227, 243, 252, 267, 271, 276, 283, 288, 298, 306, 318, 335, 349, 355, 367, 372, 380, 386, 391}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
