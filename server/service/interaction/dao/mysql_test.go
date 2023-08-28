package dao

import (
	"GoYin/server/service/interaction/config"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"testing"
)

func BenchmarkMysqlManager_GetFavoriteCountByVideoId(b *testing.B) {
	m := GetMysqlManager()
	for n := 0; n < b.N; n++ {
		m.GetFavoriteCountByVideoId(1693551100440887296)
	}
}

func BenchmarkMysqlManager_GetFavoriteVideoCountByUserId(b *testing.B) {
	m := GetMysqlManager()
	for n := 0; n < b.N; n++ {
		m.GetFavoriteVideoCountByUserId(1693512896484478976)
	}
}

func BenchmarkMysqlManager_GetFavoriteVideoIdList(b *testing.B) {
	m := GetMysqlManager()
	for n := 0; n < b.N; n++ {
		m.GetFavoriteVideoIdList(context.TODO(), 1693512896484478976)
	}
}

func BenchmarkMysqlManager_GetFavoriteCount(b *testing.B) {
	m := GetMysqlManager()
	for n := 0; n < b.N; n++ {
		m.GetFavoriteCount(context.TODO(), 1693551100440887296)
	}
}

func BenchmarkMysqlManager_JudgeIsFavoriteCount(b *testing.B) {
	m := GetMysqlManager()
	for n := 0; n < b.N; n++ {
		m.JudgeIsFavoriteCount(context.TODO(), 1693551100440887296, 1693612053983399936)
	}
}

func BenchmarkMysqlManager_GetComment(b *testing.B) {
	m := GetMysqlManager()
	for n := 0; n < b.N; n++ {
		m.GetComment(context.TODO(), 1693551100440887296)
	}
}

func BenchmarkMysqlManager_GetCommentCount(b *testing.B) {
	m := GetMysqlManager()
	for n := 0; n < b.N; n++ {
		m.GetCommentCount(context.TODO(), 1693551100440887296)
	}
}

func GetMysqlManager() *MysqlManager {
	InitNacos()
	c := config.GlobalServerConfig.MysqlInfo
	db, _ := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Name)),
		&gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	m := MysqlManager{
		commentDb:  db,
		favoriteDb: db,
	}
	return &m
}

func InitNacos() {
	v := viper.New()
	v.SetConfigFile("../config.yaml")
	if err := v.ReadInConfig(); err != nil {
		klog.Fatalf("read viper config failed: %s", err)
	}
	if err := v.Unmarshal(&config.GlobalNacosConfig); err != nil {
		klog.Fatalf("unmarshal err failed: %s", err)
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(config.GlobalNacosConfig.Host, config.GlobalNacosConfig.Port),
	}

	cc := constant.ClientConfig{
		NamespaceId:         config.GlobalNacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		Username:            config.GlobalNacosConfig.User,
		Password:            config.GlobalNacosConfig.Password,
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		klog.Fatalf("create config client failed: %s", err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: config.GlobalNacosConfig.DataId,
		Group:  config.GlobalNacosConfig.Group,
	})
	if err != nil {
		klog.Fatalf("get config failed: %s", err.Error())
	}

	err = sonic.Unmarshal([]byte(content), &config.GlobalServerConfig)
	if err != nil {
		klog.Fatalf("nacos config failed: %s", err)
	}
}
