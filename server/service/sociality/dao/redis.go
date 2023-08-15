package dao

import (
	"GoYin/server/service/sociality/model"
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisManager struct {
	redisClient *redis.Client
}

func (r RedisManager) Action(ctx context.Context, userId, toUserId int64, actionType int8) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) GetUserIdList(ctx context.Context, userId int64, option int8) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) GetSocialInfo(ctx context.Context, userId int64) (*model.SocialInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) BatchGetSocialInfo(ctx context.Context, userId []int64) ([]*model.SocialInfo, error) {
	//TODO implement me
	panic("implement me")
}

func NewRedisManager(client *redis.Client) *RedisManager {
	return &RedisManager{redisClient: client}
}
