package _http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitly/go-simplejson"
)

type Status struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"params"`
}

func StatusParser() Parser {
	return func(data []byte) (result interface{}, err error) {
		result = new(Status)
		err = json.Unmarshal(data, result)
		if err != nil {
			return
		}
		return
	}
}

func StatusDataParser(pointer interface{}) Parser {
	return func(data []byte) (result interface{}, err error) {
		j, err := simplejson.NewJson(data)
		if err != nil {
			err = fmt.Errorf("simplejson.NewJson error: %s", err)
			return
		}
		code, err := j.Get("code").Int()
		if err != nil {
			err = fmt.Errorf("get code error: %s", err)
			return
		}
		message, err := j.Get("types").String()
		if err != nil {
			err = fmt.Errorf("get types error: %s", err)
			return
		}
		if code != http.StatusOK {
			if message != "" {
				message = fmt.Sprintf("(%s)", message)
			}
			err = fmt.Errorf("status error: %d%s", code, message)
			return
		}

		_data, err := j.Get("params").Encode()
		if err != nil {
			return
		}
		err = json.Unmarshal(_data, pointer)
		if err != nil {
			err = fmt.Errorf("json.Unmarshal _data error: %s", err)
			return
		}
		result = pointer
		return
	}
}
