package dao

import (
	"GoYin/server/service/interaction/model"
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisManager struct {
	redisClient *redis.Client
}

func (r RedisManager) FavoriteAction(ctx context.Context, userId, videoId int64) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) UnFavoriteAction(ctx context.Context, userId, videoId int64) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) GetFavoriteVideoIdList(ctx context.Context, userId int64) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) Comment(ctx context.Context, comment *model.Comment) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) DeleteComment(ctx context.Context, commentId int64) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) GetComment(ctx context.Context, videoId int64) ([]*model.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) GetFavoriteCount(ctx context.Context, videoId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) GetCommentCount(ctx context.Context, videoId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisManager) JudgeIsFavoriteCount(ctx context.Context, videoId, userId int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewRedisManager(client *redis.Client) *RedisManager {
	return &RedisManager{redisClient: client}
}
