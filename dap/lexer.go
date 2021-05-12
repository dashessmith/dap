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
	content      []rune
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
	for res = this.getimpl(); res.Typ == ttBlank; res = this.getimpl() {
	}
	return
}

func (this *RuneLexer) getimpl() (res *Token) {
	if len(this.toks) > this.tIdx {
		res = this.toks[this.tIdx]
		this.tIdx++
		return
	}
	tok := this.fetch()
	this.toks = append(this.toks, tok)
	this.tIdx++
	return tok
}

func (this *RuneLexer) peek() (res *Token) {
	for res = this.peekimpl(); res.Typ == ttBlank; res = this.get() {
	}
	return
}
func (this *RuneLexer) peekimpl() *Token {
	if len(this.toks) > this.tIdx {
		return this.toks[this.tIdx]
	}
	tok := this.fetch()
	this.toks = append(this.toks, tok)
	return tok
}

func (this *RuneLexer) _getc() rune {
	if this.cIdx >= len(this.content) {
		return 0
	}
	r := this.content[this.cIdx]
	this.cIdx++
	return r
}

func (this *RuneLexer) _peekc() rune {
	if this.cIdx >= len(this.content) {
		return 0
	}
	return this.content[this.cIdx]
}

func (this *RuneLexer) _begin() *RuneLexer {
	return &RuneLexer{
		emitor:  this,
		toks:    this.toks,
		tIdx:    this.tIdx,
		line:    this.line,
		col:     this.col,
		content: this.content,
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
	if this.cIdx >= len(this.content) {
		return &Token{
			Typ: ttEOF,
		}
	}
	r := this._peekc()
	runes := []rune{r}
	if unicode.IsSpace(r) {
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
	if unicode.IsLetter(r) {
		for r = this._peekc(); unicode.IsLetter(r); r = this._peekc() {
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
		for r = this._peekc(); unicode.IsLetter(r); r = this._peekc() {
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
	log.Panicf("unknown char %s\n", string(r))
	return nil
}
