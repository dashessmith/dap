package dap

import (
	"log"
	"testing"
	//"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
)

func TestMain(m *testing.M) {
	log.Printf(`test main
	`)
	// profiler.Start(profiler.Config{
	// 	ApplicationName: "simple.golang.app123",

	// 	// replace this with the address of pyroscope server
	// 	ServerAddress: "http://localhost:4040",

	// 	// by default all profilers are enabled,
	// 	//   but you can select the ones you want to use:
	// 	ProfileTypes: []profiler.ProfileType{
	// 		profiler.ProfileCPU,
	// 		profiler.ProfileAllocObjects,
	// 		profiler.ProfileAllocSpace,
	// 		profiler.ProfileInuseObjects,
	// 		profiler.ProfileInuseSpace,
	// 	},
	// })
	m.Run()
}

func Test_parser(t *testing.T) {

	parser := Parser{
		Lexer: &RuneLexer{
			Content: []rune(`
			student {
				birth xx.string
				score list
				enter_date
			 } 
			 student.new(){ 
				if this.birth > 100 {
					
				}
			 } 
			main(){
				
			}
		`),
		},
	}
	imports, classes, methods, functions, err := parser.Parse()
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
