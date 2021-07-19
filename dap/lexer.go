package dap

import (
	"log"
	"unicode"
)

type Lexer interface {
	peek() *Token
	get() *Token
	begin() Lexer
	done()
	commit()
}

type RuneLexer struct {
	emitor       *RuneLexer
	hasCommited  bool
	toks         []*Token
	tIdx         int
	line         int
	col          int
	Content      []rune
	cIdx         int
	_hasCommited bool
}

func (rl *RuneLexer) begin() Lexer {
	ret := &RuneLexer{}
	*ret = *rl
	ret.emitor = rl
	return ret
}

func (rl *RuneLexer) commit() {
	rl.hasCommited = true
}

func (rl *RuneLexer) done() {
	rl.emitor.toks = rl.toks
	rl.emitor.cIdx = rl.cIdx
	if rl.hasCommited {
		rl.emitor.tIdx = rl.tIdx
		rl.emitor.line = rl.line
		rl.emitor.col = rl.col
	}
}

func (rl *RuneLexer) get() (res *Token) {
	return rl.getimpl()
}

func (rl *RuneLexer) getimpl() (res *Token) {
	defer func() {
		rl.tIdx++
	}()
	for len(rl.toks) <= rl.tIdx {
		tok := rl.fetch()
		rl.toks = append(rl.toks, tok)
	}
	return rl.toks[rl.tIdx]
}

func (rl *RuneLexer) peek() (res *Token) {
	return rl.peekimpl()
}

func (rl *RuneLexer) peekimpl() *Token {
	for len(rl.toks) <= rl.tIdx {
		tok := rl.fetch()
		rl.toks = append(rl.toks, tok)
	}
	return rl.toks[rl.tIdx]
}

func (rl *RuneLexer) _getc() rune {
	if rl.cIdx >= len(rl.Content) {
		return 0
	}
	r := rl.Content[rl.cIdx]
	rl.cIdx++
	return r
}

func (rl *RuneLexer) _peekc() rune {
	if rl.cIdx >= len(rl.Content) {
		return 0
	}
	return rl.Content[rl.cIdx]
}

func (rl *RuneLexer) _begin() *RuneLexer {
	return &RuneLexer{
		emitor:  rl,
		toks:    rl.toks,
		tIdx:    rl.tIdx,
		line:    rl.line,
		col:     rl.col,
		Content: rl.Content,
		cIdx:    rl.cIdx,
	}
}

func (rl *RuneLexer) _done() {
	if rl._hasCommited {
		rl.emitor.cIdx = rl.cIdx
		rl.emitor.toks = rl.toks
		rl.emitor.line = rl.line
		rl.emitor.col = rl.col
	}
}

func (rl *RuneLexer) _commit() {
	rl._hasCommited = true
}

func (rl *RuneLexer) _trans(f func(tx *RuneLexer) bool) {
	tx := rl._begin()
	defer tx._done()
	if f(tx) {
		tx._commit()
	}
}

func (rl *RuneLexer) consume(str string, typ TokenType) (res *Token) {
	rl._trans(func(tx *RuneLexer) bool {
		var runes []rune
		for r := tx._getc(); r != 0; r = tx._getc() {
			runes = append(runes, r)
			if len(runes) >= len(str) {
				break
			}
		}
		if string(runes) != str {
			return false
		}
		res = &Token{
			Val: str,
			Typ: typ,
		}
		return true
	})
	return
}

func (rl *RuneLexer) fetch() (res *Token) {
	if rl.cIdx >= len(rl.Content) {
		return &Token{
			Typ: ttEOF,
		}
	}
	r := rl._peekc()
	if r == '\n' {
		rl._getc()
		return &Token{
			Typ: ttLineEnd,
		}
	}
	if unicode.IsSpace(r) {
		runes := []rune{}
		for r := rl._peekc(); unicode.IsSpace(r); r = rl._peekc() {
			rl._getc()
			runes = append(runes, r)
		}
		return &Token{
			Line: 0,
			Col:  0,
			Val:  string(runes),
			Typ:  ttBlank,
		}
	}
	if unicode.IsLetter(r) || r == '_' {
		runes := []rune{r}
		rl._getc()
		for r = rl._peekc(); unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r); r = rl._peekc() {
			rl._getc()
			runes = append(runes, r)
		}
		switch str := string(runes); str {
		case `import`:
			return &Token{Typ: ttImport}
		case `if`:
			return &Token{Typ: ttIf}
		case `else`:
			return &Token{Typ: ttElse}
		default:
			return &Token{Typ: ttSymbol, Val: str}
		}
	}
	if unicode.IsDigit(r) {
		runes := []rune{}
		for r = rl._peekc(); unicode.IsDigit(r); r = rl._peekc() {
			rl._getc()
			runes = append(runes, r)
		}
		return &Token{
			Typ: ttConstNumber,
			Val: string(runes),
		}
	}
	if res = rl.consume(">=", ttGTE); res != nil {
		return
	}
	if res = rl.consume(">>", ttShiftRight); res != nil {
		return
	}
	if res = rl.consume(">", ttGT); res != nil {
		return
	}
	if res = rl.consume("==", ttEqual); res != nil {
		return
	}
	if res = rl.consume("++", ttAddAdd); res != nil {
		return
	}
	if res = rl.consume("+", ttAdd); res != nil {
		return
	}
	if res = rl.consume("-", ttSub); res != nil {
		return
	}
	if res = rl.consume("*", ttMulti); res != nil {
		return
	}
	if res = rl.consume("/", ttDiv); res != nil {
		return
	}
	if res = rl.consume("(", ttLeftParenthese); res != nil {
		return
	}
	if res = rl.consume(")", ttRightParenthese); res != nil {
		return
	}
	if res = rl.consume("[", ttLeftBracket); res != nil {
		return
	}
	if res = rl.consume("]", ttRigthBracket); res != nil {
		return
	}
	if res = rl.consume("{", ttLeftCurve); res != nil {
		return
	}
	if res = rl.consume("}", ttRightCurve); res != nil {
		return
	}
	if res = rl.consume("==", ttEqual); res != nil {
		return
	}
	if res = rl.consume("=", ttAssign); res != nil {
		return
	}
	if res = rl.consume("!=", ttNotEqual); res != nil {
		return
	}
	if res = rl.consume("!", ttNot); res != nil {
		return
	}
	if res = rl.consume(".", ttDot); res != nil {
		return
	}
	if res = rl.consume(",", ttComma); res != nil {
		return
	}
	if res = rl.consume("&&=", ttLogicAndAssign); res != nil {
		return
	}
	if res = rl.consume("&&", ttLogicAnd); res != nil {
		return
	}
	if res = rl.consume("&=", ttBitwiseAndAssign); res != nil {
		return
	}
	if res = rl.consume("&", ttBitwiseAnd); res != nil {
		return
	}
	if res = rl.consume("||=", ttLogicOrAssign); res != nil {
		return
	}
	if res = rl.consume("||", ttLogicOr); res != nil {
		return
	}
	if res = rl.consume("|=", ttBitwiseOrAssign); res != nil {
		return
	}
	if res = rl.consume("|", ttBitwiseOr); res != nil {
		return
	}
	log.Panicf("unknown char %s\n", string(r))
	return nil
}
