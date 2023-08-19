package dao

import (
	"GoYin/server/service/interaction/model"
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type RedisManager struct {
	redisClient *redis.Client
}

func (r RedisManager) FavoriteAction(ctx context.Context, userId, videoId int64) error {
	err := r.redisClient.LPush(ctx, "user_video_id:"+strconv.FormatInt(userId, 10), videoId).Err()
	if err != nil {
		klog.Error("redis lPush failed,", err)
		return err
	}
	err = r.redisClient.LPush(ctx, "video_user_id:"+strconv.FormatInt(videoId, 10), userId).Err()
	if err != nil {
		klog.Error("redis lPush failed,", err)
		return err
	}
	return nil
}

func (r RedisManager) UnFavoriteAction(ctx context.Context, userId, videoId int64) error {
	err := r.redisClient.LRem(ctx, "user_video_id:"+strconv.FormatInt(userId, 10), 0, videoId).Err()
	if err != nil {
		klog.Error("redis lRem failed,", err)
		return err
	}
	err = r.redisClient.LRem(ctx, "video_user_id:"+strconv.FormatInt(videoId, 10), 0, userId).Err()
	if err != nil {
		klog.Error("redis lRem failed,", err)
		return err
	}
	return nil
}

func (r RedisManager) GetFavoriteVideoIdList(ctx context.Context, userId int64) ([]int64, error) {
	res, err := r.redisClient.LRange(ctx, "user_video_id:"+strconv.FormatInt(userId, 10), 0, -1).Result()
	if err != nil {
		klog.Error("redis get favorite videoIdList failed,", err)
		return nil, err
	}
	var idList []int64
	var id int64
	for _, v := range res {
		id, err = strconv.ParseInt(v, 0, 64)
		if err != nil {
			klog.Error("redis transform string into int64 failed,", err)
			return nil, err
		}
		idList = append(idList, id)
	}
	return idList, nil
}

func (r RedisManager) Comment(ctx context.Context, comment *model.Comment) error {
	commentJson, err := sonic.Marshal(comment)
	if err != nil {
		klog.Error("redis marshal comment failed,", err)
		return err
	}
	err = r.redisClient.LPush(ctx, "video_comment:"+strconv.FormatInt(comment.VideoId, 10), commentJson).Err()
	if err != nil {
		klog.Error("redis insert comment failed,", err)
		return err
	}
	err = r.redisClient.Set(ctx, "comment:"+strconv.FormatInt(comment.ID, 10), commentJson, 0).Err()
	if err != nil {
		klog.Error("redis insert comment failed,", err)
		return err
	}
	return nil
}

func (r RedisManager) DeleteComment(ctx context.Context, commentId int64) error {
	commentJson, err := r.redisClient.Get(ctx, "comment:"+strconv.FormatInt(commentId, 10)).Bytes()
	if err != nil && err != redis.Nil {
		klog.Error("redis get commentJson failed,", err)
		return err
	}
	r.redisClient.Del(ctx, "comment:"+strconv.FormatInt(commentId, 10))
	var comment *model.Comment
	err = sonic.Unmarshal(commentJson, &comment)
	if err != nil {
		klog.Error("redis unmarshal comment failed,", err)
		return err
	}
	err = r.redisClient.LRem(ctx, "video_comment:"+strconv.FormatInt(comment.VideoId, 10), 0, commentJson).Err()
	if err != nil && err != redis.Nil {
		klog.Error("redis delete comment failed,", err)
		return err
	}
	return nil
}

func (r RedisManager) GetComment(ctx context.Context, videoId int64) ([]*model.Comment, error) {
	res, err := r.redisClient.LRange(ctx, "video_comment:"+strconv.FormatInt(videoId, 10), 0, -1).Result()
	if err != nil {
		klog.Error("redis get comment failed,", err)
		return nil, err
	}
	var commentList []*model.Comment
	for _, v := range res {
		var comment *model.Comment
		err = sonic.UnmarshalString(v, &comment)
		if err != nil {
			klog.Error("redis unmarshal comment failed,", err)
			return nil, err
		}
		commentList = append(commentList, comment)
	}
	return commentList, nil
}

func (r RedisManager) GetFavoriteCount(ctx context.Context, videoId int64) (int64, error) {
	count, err := r.redisClient.LLen(ctx, "video_user_id:"+strconv.FormatInt(videoId, 10)).Result()
	if err != nil {
		klog.Error("redis get favorite count failed,", err)
		return 0, err
	}
	return count, nil
}

func (r RedisManager) GetCommentCount(ctx context.Context, videoId int64) (int64, error) {
	lenth, err := r.redisClient.LLen(ctx, "video_comment:"+strconv.FormatInt(videoId, 10)).Result()
	if err != nil {
		klog.Error("redis get comment count failed,", err)
		return 0, err
	}
	return lenth, nil
}

func (r RedisManager) JudgeIsFavoriteCount(ctx context.Context, videoId, userId int64) (bool, error) {
	res, err := r.redisClient.LRange(ctx, "user_video_id:"+strconv.FormatInt(userId, 10), 0, -1).Result()
	if err != nil {
		klog.Error("redis get user_video failed,", err)
		return false, err
	}
	for _, v := range res {
		if v == strconv.FormatInt(videoId, 10) {
			return true, nil
		}
	}
	return false, nil
}

func NewRedisManager(client *redis.Client) *RedisManager {
	return &RedisManager{redisClient: client}
}
