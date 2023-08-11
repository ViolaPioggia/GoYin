package main

import (
	sociality "GoYin/server/kitex_gen/sociality/socialityservice"
	"GoYin/server/service/sociality/config"
	"GoYin/server/service/sociality/dao"
	"GoYin/server/service/sociality/initialize"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"net"

	"log"
)

func main() {
	initialize.InitLogger()
	r, info := initialize.InitNacos()
	db := initialize.InitDB()
	rdb := initialize.InitRedis()
	impl := &SocialityServiceImpl{
		RedisManager: dao.NewRedisManager(rdb),
		MysqlManager: dao.NewMysqlManager(db),
	}
	svr := sociality.NewServer(impl,
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
