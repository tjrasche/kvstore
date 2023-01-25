package membership

import (
	"bytes"
	"io"
	"net/http"
)

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

	_, err := n.c.Post(n.address+"/replicate/"+key, "Content-Type: text/plain", bytes.NewBuffer([]byte(value)))

	if err != nil {
		return err
	}

	return nil
}

func (n Node) Get(key string) (value string, err error) {
	result, err := n.c.Get(n.address + "/" + key)

	if err != nil {
		return "", err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (n Node) Store(key string, value string) {
	panic("not implemented")

}
