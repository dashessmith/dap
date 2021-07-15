package dap

import (
	"io/ioutil"
	"strings"
)

func Run(pat string) (err error) {
	files, err := ioutil.ReadDir(pat)
	if err != nil {
		return
	}
	for _, f := range files {
		if f.Mode().IsRegular() {
			if strings.HasSuffix(f.Name(), ".dap") {
			}
		}
	}
	return
}
