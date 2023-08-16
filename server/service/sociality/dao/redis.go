package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/sociality/model"
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type RedisManager struct {
	redisClient *redis.Client
}

func (r RedisManager) Action(ctx context.Context, userId, toUserId int64, actionType int8) error {
	if actionType == consts.Follow {
		err := r.redisClient.SAdd(ctx, "user_follow:"+strconv.FormatInt(userId, 10), toUserId).Err()
		if err != nil {
			klog.Error("redis follow user failed,", err)
			return err
		}
		err = r.redisClient.SAdd(ctx, "user_follower:"+strconv.FormatInt(toUserId, 10), userId).Err()
		if err != nil {
			klog.Error("redis follower user failed,", err)
			return err
		}
	} else if actionType == consts.UnFollow {
		err := r.redisClient.SRem(ctx, "user_follow:"+strconv.FormatInt(userId, 10), toUserId).Err()
		if err != nil {
			klog.Error("redis unfollow user failed,", err)
			return err
		}
		err = r.redisClient.SRem(ctx, "user_follower:"+strconv.FormatInt(toUserId, 10), userId).Err()
		if err != nil {
			klog.Error("redis unfollower user failed,", err)
			return err
		}
	}
	return nil
}

func (r RedisManager) GetUserIdList(ctx context.Context, userId int64, option int8) ([]int64, error) {
	if option == consts.FollowList {
		res, err := r.redisClient.SMembers(ctx, "user_follow:"+strconv.FormatInt(userId, 10)).Result()
		if err != nil {
			klog.Error("redis get followList failed,", err)
			return nil, err
		}
		var followList []int64
		var follow int64
		for _, v := range res {
			follow, _ = strconv.ParseInt(v, 10, 64)
			followList = append(followList, follow)
		}
		return followList, nil
	} else if option == consts.FollowerList {
		res, err := r.redisClient.SMembers(ctx, "user_follower:"+strconv.FormatInt(userId, 10)).Result()
		if err != nil {
			klog.Error("redis get followerList failed,", err)
			return nil, err
		}
		var followerList []int64
		var follower int64
		for _, v := range res {
			follower, _ = strconv.ParseInt(v, 10, 64)
			followerList = append(followerList, follower)
		}
		return followerList, nil
	} else if option == consts.FriendsList {
		res, err := r.redisClient.SInter(ctx, "user_follower:"+strconv.FormatInt(userId, 10), "user_follow:"+strconv.FormatInt(userId, 10)).Result()
		if err != nil {
			klog.Error("redis sinter get friend failed,", err)
			return nil, err
		}
		var friendList []int64
		var friend int64
		for _, v := range res {
			friend, _ = strconv.ParseInt(v, 10, 64)
			friendList = append(friendList, friend)
		}
		return friendList, nil
	}
	return nil, nil
}

func (r RedisManager) GetSocialInfo(ctx context.Context, userId int64, viewerId int64) (*model.SocialInfo, error) {
	followList, err := r.GetUserIdList(ctx, userId, consts.FollowList)
	if err != nil {
		return nil, err
	}
	followerList, err := r.GetUserIdList(ctx, userId, consts.FollowerList)
	if err != nil {
		return nil, err
	}
	var flag bool
	for _, v := range followerList {
		if v == viewerId {
			flag = true
		}
	}
	return &model.SocialInfo{
		FollowCount:   int64(len(followList)),
		FollowerCount: int64(len(followerList)),
		IsFollow:      flag,
	}, nil
}

func (r RedisManager) BatchGetSocialInfo(ctx context.Context, userId []int64, viewerId int64) ([]*model.SocialInfo, error) {
	var res []*model.SocialInfo
	for _, v := range userId {
		socialInfo, err := r.GetSocialInfo(ctx, v, viewerId)
		if err != nil {
			return nil, err
		}
		res = append(res, socialInfo)
	}
	return res, nil
}

func NewRedisManager(client *redis.Client) *RedisManager {
	return &RedisManager{redisClient: client}
}
