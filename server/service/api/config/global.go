package config

import (
	"GoYin/server/kitex_gen/chat/chatservice"
	"GoYin/server/kitex_gen/interaction/interactionserver"
	"GoYin/server/kitex_gen/sociality/socialityservice"
	"GoYin/server/kitex_gen/user/userservice"
	"GoYin/server/kitex_gen/video/videoservice"
	"github.com/minio/minio-go/v7"
)

var (
	GlobalServerConfig = &ServerConfig{}
	GlobalNacosConfig  = &NacosConfig{}

	GlobalChatClient        chatservice.Client
	GlobalUserClient        userservice.Client
	GlobalVideoClient       videoservice.Client
	GlobalSocialClient      socialityservice.Client
	GlobalInteractionClient interactionserver.Client
	GlobalMinioClient       *minio.Client
)
