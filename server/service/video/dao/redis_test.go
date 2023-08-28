package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/video/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

func BenchmarkRedisManager_GetBasicVideoListByLatestTime(b *testing.B) {
	m := GetRedisManager()
	for n := 0; n < b.N; n++ {
		m.GetBasicVideoListByLatestTime(context.TODO(), 0, time.Now().Unix())
	}
}

func BenchmarkRedisManager_GetPublishedVideoListByUserId(b *testing.B) {
	m := GetRedisManager()
	for n := 0; n < b.N; n++ {
		m.GetPublishedVideoListByUserId(context.TODO(), 1693512896484478976)
	}
}

func BenchmarkRedisManager_GetFavoriteVideoListByUserId(b *testing.B) {
	m := GetRedisManager()
	for n := 0; n < b.N; n++ {
		m.GetFavoriteVideoListByUserId(context.TODO(), 1693512896484478976)
	}
}

func BenchmarkRedisManager_GetPublishedVideoIdListByUserId(b *testing.B) {
	m := GetRedisManager()
	for n := 0; n < b.N; n++ {
		m.GetPublishedVideoIdListByUserId(context.TODO(), 1693512896484478976)
	}
}

func BenchmarkRedisManager_GetVideoByVideoId(b *testing.B) {
	m := GetRedisManager()
	for n := 0; n < b.N; n++ {
		m.GetVideoByVideoId(context.TODO(), 1693551100440887296)
	}
}

func GetRedisManager() *RedisManager {
	InitNacos()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.GlobalServerConfig.RedisInfo.Host, config.GlobalServerConfig.RedisInfo.Port),
		Password: config.GlobalServerConfig.RedisInfo.Password,
		DB:       consts.RedisUserClientDB,
	})
	return &RedisManager{redisClient: client}
}
