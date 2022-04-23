package configx

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

const (
	Local ConfigType = "config-local"
	Redis ConfigType = "config-redis"
)

const (
	bootstrapFile = "bootstrap.yaml"
)

type ConfigType string

func LoadConfig(key string, typ ConfigType, obj interface{}) error {
	switch typ {
	case Local:
		bytes, err := ioutil.ReadFile(key)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(bytes, obj)
	case Redis:
		c, err := getConfig()
		if err != nil {
			return err
		}
		return redisParse(c, key, obj)
	default:
		panic("not support ConfigType," + typ)
	}
	return nil
}

func WatchConfig(key string, handler func(configVal string, err error)) {
	var oldVal string
	for {
		time.Sleep(5 * time.Second)
		temp, err := get(key)
		if err != nil {
			handler("", err)
			continue
		}
		if oldVal == "" {
			oldVal = temp
		} else if oldVal != temp {
			handler(temp, err)
		}
	}
}

func getConfig() (config, error) {
	var c config
	bytes, err := ioutil.ReadFile(bootstrapFile)
	if err != nil {
		return c, err
	}
	return c, yaml.Unmarshal(bytes, &c)
}

type config struct {
	Addr string `json:"addr" yaml:"addr"`
	Port int    `json:"port" yaml:"port"`
	Auth string `json:"auth" yaml:"auth"`
	Db   int    `json:"db" yaml:"db"`
}
