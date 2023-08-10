package main

import (
	"GoYin/server/common/middleware"
	user "GoYin/server/kitex_gen/user/userservice"
	"GoYin/server/service/user/config"
	"GoYin/server/service/user/dao"
	"GoYin/server/service/user/initialize"
	"GoYin/server/service/user/pkg"
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
	socialClient := initialize.InitSocial()
	interactionClient := initialize.InitInteraction()
	impl := &UserServiceImpl{
		Jwt:                middleware.NewJWT(config.GlobalServerConfig.Name),
		InteractionManager: pkg.NewInteractionManager(interactionClient),
		SocialManager:      pkg.NewSocialManager(socialClient),
		RedisManager:       dao.NewRedisManager(rdb),
		MysqlManager:       dao.NewUser(db),
	}
	// Create new server.
	srv := user.NewServer(impl,
		server.WithServiceAddr(utils.NewNetAddr("tcp", net.JoinHostPort(config.GlobalServerConfig.Host, config.GlobalServerConfig.Port))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.Name}),
	)

	err := srv.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
