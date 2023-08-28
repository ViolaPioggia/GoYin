package config

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Salt     string `mapstructure:"salt" json:"salt"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type ServerConfig struct {
	Name               string               `mapstructure:"name" json:"name"`
	Host               string               `mapstructure:"host" json:"host"`
	Port               string               `mapstructure:"port" json:"port"`
	JWTInfo            JWTConfig            `mapstructure:"jwt" json:"jwt"`
	MysqlInfo          MysqlConfig          `mapstructure:"mysql" json:"mysql"`
	RedisInfo          RedisConfig          `mapstructure:"redis" json:"redis"`
	OtelInfo           OtelConfig           `mapstructure:"otel" json:"otel"`
	SocialitySrvInfo   RPCSrvConfig         `mapstructure:"sociality_srv" json:"sociality_srv"`
	ChatSrvInfo        RPCSrvConfig         `mapstructure:"chat_srv" json:"chat_srv"`
	InteractionSrvInfo RPCSrvConfig         `mapstructure:"interaction_srv" json:"interaction_srv"`
	CbRule             CircuitBreakerConfig `mapstructure:"cb_rule" json:"cb_rule"`
}

type RPCSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type CircuitBreakerConfig struct {
	Resource                     string  `json:"Resource"`
	Strategy                     uint32  `json:"Strategy"`
	RetryTimeoutMs               uint32  `json:"RetryTimeoutMs"`
	MinRequestAmount             uint64  `json:"MinRequestAmount"`
	StatIntervalMs               uint32  `json:"StatIntervalMs"`
	StatSlidingWindowBucketCount uint32  `json:"StatSlidingWindowBucketCount"`
	MaxAllowedRtMs               uint64  `json:"MaxAllowedRtMs"`
	Threshold                    float64 `json:"Threshold"`
}
