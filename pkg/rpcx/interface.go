package rpcx

type Register interface {
	Discov(serviceName string, port string) <-chan error
}
