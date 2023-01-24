package membership

import "net/http"

type Node struct {
	Id      string
	address string
	c       http.Client
}

func NewNode(id string, address string) *Node {
	return &Node{
		Id:      id,
		address: address,
		c:       http.Client{},
	}
}

func (n Node) Replicate(key string, value string) error {
	panic("not implemented")
}

func (n Node) Get(key string) (err error, value string) {
	panic("not implemented")
}

func (n Node) Store(key string, value string) {
	panic("not implemented")

}
