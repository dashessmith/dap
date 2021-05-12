package dap

import (
	"log"
	"testing"
)

func Test_parser(t *testing.T) {
	parser := Parser{
		Lexer: &RuneLexer{
			content: []rune(`
			student {
				birth
				score
				enter_date
			 } 
			 student.new(){ 

			 } 
			main(){
				
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
