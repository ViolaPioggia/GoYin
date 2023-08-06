package dao

import (
	"GoYin/server/service/user/models"
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisManager struct {
	redisClient *redis.Client
}

func NewRedisManager(client *redis.Client) *RedisManager {
	return &RedisManager{redisClient: client}
}

func (r *RedisManager) CreateUser(ctx context.Context, user *models.User) error {
	return nil
}

func (r *RedisManager) GetUserById(ctx context.Context, id int64) (*models.User, error) {
	return nil, nil
}

func (r *RedisManager) BatchGetUserById(ctx context.Context, id []int64) ([]*models.User, error) {
	return nil, nil
}
