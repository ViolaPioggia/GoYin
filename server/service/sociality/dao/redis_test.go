package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/sociality/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkRedisManager_GetUserIdList(b *testing.B) {
	m := GetReidsManager()
	rand.Seed(time.Now().Unix())
	for n := 0; n < b.N; n++ {
		m.GetUserIdList(context.TODO(), 1693512896484478976, int8(rand.Intn(12)%3))
	}
}

func BenchmarkRedisManager_GetSocialInfo(b *testing.B) {
	m := GetReidsManager()
	for n := 0; n < b.N; n++ {
		m.GetSocialInfo(context.TODO(), 1693512896484478976, 1693526693303554048)
	}
}

func BenchmarkRedisManager_BatchGetSocialInfo(b *testing.B) {
	m := GetReidsManager()
	userId := []int64{1693526693303554048, 1693526693303554048, 1693526693303554048}
	for n := 0; n < b.N; n++ {
		m.BatchGetSocialInfo(context.TODO(), userId, 1693512896484478976)
	}
}

func GetReidsManager() *RedisManager {
	InitNacos()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.GlobalServerConfig.RedisInfo.Host, config.GlobalServerConfig.RedisInfo.Port),
		Password: config.GlobalServerConfig.RedisInfo.Password,
		DB:       consts.RedisUserClientDB,
	})
	return &RedisManager{
		redisClient: client,
	}
}
