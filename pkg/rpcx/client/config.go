package client

type RpcClientConf struct {
	Addr        string `yaml:"addr"`
	ServiceName string `yaml:"service_name"`
}
