package dao

import (
	"GoYin/server/service/video/model"
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisManager struct {
	redisClient *redis.Client
}

func (r RedisManager) GetBasicVideoListByLatestTime(ctx context.Context, userId, latestTime int64) ([]*model.Video, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) GetPublishedVideoListByUserId(ctx context.Context, userId int64) ([]*model.Video, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) GetFavoriteVideoListByUserId(ctx context.Context, userId int64) ([]*model.Video, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) GetPublishedVideoIdListByUserId(ctx context.Context, userId int64) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) PublishVideo(ctx context.Context, video *model.Video) error {
	//TODO implement me
	panic("implement me")
}

func NewRedisManager(client *redis.Client) *RedisManager {
	return &RedisManager{redisClient: client}
}
