package utils

import "encoding/json"

func JsonStr(x interface{}) string {
	bs, _ := json.MarshalIndent(x, ``, `    `)
	return string(bs)
}
