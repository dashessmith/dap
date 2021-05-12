package dap

import (
	"log"
	"testing"
)

func Test_parser(t *testing.T) {
	parser := Parser{
		Lexer: &RuneLexer{
			content: []rune(`fdsa {} `),
		},
	}
	imports, classes, methods, err := parser.parse()
	t.Logf("\n%v\n%v\n%v\n%v\n", imports, classes, methods, err)
}

func Test_x(t *testing.T) {
	tt := ttImport
	log.Printf("%s\n", tt)
}
