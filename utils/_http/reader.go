package _http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Reader interface {
	Key() string
	Value() interface{}
	AsPayload() bool
}

type paramReader struct {
	asPayload bool
	key       string
	value     interface{}
}

func (p paramReader) Key() string {
	return p.key
}

func (p paramReader) Value() interface{} {
	return p.value
}

func (p paramReader) AsPayload() bool {
	return p.asPayload
}

type fileReader struct {
	key  string
	path string
}

func (p fileReader) Key() string {
	return p.key
}

func (p fileReader) Value() interface{} {
	return p.path
}

func (p fileReader) AsPayload() bool {
	return false
}

func Param(key string, value interface{}) Reader {
	return paramReader{
		key:       key,
		value:     value,
		asPayload: false,
	}
}

func Payload(v interface{}) Reader {
	return paramReader{
		value:     v,
		asPayload: true,
	}
}

func File(key string, path string) Reader {
	return fileReader{
		key:  key,
		path: path,
	}
}

func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func newReaders(list []Reader) (*readers, error) {
	var (
		err error
		r   = &readers{}
	)

	for _, v := range list {
		if v.AsPayload() {
			if r.payload {
				return nil, fmt.Errorf("mutilple reades as payload")
			} else {
				r.payload = true
				r.item = v
			}
		} else {
			r.list = append(r.list, v)
		}
	}

	return r, err
}

type readers struct {
	payload bool
	item    Reader
	list    []Reader
}

func (p readers) Json() io.Reader {

	if p.payload {
		d, err := json.Marshal(p.item.Value())
		if err != nil {
			panic(fmt.Errorf("encode json error: %s", err))
		}
		return bytes.NewReader(d)
	}

	var builder strings.Builder
	builder.WriteString("{")
	l := len(p.list)
	for i, v := range p.list {
		switch t := v.(type) {
		case paramReader:
			_param := v.(paramReader)
			builder.WriteString("\"")
			builder.WriteString(_param.key)
			builder.WriteString("\":")
			d, err := json.Marshal(_param.value)
			if err != nil {
				panic(fmt.Errorf("encode '%s' to json error: %s", _param.key, err))
			}
			builder.Write(d)
		default:
			panic(fmt.Errorf("unsupport encode '%s' type '%s' to json", v.Key(), t))
		}
		if i != l-1 {
			builder.WriteString(",")
		}

	}
	builder.WriteString("}")

	return strings.NewReader(builder.String())
}

func (p readers) Multipart(req *Request) io.Reader {

	if p.payload {
		switch t := p.item.Value().(type) {
		case []byte:
			return bytes.NewReader(p.item.Value().([]byte))
		default:
			panic(fmt.Errorf("unsupport %s cast to multipart data", t))
		}
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	defer func() {
		err := writer.Close()
		if err != nil {
			panic(fmt.Errorf("encode multipart data error: %s", err))
		}
	}()
	for _, v := range p.list {
		switch t := v.(type) {
		case paramReader:
			err := writer.WriteField(v.Key(), fmt.Sprintf("%v", v.Value()))
			if err != nil {
				panic(fmt.Errorf("write field error: %s", err))
			}
		case fileReader:
			err := p.WriteFile(writer, v.(fileReader))
			if err != nil {
				panic(fmt.Errorf("write file '%s' path '%v' error: %s", v.Key(), v.Value(), err))
			}
		default:
			panic(fmt.Errorf("unsupport %s cast to multipart field", t))
		}
	}
	req.SetContentType(writer.FormDataContentType())
	return payload
}

func (p readers) WriteFile(writer *multipart.Writer, reader fileReader) (err error) {

	file, err := os.Open(reader.path)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()

	part, err := writer.CreateFormFile(reader.key, filepath.Base(reader.path))
	_, err = io.Copy(part, file)
	if err != nil {
		return
	}

	return
}

func (p readers) UrlEncoded() io.Reader {

	if p.payload {
		panic(fmt.Errorf("unsupport cast as payload reader to url values"))
	}

	params := url.Values{}
	for _, v := range p.list {
		switch v.(type) {
		case paramReader:
			params.Add(v.(paramReader).key, fmt.Sprintf("%v", v.Value()))
		case fileReader:
			params.Add(v.(fileReader).key, fmt.Sprintf("%v", v.Value()))
		}
	}

	return strings.NewReader(params.Encode())
}

type Json map[string]interface{}

func (p Json) Key() string {
	return ""
}

func (p Json) AsPayload() bool {
	return true
}

func (p Json) Value() (interface{}, error) {
	return json.Marshal(p)
}
