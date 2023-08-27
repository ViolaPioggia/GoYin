package main

import (
	chat "GoYin/server/kitex_gen/chat/chatservice"
	"GoYin/server/service/chat/config"
	"GoYin/server/service/chat/dao"
	"GoYin/server/service/chat/initialize"
	"GoYin/server/service/chat/pkg"
	"context"
	"errors"
	"fmt"
	kitexSentinel "github.com/alibaba/sentinel-golang/pkg/adapters/kitex"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"net"
)

func main() {
	initialize.InitLogger()
	r, info := initialize.InitNacos()
	initialize.Sentinel()
	db := initialize.InitDB()
	publisher := initialize.InitProducer()
	subscriber := initialize.InitSubscriber()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	go func() {
		err := pkg.SubscriberManager.Subscribe(*pkg.NewSubscriberManager(subscriber), context.Background(), dao.NewMysqlManager(db))
		if err != nil {
			klog.Error(err)
		}
	}()
	impl := &ChatServiceImpl{
		MysqlManager: dao.NewMysqlManager(db),
		Publisher:    pkg.NewPublishManager(publisher),
		Subscriber:   pkg.NewSubscriberManager(subscriber),
	}
	svr := chat.NewServer(impl,
		server.WithServiceAddr(utils.NewNetAddr("tcp", net.JoinHostPort(config.GlobalServerConfig.Host, config.GlobalServerConfig.Port))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithMiddleware(kitexSentinel.SentinelServerMiddleware(
			kitexSentinel.WithResourceExtract(func(ctx context.Context, req, resp interface{}) string {
				return config.GlobalServerConfig.CbRule.Resource
			}),
			kitexSentinel.WithBlockFallback(func(ctx context.Context, req, resp interface{}, blockErr error) error {
				return errors.New("service block")
			}),
		)),
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
		fmt.Println(err)
	}
}
