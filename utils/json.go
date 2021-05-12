package utils

import "encoding/json"


func JsonStr(x interface{}) string {
	bs, _ := json.Marshal(x)
	return string(bs)
}
