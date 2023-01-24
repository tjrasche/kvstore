package membership

type Handler interface {
	Members() ([]Node, error)
	CalculateReplication(key string, value string) []Node
}
