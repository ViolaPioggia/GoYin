package initialize

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/user/config"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/utils"
	nacos "github.com/kitex-contrib/registry-nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"net"
)

func InitNacos() (registry.Registry, *registry.Info) {
	v := viper.New()
	v.SetConfigFile(consts.UserConfigPath)
	if err := v.ReadInConfig(); err != nil {
		klog.Fatalf("read viper config failed: %s", err)
	}
	if err := v.Unmarshal(&config.GlobalNacosConfig); err != nil {
		klog.Fatalf("unmarshal err failed: %s", err)
	}
	klog.Infof("Config Info: %v", config.GlobalNacosConfig)

	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := constant.ClientConfig{
		NamespaceId:         config.GlobalNacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "info",
		Username:            config.GlobalNacosConfig.Namespace,
		Password:            config.GlobalNacosConfig.Namespace,
	}

	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		klog.Errorf("create registry err: %s", err.Error())
	}

	r := nacos.NewNacosRegistry(cli, nacos.WithGroup("user_group")) //group?

	sf, err := snowflake.NewNode(consts.NacosSnowflakeNode)
	if err != nil {
		klog.Fatalf("generate nacos service name failed:%s", err)
	}
	info := &registry.Info{
		ServiceName: config.GlobalServerConfig.Name,
		Addr:        utils.NewNetAddr("tcp", net.JoinHostPort(config.GlobalServerConfig.Host, config.GlobalServerConfig.Port)),
		Tags:        map[string]string{"ID": sf.Generate().Base36()},
	}
	return r, info
}
