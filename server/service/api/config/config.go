package config

import "github.com/ip2location/ip2location-go/v9"

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
type MinioConfig struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id" json:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key" json:"secret_access_key"`
	Bucket          string `mapstructure:"bucket" json:"bucket"`
	UrlPrefix       string `mapstructure:"url_prefix" json:"url_prefix"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}

type ServerConfig struct {
	Name               string          `mapstructure:"name" json:"name"`
	Host               string          `mapstructure:"host" json:"host"`
	Port               int             `mapstructure:"port" json:"port"`
	MinioInfo          MinioConfig     `mapstructure:"minio" json:"minio"`
	JWTInfo            JWTConfig       `mapstructure:"jwt" json:"jwt"`
	OtelInfo           OtelConfig      `mapstructure:"otel" json:"otel"`
	ChatSrvInfo        RPCSrvConfig    `mapstructure:"chat_srv" json:"chat_srv"`
	UserSrvInfo        RPCSrvConfig    `mapstructure:"user_srv" json:"user_srv"`
	InteractionSrvInfo RPCSrvConfig    `mapstructure:"interaction_srv" json:"interaction_srv"`
	SocialitySrvInfo   RPCSrvConfig    `mapstructure:"sociality_srv" json:"sociality_srv"`
	VideoSrvInfo       RPCSrvConfig    `mapstructure:"video_srv" json:"video_srv"`
	IpInfo             *ip2location.DB `mapstructure:"ip" json:"ip"`
	FlowRule           FlowRule        `mapstructure:"flow_rule" json:"flow_rule"`
}

type RPCSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type FlowRule struct {
	Resource               string `json:"resource"`
	Threshold              uint32 `json:"threshold"`
	TokenCalculateStrategy int32  `json:"token_calculate_strategy"`
	ControlBehavior        int32  `json:"control_behavior"`
	StatIntervalInMs       uint32 `json:"stat_interval_in_ms"`
}
