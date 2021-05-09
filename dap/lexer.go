package dap

import (
	"dap/utils"
)

type Lexer interface {
	//peek() *Token
	get() *Token
	//fetch() *Token
	begin() Lexer
	done()
	commit()
	trans(f func(l Lexer) bool)
}

type RuneLexer struct {
	tok     *Token
	idx     int
	line    int
	col     int
	content []rune
}

func (l *RuneLexer) peek() *Token {
	if l.tok != nil {
		return l.tok
	}
	l.tok = l.fetch()
	return l.tok
}

func (l *RuneLexer) fetch() (res *Token) {
	runes := []rune{l.content[l.idx]}
	switch l.content[l.idx] {
	case ' ':
		for l.idx++; utils.RuneIn(l.content[l.idx], ' ', '\t', '\r', '\n'); l.idx++ {
			runes = append(runes, l.content[l.idx])
		}
		return &Token{
			line: 0,
			col:  0,
			val:  string(runes),
		}
	default:
		for l.idx++; !utils.RuneIn(l.content[l.idx], ' ', '\t', '\r', '\n'); l.idx++ {
			runes = append(runes, l.content[l.idx])
		}
		return &Token{
			line: 0,
			col:  0,
			val:  string(runes),
		}
	}
}

func (l *RuneLexer) get() (res *Token) {
	if l.tok != nil {
		res, l.tok = l.tok, nil
		return
	}
	return l.fetch()
}
