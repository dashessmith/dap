package dap

import "testing"

func Test_parser(t *testing.T) {
	parser := Parser{
		Lexer: RuneLexer{},
	}
	parser.parse()
}
