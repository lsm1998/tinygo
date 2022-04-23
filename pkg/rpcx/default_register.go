package rpcx

type defaultRegister struct {
}

var DefaultRegister = &defaultRegister{}

func (*defaultRegister) Discov(serviceName string, port string) <-chan error {
	return nil
}
