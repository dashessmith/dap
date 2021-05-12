package dap

import (
	"dap/utils"
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
	if len(this.toks) > this.tIdx {
		return this.toks[this.tIdx]
	}
	tok := this.fetch()
	this.toks = append(this.toks, tok)
	this.tIdx++
	return tok
}

func (this *RuneLexer) peek() *Token {
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

func (this *RuneLexer) _done() *RuneLexer {
	if this._hasCommited {
		this.emitor.cIdx = this.cIdx
		this.emitor.toks = this.toks
		this.emitor.line = this.line
		this.emitor.col = this.col
	}
	return this.emitor
}

func (this *RuneLexer) _commit() {
	this._hasCommited = true
}

func (this *RuneLexer) _trans(f func(l *RuneLexer) bool) {
	l := this._begin()
	defer l._done()
	if f(l) {
		l._commit()
	}
}

func (this *RuneLexer) fetch() (res *Token) {
	if this.cIdx >= len(this.content) {
		return &Token{
			typ: ttEOF,
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
			line: 0,
			col:  0,
			val:  string(runes),
			typ:  ttBlank,
		}
	}

	switch this.content[this.cIdx] {
	case ' ':

		for this.cIdx++; utils.RuneIn(this.content[this.cIdx], ' ', '\t', '\r', '\n'); this.cIdx++ {
			runes = append(runes, this.content[this.cIdx])
		}

	default:
		for this.cIdx++; !utils.RuneIn(this.content[this.cIdx], ' ', '\t', '\r', '\n'); this.cIdx++ {
			runes = append(runes, this.content[this.cIdx])
		}
		return &Token{
			line: 0,
			col:  0,
			val:  string(runes),
		}
	}
}
