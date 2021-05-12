package dap

import "dap/utils"

type Import struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (this *Import) String() string {
	return utils.JsonStr(this)
}
