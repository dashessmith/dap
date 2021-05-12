package dap

import (
	"log"
	"testing"
)

func Test_parser(t *testing.T) {
	parser := Parser{
		Lexer: &RuneLexer{
			content: []rune(`
			fdsa {
			} 
			foo(){ 
				if 1 + 1 {
					1-2
				}
				(){}
			} 
			fdsa.__x123x(a xx.string, b){
			}
		`),
		},
	}
	imports, classes, methods, functions, err := parser.parse()
	t.Logf(`
	imports		%s
	classes		%s
	methods		%s
	functions	%s
	err 		%v
	`, imports, classes, methods, functions, err)
}

func Test_x(t *testing.T) {
	tt := ttImport
	log.Printf("%s\n", tt)
}
