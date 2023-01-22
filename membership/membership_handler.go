package membership

type Handler interface {
	Members()
	CalculateReplication(key string, value string) []string
}
