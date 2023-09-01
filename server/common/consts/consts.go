package consts

const (
	UserSnowflakeNode    = 1
	NacosSnowflakeNode   = 2
	CommentSnowFlakeNode = 3
	VideoSnowFlakeNode   = 4
	MinioSnowFlakeNode   = 5

	MysqlAlreadyExists = "useralreadyexists"

	KlogFilePath                = "./tmp/klog/logs"
	HlogFilePath                = "./tmp/hlog/logs"
	UserSentinelFilePath        = "./tmp/circuit/user"
	ChatSentinelFilePath        = "./tmp/circuit/chat"
	InteractionSentinelFilePath = "./tmp/circuit/interaction"
	SocialSentinelFilePath      = "./tmp/circuit/social"
	VideoSentinelFilePath       = "./tmp/circuit/video"
	ApiSentinelFilePath         = "./tmp/flow/api"

	NacosLogDir   = "tmp/nacos/log"
	NacosCacheDir = "tmp/nacos/cache"
	NacosLogLevel = "debug"

	UserConfigPath        = "./server/service/user/config.yaml"
	SocialityConfigPath   = "./server/service/sociality/config.yaml"
	InteractionConfigPath = "./server/service/interaction/config.yaml"
	VideoConfigPath       = "./server/service/video/config.yaml"
	ChatConfigPath        = "./server/service/chat/config.yaml"
	ApiConfigPath         = "./server/service/api/config.yaml"

	RedisSocialClientDB   = 1
	RedisUserClientDB     = 2
	RedisVideoClientDB    = 3
	RedisCommentClientDB  = 4
	RedisFavoriteClientDB = 5

	FollowList   = 0
	FollowerList = 1
	FriendsList  = 2

	Follow   = 1
	UnFollow = 2

	IsLike = 1
	Like   = 1
	UnLike = 2

	DeleteComment = 2
	Comment       = 1

	SentMessage    = 1
	ReceiveMessage = 0

	NsqChatTopic        = "chat"
	NsqVideoTopic       = "video"
	NsqSocialityTopic   = "sociality"
	NsqChatChannel      = "1"
	NsqVideoChannel     = "2"
	NsqSocialityChannel = "3"

	MySqlDSN           = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	MySQLImage         = "mysql:latest"
	MySQLContainerPort = "3306/tcp"
	MySQLContainerIP   = "127.0.0.1"
	MySQLPort          = "0"
	MySQLAdmin         = "root"
	DockerTestMySQLPwd = "123456"
	DockerTestMySQLDb  = "GoYin"

	RedisImage         = "redis:latest"
	RedisContainerPort = "6379/tcp"
	RedisContainerIP   = "127.0.0.1"
	RedisPort          = "0"
)
