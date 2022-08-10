package _http

type Parser func(data []byte) (result interface{}, err error)
