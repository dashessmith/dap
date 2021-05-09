package dap

type Package struct {
	classes   map[string]Class
	functions map[string]Function
	variables map[string]Object
}
