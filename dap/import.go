package dap

import "dap/utils"

type Import struct {
	name string
	path string
}

func (this *Import) String() string {
	return utils.JsonStr(this)
}
