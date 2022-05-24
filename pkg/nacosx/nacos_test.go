package nacosx

import (
	"fmt"
	"github.com/lsm1998/tinygo"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"testing"
	"time"
)

/***
api https://github.com/nacos-group/nacos-sdk-go/blob/master/README_CN.md
*/

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

func TestRegisterInstance(t *testing.T) {
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "console1.nacos.io",
			ContextPath: "/nacos",
			Port:        80,
			Scheme:      "http",
		},
		{
			IpAddr:      "console2.nacos.io",
			ContextPath: "/nacos",
			Port:        80,
			Scheme:      "http",
		},
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         "5af3b7b4-cee9-4fff-b96c-d724b3824317", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "error",
		Username:            "nacos",
		Password:            "nacos",
	}

	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	addr, _ := tinygo.IpAddr()
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          addr,
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
		ClusterName: "cluster-a", // 默认值DEFAULT
		GroupName:   "group-a",   // 默认值DEFAULT_GROUP
	})
	fmt.Println("success:", success)
	time.Sleep(10 * time.Second)
	success, err = namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          addr,
		Port:        8848,
		ServiceName: "demo.go",
		Ephemeral:   true,
		Cluster:     "cluster-a", // 默认值DEFAULT
		GroupName:   "group-a",   // 默认值DEFAULT_GROUP
	})
	fmt.Println("success:", success)
}

func TestGetService(t *testing.T) {
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "console1.nacos.io",
			ContextPath: "/nacos",
			Port:        80,
			Scheme:      "http",
		},
		{
			IpAddr:      "console2.nacos.io",
			ContextPath: "/nacos",
			Port:        80,
			Scheme:      "http",
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         "5af3b7b4-cee9-4fff-b96c-d724b3824317", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "error",
		Username:            "nacos",
		Password:            "nacos",
	}

	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// 获取一个
	services, err := namingClient.GetService(vo.GetServiceParam{
		ServiceName: "demo.go",
		Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
		GroupName:   "group-a",             // 默认值DEFAULT_GROUP
	})
	fmt.Println(services)

	// SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0
	instances, err := namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",             // 默认值DEFAULT_GROUP
		Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
	})
	fmt.Println(instances)

	// SelectOneHealthyInstance将会按加权随机轮询的负载均衡策略返回一个健康的实例
	// 实例必须满足的条件：health=true,enable=true and weight>0
	instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",             // 默认值DEFAULT_GROUP
		Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
	})
	fmt.Println(instance)

	// Subscribe key=serviceName+groupName+cluster
	// 注意:我们可以在相同的key添加多个SubscribeCallback.
	err = namingClient.Subscribe(&vo.SubscribeParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",             // 默认值DEFAULT_GROUP
		Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			fmt.Println(services, err)
		},
	})

}
