package consts

const (
	UserSnowflakeNode    = 1
	NacosSnowflakeNode   = 2
	CommentSnowFlakeNode = 3
	VideoSnowFlakeNode   = 4

	MysqlAlreadyExists = "useralreadyexists"

	KlogFilePath = "./tmp/klog/logs"
	HlogFilePath = "./tmp/hlog/logs"

	NacosLogDir   = "tmp/nacos/log"
	NacosCacheDir = "tmp/nacos/cache"
	NacosLogLevel = "debug"

	UserConfigPath = "./server/service/user/config.yaml"
	ApiConfigPath  = "./server/service/api/config.yaml"

	RedisUserClientDB = 1

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

	NsqChatTopic    = "chat"
	NsqVideoTopic   = "video"
	NsqChatChannel  = "channel 1"
	NsqVideoChannel = "channel 2"
)
