package dao

import (
	"GoYin/server/service/video/model"
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type RedisManager struct {
	redisClient *redis.Client
}

func (r RedisManager) GetBasicVideoListByLatestTime(ctx context.Context, userId, latestTime int64) ([]*model.Video, error) {
	res, err := r.redisClient.LRange(ctx, "video", 0, -1).Result()
	if err != nil {
		klog.Error("redis get basic videoList failed,", err)
		return nil, err
	}
	var videoList []*model.Video
	for _, v := range res {
		var video *model.Video
		err = sonic.UnmarshalString(v, &video)
		if err != nil {
			klog.Error("redis unmarshal video failed,", err)
			return nil, err
		}
		if video.CreateTime-latestTime >= 0 {
			videoList = append(videoList, video)
		} else {
			return videoList, nil
		}
	}
	return videoList, nil
}

func (r RedisManager) GetPublishedVideoListByUserId(ctx context.Context, userId int64) ([]*model.Video, error) {
	videoJson, err := r.redisClient.LRange(ctx, "user_video:"+strconv.FormatInt(userId, 10), 0, -1).Result()
	if err != nil && err != redis.Nil {
		klog.Error("redis get published videoList by userId failed,", err)
		return nil, err
	}
	var res []*model.Video
	for _, v := range videoJson {
		var video *model.Video
		err = sonic.UnmarshalString(v, &video)
		if err != nil {
			klog.Error("redis unmarshal video failed,", err)
			return nil, err
		}
		res = append(res, video)
	}
	return res, nil
}

func (r RedisManager) GetFavoriteVideoListByUserId(ctx context.Context, userId int64) ([]*model.Video, error) {
	res, err := r.redisClient.LRange(ctx, "user_video_id:"+strconv.FormatInt(userId, 10), 0, -1).Result()
	if err != nil {
		klog.Error("redis get favorite videoIdList failed,", err)
		return nil, err
	}
	var videoList []*model.Video
	var id int64
	for _, v := range res {
		id, err = strconv.ParseInt(v, 0, 64)
		if err != nil {
			klog.Error("redis transform string into int64 failed,", err)
			return nil, err
		}
		video, err := r.GetVideoByVideoId(ctx, id)
		if err != nil {
			klog.Error("redis get video by video id failed,", err)
		}
		videoList = append(videoList, video)
	}
	return videoList, nil
}

func (r RedisManager) GetPublishedVideoIdListByUserId(ctx context.Context, userId int64) ([]int64, error) {
	videoJson, err := r.redisClient.LRange(ctx, "user_video:"+strconv.FormatInt(userId, 10), 0, -1).Result()
	if err != nil && err != redis.Nil {
		klog.Error("redis get published videoList by userId failed,", err)
		return nil, err
	}
	var video *model.Video
	var idList []int64
	for _, v := range videoJson {
		err = sonic.UnmarshalString(v, &video)
		if err != nil {
			klog.Error("redis unmarshal video failed,", err)
			return nil, err
		}
		idList = append(idList, video.ID)
	}
	return idList, nil
}

func (r RedisManager) PublishVideo(ctx context.Context, video *model.Video) error {
	videoJson, err := sonic.MarshalString(video)
	if err != nil {
		klog.Error("redis marshal video failed,", err)
		return err
	}
	err = r.redisClient.LPush(ctx, "video", videoJson).Err()
	if err != nil {
		klog.Error("redis publish video failed,", err)
		return err
	}
	err = r.redisClient.Set(ctx, "video:"+strconv.FormatInt(video.ID, 10), videoJson, 0).Err()
	if err != nil {
		klog.Error("redis publish video failed,", err)
		return err
	}
	err = r.redisClient.LPush(ctx, "user_video:"+strconv.FormatInt(video.AuthorId, 10), videoJson).Err()
	if err != nil {
		klog.Error("redis publish video failed,", err)
		return err
	}
	return nil
}
func (r RedisManager) GetVideoByVideoId(ctx context.Context, vid int64) (*model.Video, error) {
	videoJson, err := r.redisClient.Get(ctx, "video:"+strconv.FormatInt(vid, 10)).Bytes()
	var video *model.Video
	err = sonic.Unmarshal(videoJson, &video)
	if err != nil {
		klog.Error("redis get video by id failed,", err)
		return nil, err
	}
	return video, nil
}
func NewRedisManager(client *redis.Client) *RedisManager {
	return &RedisManager{redisClient: client}
}
