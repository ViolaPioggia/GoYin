package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/user/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

func BenchmarkRedisManager_GetUserById(b *testing.B) {
	m := GetRedisManager()
	for n := 0; n < b.N; n++ {
		m.GetUserById(context.TODO(), 1693512896484478976)
	}
}

func BenchmarkRedisManager_BatchGetUserById(b *testing.B) {
	m := GetRedisManager()
	id := []int64{1693512896484478976, 1693526693303554048, 1693602654074179584}
	for n := 0; n < b.N; n++ {
		m.BatchGetUserById(context.TODO(), id)
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
