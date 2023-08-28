package main

import (
	sociality "GoYin/server/kitex_gen/sociality/socialityservice"
	"GoYin/server/service/sociality/config"
	"GoYin/server/service/sociality/dao"
	"GoYin/server/service/sociality/initialize"
	"GoYin/server/service/sociality/pkg"
	"context"
	"errors"
	kitexSentinel "github.com/alibaba/sentinel-golang/pkg/adapters/kitex"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"net"

	"log"
)

func main() {
	initialize.InitLogger()
	r, info := initialize.InitNacos()
	initialize.Sentinel()
	db := initialize.InitDB()
	rdb := initialize.InitRedis()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	publisherClient := initialize.InitProducer()
	subscriberClient := initialize.InitSubscriber()
	go func() {
		err := pkg.SubscriberManager.Subscribe(*pkg.NewSubscriberManager(subscriberClient), context.Background(), dao.NewMysqlManager(db))
		if err != nil {
			klog.Error(err)
		}
	}()
	impl := &SocialityServiceImpl{
		Publisher:    pkg.NewPublishManager(publisherClient),
		RedisManager: dao.NewRedisManager(rdb),
		MysqlManager: dao.NewMysqlManager(db),
	}
	svr := sociality.NewServer(impl,
		server.WithServiceAddr(utils.NewNetAddr("tcp", net.JoinHostPort(config.GlobalServerConfig.Host, config.GlobalServerConfig.Port))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.Name}),
		server.WithMiddleware(kitexSentinel.SentinelServerMiddleware(
			kitexSentinel.WithResourceExtract(func(ctx context.Context, req, resp interface{}) string {
				return config.GlobalServerConfig.CbRule.Resource
			}),
			kitexSentinel.WithBlockFallback(func(ctx context.Context, req, resp interface{}, blockErr error) error {
				return errors.New("service block")
			}),
		)))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
