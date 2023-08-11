package main

import (
	interaction "GoYin/server/kitex_gen/interaction/interactionserver"
	"GoYin/server/service/interaction/config"
	"GoYin/server/service/interaction/dao"
	"GoYin/server/service/interaction/initialize"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
)

func main() {
	initialize.InitLogger()
	r, info := initialize.InitNacos()
	db := initialize.InitDB()
	rdb := initialize.InitRedis()
	impl := &InteractionServerImpl{
		RedisManager:    dao.NewRedisManager(rdb),
		CommentManager:  dao.NewMysqlManager(db),
		FavoriteManager: dao.NewMysqlManager(db),
	}
	svr := interaction.NewServer(impl,
		server.WithServiceAddr(utils.NewNetAddr("tcp", net.JoinHostPort(config.GlobalServerConfig.Host, config.GlobalServerConfig.Port))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.Name}))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
