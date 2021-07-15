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

func (this *RuneLexer) begin() Lexer {
	ret := &RuneLexer{}
	*ret = *this
	ret.emitor = this
	return ret
}

func (this *RuneLexer) commit() {
	this.hasCommited = true
}

func (this *RuneLexer) done() {
	this.emitor.toks = this.toks
	this.emitor.cIdx = this.cIdx
	if this.hasCommited {
		this.emitor.tIdx = this.tIdx
		this.emitor.line = this.line
		this.emitor.col = this.col
	}
}

func (this *RuneLexer) get() (res *Token) {
	return this.getimpl()
}

func (this *RuneLexer) getimpl() (res *Token) {
	defer func() {
		this.tIdx++
	}()
	for len(this.toks) <= this.tIdx {
		tok := this.fetch()
		this.toks = append(this.toks, tok)
	}
	return this.toks[this.tIdx]
}

func (this *RuneLexer) peek() (res *Token) {
	return this.peekimpl()
}

func (this *RuneLexer) peekimpl() *Token {
	for len(this.toks) <= this.tIdx {
		tok := this.fetch()
		this.toks = append(this.toks, tok)
	}
	return this.toks[this.tIdx]
}

func (this *RuneLexer) _getc() rune {
	if this.cIdx >= len(this.Content) {
		return 0
	}
	r := this.Content[this.cIdx]
	this.cIdx++
	return r
}

func (this *RuneLexer) _peekc() rune {
	if this.cIdx >= len(this.Content) {
		return 0
	}
	return this.Content[this.cIdx]
}

func (this *RuneLexer) _begin() *RuneLexer {
	return &RuneLexer{
		emitor:  this,
		toks:    this.toks,
		tIdx:    this.tIdx,
		line:    this.line,
		col:     this.col,
		Content: this.Content,
		cIdx:    this.cIdx,
	}
}

func (this *RuneLexer) _done() {
	if this._hasCommited {
		this.emitor.cIdx = this.cIdx
		this.emitor.toks = this.toks
		this.emitor.line = this.line
		this.emitor.col = this.col
	}
}

func (this *RuneLexer) _commit() {
	this._hasCommited = true
}

func (this *RuneLexer) _trans(f func(tx *RuneLexer) bool) {
	tx := this._begin()
	defer tx._done()
	if f(tx) {
		tx._commit()
	}
}

func (this *RuneLexer) consume(str string, typ TokenType) (res *Token) {
	this._trans(func(tx *RuneLexer) bool {
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

func (this *RuneLexer) fetch() (res *Token) {
	if this.cIdx >= len(this.Content) {
		return &Token{
			Typ: ttEOF,
		}
	}
	r := this._peekc()
	if r == '\n' {
		this._getc()
		return &Token{
			Typ: ttLineEnd,
		}
	}
	if unicode.IsSpace(r) {
		runes := []rune{}
		for r := this._peekc(); unicode.IsSpace(r); r = this._peekc() {
			this._getc()
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
		this._getc()
		for r = this._peekc(); unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r); r = this._peekc() {
			this._getc()
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
		for r = this._peekc(); unicode.IsDigit(r); r = this._peekc() {
			this._getc()
			runes = append(runes, r)
		}
		return &Token{
			Typ: ttConstNumber,
			Val: string(runes),
		}
	}
	if res = this.consume(">=", ttGTE); res != nil {
		return
	}
	if res = this.consume(">>", ttShiftRight); res != nil {
		return
	}
	if res = this.consume(">", ttGT); res != nil {
		return
	}
	if res = this.consume("==", ttEqual); res != nil {
		return
	}
	if res = this.consume("++", ttAddAdd); res != nil {
		return
	}
	if res = this.consume("+", ttAdd); res != nil {
		return
	}
	if res = this.consume("-", ttSub); res != nil {
		return
	}
	if res = this.consume("*", ttMulti); res != nil {
		return
	}
	if res = this.consume("/", ttDiv); res != nil {
		return
	}
	if res = this.consume("(", ttLeftParenthese); res != nil {
		return
	}
	if res = this.consume(")", ttRightParenthese); res != nil {
		return
	}
	if res = this.consume("[", ttLeftBracket); res != nil {
		return
	}
	if res = this.consume("]", ttRigthBracket); res != nil {
		return
	}
	if res = this.consume("{", ttLeftCurve); res != nil {
		return
	}
	if res = this.consume("}", ttRightCurve); res != nil {
		return
	}
	if res = this.consume("==", ttEqual); res != nil {
		return
	}
	if res = this.consume("=", ttAssign); res != nil {
		return
	}
	if res = this.consume("!=", ttNotEqual); res != nil {
		return
	}
	if res = this.consume("!", ttNot); res != nil {
		return
	}
	if res = this.consume(".", ttDot); res != nil {
		return
	}
	if res = this.consume(",", ttComma); res != nil {
		return
	}
	if res = this.consume("&&=", ttLogicAndAssign); res != nil {
		return
	}
	if res = this.consume("&&", ttLogicAnd); res != nil {
		return
	}
	if res = this.consume("&=", ttBitwiseAndAssign); res != nil {
		return
	}
	if res = this.consume("&", ttBitwiseAnd); res != nil {
		return
	}
	if res = this.consume("||=", ttLogicOrAssign); res != nil {
		return
	}
	if res = this.consume("||", ttLogicOr); res != nil {
		return
	}
	if res = this.consume("|=", ttBitwiseOrAssign); res != nil {
		return
	}
	if res = this.consume("|", ttBitwiseOr); res != nil {
		return
	}
	log.Panicf("unknown char %s\n", string(r))
	return nil
}
