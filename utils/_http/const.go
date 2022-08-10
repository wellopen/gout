package _http

type Method string

const (
	POST   Method = "POST"
	GET    Method = "GET"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	HEAD   Method = "HEAD"
)
const (
	CONTENT_TYPE          = "Content-Type"
	JSON                  = "application/json"
	JSON_UTF8             = "application/json; charset=utf-8"
	X_WWW_FORM_URLENCODED = "application/x-www-form-urlencoded"
	FORM_DATA             = "multipart/form-data"
	TEXT_XML              = "text/xml"
)
