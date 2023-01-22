package membership

import "net/http"

type Node struct {
	address string
	c       http.Client
}

func (n Node) Replicate(key string, value string) error {
	panic("not implemented")
}

func (n Node) Get(key string) (err error, value string) {
	panic("not implemented")
}
