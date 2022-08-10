package _http

import (
	"fmt"
	"net/http"
	"github.com/wellopen/send/utils/_json"
	"sync"
	"time"
)

func NewResponse() *Response {
	return &Response{}
}

type Response struct {
	request *Request
	err     error
	header  http.Header
	cost    time.Duration
	body    []byte
	origin  *http.Response
}

func (p *Response) Header() http.Header {
	if p.header == nil {
		p.header = http.Header{}
	}
	return p.header
}

func (p *Response) Cost() time.Duration {
	return p.cost
}

func (p *Response) Code() int {
	if p.origin == nil {
		return 0
	}
	return p.origin.StatusCode
}

func (p *Response) Status() string {
	if p.origin == nil {
		return ""
	}
	return p.origin.Status
}

func (p *Response) Request() *Request {
	return p.request
}

func (p *Response) setError(err error) {
	if p.err == nil {
		p.err = err
	}
}

func (p *Response) Error() error {
	return p.err
}

func (p *Response) Body() []byte {
	return p.body
}

func (p *Response) String() string {
	return string(p.body)
}

func (p *Response) Parse(parser Parser) (data interface{}, err error) {
	return parser(p.body)
}

func (p *Response) Json() (*_json.Json, error) {
	if p.err != nil {
		return nil, p.err
	}
	if len(p.body) == 0 {
		return nil, fmt.Errorf("response body is empty")
	}
	return _json.NewJson(p.body)
}

func (p *Response) Name() string {
	return p.request.name
}

type Responses struct {
	err   error
	lock  sync.Mutex
	_list []*Response
	_map  map[string]*Response
}

func (p *Responses) Error() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.err
}

func (p *Responses) Add(response *Response) {
	if response == nil {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	if p._map == nil {
		p._map = map[string]*Response{}
	}

	_, ok := p._map[response.request.name]
	if ok {
		if p.err == nil {
			p.err = fmt.Errorf("name '%s' duplicated", response.request.name)
		}
		return
	}

	p._list = append(p._list, response)
	p._map[response.request.name] = response
	return
}

func (p *Responses) Get(name string) *Response {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p._map == nil {
		return nil
	}
	return p._map[name]
}

func (p *Responses) All() []*Response {
	p.lock.Lock()
	defer p.lock.Unlock()
	return append(make([]*Response, 0), p._list...)
}
