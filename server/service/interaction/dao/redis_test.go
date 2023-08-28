package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/interaction/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

func BenchmarkRedisManager_GetFavoriteVideoIdList(b *testing.B) {
	m := GetRedisManager()
	for n := 0; n < b.N; n++ {
		m.GetFavoriteVideoIdList(context.TODO(), 1693512896484478976)
	}
}

func BenchmarkRedisManager_GetComment(b *testing.B) {
	m := GetRedisManager()
	for n := 0; n < b.N; n++ {
		m.GetComment(context.TODO(), 1693551100440887296)
	}
}

func BenchmarkRedisManager_GetFavoriteCount(b *testing.B) {
	m := GetRedisManager()
	for n := 0; n < b.N; n++ {
		m.GetFavoriteCount(context.TODO(), 1693551100440887296)
	}
}

func BenchmarkRedisManager_GetCommentCount(b *testing.B) {
	m := GetRedisManager()
	for n := 0; n < b.N; n++ {
		m.GetCommentCount(context.TODO(), 1693551100440887296)
	}
}

func BenchmarkRedisManager_JudgeIsFavoriteCount(b *testing.B) {
	m := GetRedisManager()
	for n := 0; n < b.N; n++ {
		m.JudgeIsFavoriteCount(context.TODO(), 1693551100440887296, 1693512896484478976)
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
