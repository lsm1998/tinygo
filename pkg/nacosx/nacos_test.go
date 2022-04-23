package nacosx

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"testing"
)

func TestServerConfig(t *testing.T) {
	configClient := Must(WithConfig(Config{
		Endpoints:   []string{"127.0.0.1:8848"},
		Username:    "nacos",
		Password:    "nacos",
		NamespaceId: "5af3b7b4-cee9-4fff-b96c-d724b3824317",
	}))

	//获取配置信息
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "demo",
		Group:  "DEFAULT_GROUP"})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(content)

	//监听配置
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "demo",
		Group:  "DEFAULT_GROUP",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	select {}
}
