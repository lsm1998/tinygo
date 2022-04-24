package configx

import (
	"github.com/lsm1998/tinygo/pkg/nacosx"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var configClient config_client.IConfigClient

type nacosConfig struct {
	nacosx.Config `json:",inline" yaml:",inline"`
	Group         string `json:"group" yaml:"group"`
	DataId        string `json:"data_id" yaml:"data_id"`
}

func nacosParse(c string, key string, obj interface{}, watch ...func(string, error)) error {
	config := nacosConfig{}
	err := yaml.Unmarshal([]byte(c), &config)
	if err != nil {
		return err
	}
	if config.DataId == "" {
		config.DataId = key
	}
	configClient = nacosx.Must(nacosx.WithConfig(config.Config))
	content, err := configClient.GetConfig(vo.ConfigParam{
		Group:  config.Group,
		DataId: config.DataId,
	})
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal([]byte(content), obj); err != nil {
		return err
	}
	if len(watch) == 0 {
		return nil
	}
	go nacosWatch(config, watch...)
	return nil
}

func nacosWatch(c nacosConfig, watch ...func(string, error)) {
	if err := configClient.ListenConfig(vo.ConfigParam{
		DataId: c.DataId,
		Group:  c.Group,
		OnChange: func(namespace, group, dataId, data string) {
			for _, v := range watch {
				v(data, nil)
			}
			logrus.Debugf("nacosWatch OnChange,data=%v", data)
		},
	}); err != nil {
		logrus.Errorf("nacosWatch fail,err=%v", err)
		for _, v := range watch {
			v("", nil)
		}
	}
}
