package _json

import "github.com/valyala/fastjson"

type Value struct {
	*fastjson.Value
}

func Parse(v []byte) (parser *Value, err error) {
	a, err := fastjson.ParseBytes(v)
	if err != nil {
		return
	}
	parser = &Value{Value: a}
	return
}

func ParseString(v string) (parser *Value, err error) {
	a, err := fastjson.Parse(v)
	if err != nil {
		return
	}
	parser = &Value{Value: a}
	return
}
