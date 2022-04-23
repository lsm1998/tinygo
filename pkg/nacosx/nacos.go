package nacosx

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/cast"
	"strings"
)

type Options struct {
	conf Config
}

type Option func(*Options)

// Config nacos的配置字段
type Config struct {
	Endpoints           []string `json:"endpoints" yaml:"endpoints"`
	NamespaceId         string   `json:"namespace_id" yaml:"namespace_id"`
	Username            string   `json:"username" yaml:"timeout"`
	Password            string   `json:"password" yaml:"password"`
	LogDir              string   `json:"log_dir" yaml:"log_dir"`
	LogLevel            string   `json:"log_level" yaml:"log_level"`
	TimeoutMs           uint64   `json:"timeout_ms" yaml:"timeout_ms"`
	NotLoadCacheAtStart bool     `json:"not_load_cache_at_start" yaml:"not_load_cache_at_start"`
}

func WithConfig(conf Config) Option {
	return func(k *Options) {
		k.conf = conf
	}
}

func Must(opts ...Option) config_client.IConfigClient {
	e := &Options{}
	for _, opt := range opts {
		opt(e)
	}
	var serverConfigs []constant.ServerConfig
	for _, v := range e.conf.Endpoints {
		split := strings.Split(v, ":")
		if len(split) != 2 {
			panic("nacosx Config Endpoints split error")
		}
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr: split[0],
			Port:   cast.ToUint64(split[1]),
		})
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         e.conf.NamespaceId,
		TimeoutMs:           e.conf.TimeoutMs,
		NotLoadCacheAtStart: e.conf.NotLoadCacheAtStart,
		Username:            e.conf.Username,
		Password:            e.conf.Password,
		LogDir:              e.conf.LogDir,
		LogLevel:            e.conf.LogLevel,
	}
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
	return configClient
}
