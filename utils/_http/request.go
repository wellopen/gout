package _http

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func Send(method Method, url string, params ...Reader) *Response {
	return (&Request{}).Send(method, url, params...)
}

func Post(url string, params ...Reader) *Response {
	return (&Request{}).Post(url, params...)
}

func Put(url string, params ...Reader) *Response {
	return (&Request{}).Put(url, params...)
}

func Delete(url string, params ...Reader) *Response {
	return (&Request{}).Delete(url, params...)
}

func Get(url string, params ...Reader) *Response {
	return (&Request{}).Get(url, params...)
}

//func Download(url string, data interface{} ){
//	return nil
//}
//
//func Upload(url string, data interface{} ){
//	return nil
//}

func NewRequest(method Method, url string, params ...Reader) *Request {
	return &Request{
		method: method,
		url:    url,
		params: params,
	}
}

type Request struct {
	name            string
	method          Method
	origin          *http.Request
	payload         io.Reader
	host            string
	header          http.Header
	timeout         time.Duration
	done            bool
	params          []Reader
	url             string
	failedCallback  func(response *Response)
	successCallback func(response *Response)
}

func (p *Request) SetTimeout(timeout time.Duration) *Request {
	p.timeout = timeout
	return p
}

func (p *Request) SetSuccessCallback(successCallback func(response *Response)) *Request {
	p.successCallback = successCallback
	return p
}

func (p *Request) SetFailedCallback(failedCallback func(response *Response)) *Request {
	p.failedCallback = failedCallback
	return p
}

func (p *Request) SetUrl(url string) *Request {
	p.url = url
	return p
}

func (p *Request) Url() string {
	return p.url
}

func (p *Request) SetParams(params []Reader) *Request {
	p.params = params
	return p
}

func (p *Request) Params() []Reader {
	return p.params
}

func (p *Request) SetMethod(method Method) *Request {
	p.method = method
	return p
}

func (p *Request) Method() Method {
	return p.method
}

func (p *Request) SetName(name string) *Request {
	p.name = name
	return p
}

func (p *Request) Name() string {
	return p.name
}

func (p *Request) Origin() *http.Request {
	return p.origin
}

func (p *Request) Header() http.Header {
	if p.header == nil {
		p.header = http.Header{}
	}
	return p.header
}

func (p *Request) Copy() *Request {
	r := new(Request)
	r.header = p.header
	r.timeout = p.timeout
	return p
}

func (p *Request) SetHeader(key, value string) *Request {
	p.Header().Set(key, value)
	return p
}

func (p *Request) SetContentType(contentType string) *Request {
	p.SetHeader(CONTENT_TYPE, contentType)
	return p
}

func (p *Request) Get(url string, params ...Reader) *Response {
	return p.Send(GET, url, params...)
}

func (p *Request) Post(url string, params ...Reader) *Response {
	return p.Send(POST, url, params...)
}

func (p *Request) Put(url string, params ...Reader) *Response {
	return p.Send(PUT, url, params...)
}

func (p *Request) Delete(url string, params ...Reader) *Response {
	return p.Send(DELETE, url, params...)
}

func (p *Request) Send(method Method, url string, params ...Reader) (response *Response) {
	
	response = NewResponse()
	
	defer func() {
		e := recover()
		if e != nil {
			response.setError(fmt.Errorf("%v", e))
		}
	}()
	
	urlBuilder := strings.Builder{}
	urlBuilder.WriteString(p.host)
	urlBuilder.WriteString(url)
	
	contentType := JSON
	t := p.Header().Get(CONTENT_TYPE)
	if t == "" {
		p.SetContentType(contentType)
	} else {
		contentType = t
	}
	
	_readers, err := newReaders(params)
	if err != nil {
		response.err = err
		return
	}
	
	var body io.Reader = nil
	if len(params) > 0 {
		switch method {
		case POST, PUT:
			body = p.prepareBody(contentType, _readers)
		case GET, DELETE:
			var _params []byte
			_params, err = ioutil.ReadAll(_readers.UrlEncoded())
			if err != nil {
				response.err = fmt.Errorf("reade url encoded params error: %s", err)
				return
			}
			urlBuilder.WriteString("?")
			urlBuilder.Write(_params)
		}
	}
	
	var ctx context.Context
	if p.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), p.timeout)
		defer cancel()
	} else {
		ctx = context.Background()
	}
	
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, string(method), urlBuilder.String(), body)
	if err != nil {
		response.setError(fmt.Errorf("new http request error: %s", err))
		return
	}
	
	now := time.Now()
	defer func() {
		response.cost = time.Since(now)
	}()
	
	req.Header = p.header
	response.request = p.Copy()
	response.request.payload = body
	response.request.origin = req
	
	resp, err := client.Do(req)
	if err != nil {
		response.setError(fmt.Errorf("send http reqeust error: %s", err))
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	
	response.origin = resp
	response.body, _ = ioutil.ReadAll(resp.Body)
	
	return
}

func (p *Request) prepareBody(contentType string, params *readers) io.Reader {
	switch contentType {
	case X_WWW_FORM_URLENCODED:
		return params.UrlEncoded()
	case JSON, JSON_UTF8:
		return params.Json()
	case FORM_DATA:
		return params.Multipart(p)
	}
	panic(fmt.Errorf("content-type '%s' unsupport prepare body", contentType))
}

func (p *Request) Exec() *Response {
	return p.Send(p.method, p.url, p.params...)
}
