package _json

import (
	"encoding/json"
	"fmt"
)

func Debug(val interface{}) {
	d, _ := json.MarshalIndent(val, "", "\t")
	fmt.Println(string(d))
}
