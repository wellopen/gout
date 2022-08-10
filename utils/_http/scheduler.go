package _http

import (
	"fmt"
	"net/http"
	"github.com/wellopen/send/utils/_rand"
	"sync"
	"time"
)

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

type Scheduler struct {
	host    string
	header  http.Header
	timeout time.Duration
}

func (p *Scheduler) SetHost(host string) {
	p.host = host
}

func (p *Scheduler) SetHeader(key, value string) {
	if p.header == nil {
		p.header = http.Header{}
	}
	p.header.Set(key, value)
}

func (p *Scheduler) SetContentType(contentType string) {
	p.SetHeader(CONTENT_TYPE, contentType)
}

func (p *Scheduler) newRequest() *Request {
	req := &Request{}
	req.host = p.host
	if p.header != nil {
		req.header = p.header.Clone()
	}
	req.timeout = p.timeout
	return req
}

func (p *Scheduler) Send(method Method, url string, params ...Reader) (response *Response) {
	return p.newRequest().Send(method, url, params...)
}

func (p *Scheduler) Post(url string, params ...Reader) *Response {
	return p.Send(POST, url, params...)
}

func (p *Scheduler) Put(url string, params ...Reader) *Response {
	return p.Send(PUT, url, params...)
}

func (p *Scheduler) Delete(url string, params ...Reader) *Response {
	return p.Send(DELETE, url, params...)
}

// Series 串行执行请求，直到请求成功，或全部失败后返回
func (p *Scheduler) Series(requests ...*Request) *Response {
	if len(requests) == 0 {
		return nil
	}
	for _, v := range requests {
		resp := v.Send(v.method, v.url, v.params...)
		if resp.Error() != nil {
			return resp
		}
	}
	return &Response{err: fmt.Errorf("all failed")}
}

// Parallel 并发执行所有请求, 返回所有结果，要求每一个 Request 拥有 name 属性，且不冲突
func (p *Scheduler) Parallel(requests ...*Request) *Responses {

	n := len(requests)

	if n == 0 {
		return nil
	}

	res := new(Responses)
	wg := &sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int, wg *sync.WaitGroup, res *Responses) {
			defer wg.Done()
			res.Add(requests[i].Send(requests[i].method, requests[i].url, requests[i].params...))
		}(i, wg, res)
	}
	wg.Wait()

	return res
}

// Rand 随机执行，根据请求数抛出随机数选择一个执行，直到请求成功，或全部失败后返回
func (p *Scheduler) Rand(requests ...*Request) *Response {
	n := len(requests)
	if n == 0 {
		return nil
	}

	for {
		m := len(requests)
		if m == 0 {
			break
		}
		i, err := _rand.Int(1, m)
		if err != nil {
			return &Response{err: fmt.Errorf("rand error: %s", err)}
		}
		v := requests[i]
		resp := v.Send(v.method, v.url, v.params...)
		if resp.Error() == nil {
			return resp
		}
		requests = append(requests[:i], requests[i+1:]...)
	}
	return &Response{err: fmt.Errorf("all failed")}
}
