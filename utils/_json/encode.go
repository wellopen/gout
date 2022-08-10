package _json

import jsoniter "github.com/json-iterator/go"

func Marshal(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func MarshalMarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return jsoniter.MarshalIndent(v, prefix, indent)
}

func Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

func UnmarshalFromString(data string, v interface{}) error {
	return jsoniter.UnmarshalFromString(data, v)
}

func MarshalToString(v interface{}) (string, error) {
	return jsoniter.MarshalToString(v)
}
