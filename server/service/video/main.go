package main

import (
	video "GoYin/server/kitex_gen/video/videoservice"
	"GoYin/server/service/video/config"
	"GoYin/server/service/video/dao"
	"GoYin/server/service/video/pkg"
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
	"log"
	"net"

	"GoYin/server/service/video/initialize"
)

func main() {
	initialize.InitLogger()
	r, info := initialize.InitNacos()
	initialize.Sentinel()
	db := initialize.InitDB()
	rdb := initialize.InitRedis()
	publisher := initialize.InitProducer()
	subscriber := initialize.InitSubscriber()
	userClient := initialize.InitUser()
	interactionClient := initialize.InitInteraction()
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
	impl := &VideoServiceImpl{
		UserManager:        pkg.NewUserManager(userClient),
		InteractionManager: pkg.NewInteractionManager(interactionClient),
		MysqlManager:       dao.NewMysqlManager(db),
		RedisManager:       dao.NewRedisManager(rdb),
		Publisher:          pkg.NewPublishManager(publisher),
	}
	svr := video.NewServer(impl,
		server.WithServiceAddr(utils.NewNetAddr("tcp", net.JoinHostPort(config.GlobalServerConfig.Host, config.GlobalServerConfig.Port))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
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
