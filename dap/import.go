package dap

import "dap/utils"

type Import struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (i *Import) String() string {
	return utils.JsonStr(i)
}
