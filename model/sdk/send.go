package sdk

type Config struct {
	Base
	Server
	Path string
}

type Base struct {
	From    []string
	To      []string
	Subject []string
	Body    string
}

type Server struct {
	Username string
	Password string
	Host     string
	Port     int
}
