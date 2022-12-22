package configx

import (
	"gopkg.in/yaml.v2"
	"os"
)

const (
	Empty            = ""
	Local ConfigType = "config-local"
	Nacos ConfigType = "config-nacos"
)

const (
	bootstrapFile = "bootstrap.yaml"
)

type ConfigType string

func LoadConfig(key string, typ ConfigType, obj interface{}, watch ...func(string, error)) error {
	switch typ {
	case Local:
		bytes, err := os.ReadFile(key)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(bytes, obj)
	case Nacos:
		c, err := getConfig()
		if err != nil {
			return err
		}
		return nacosParse(c, key, obj, watch...)
	default:
		panic("not support ConfigType," + typ)
	}
	return nil
}

func getConfig() (string, error) {
	bytes, err := os.ReadFile(bootstrapFile)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
