package main

import (
	"GoYin/server/kitex_gen/base"
	sociality "GoYin/server/kitex_gen/sociality"
	"GoYin/server/service/sociality/model"
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
)

type RedisManager interface {
	Action(ctx context.Context, userId, toUserId int64, actionType int8) error
	GetUserIdList(ctx context.Context, userId int64, option int8) ([]int64, error)
	GetSocialInfo(ctx context.Context, userId int64, viewerId int64) (*model.SocialInfo, error)
	BatchGetSocialInfo(ctx context.Context, userId []int64, viewerId int64) ([]*model.SocialInfo, error)
}

type MysqlManager interface {
	HandleSocialInfo(ctx context.Context, userId int64, toUserId int64, actionType int8) error
	GetUserIdList(ctx context.Context, userId int64, option int8) ([]int64, error)
	GetSocialInfo(ctx context.Context, userId int64, viewerId int64) (*model.SocialInfo, error)
	BatchGetSocialInfo(ctx context.Context, userId []int64, viewerId int64) ([]*model.SocialInfo, error)
}
type Publisher interface {
	Publish(ctx context.Context, req *sociality.DouyinRelationActionRequest) error
}

// SocialityServiceImpl implements the last service interface defined in the IDL.
type SocialityServiceImpl struct {
	Publisher
	RedisManager
	MysqlManager
}

// Action implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) Action(ctx context.Context, req *sociality.DouyinRelationActionRequest) (resp *sociality.DouyinRelationActionResponse, err error) {
	resp = new(sociality.DouyinRelationActionResponse)
	if req.UserId == req.ToUserId {
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "can not follow yourself",
		}
		return resp, nil
	}
	err = s.Publisher.Publish(ctx, req)
	if err != nil {
		klog.Errorf("sociality publish action failed", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "sociality publish action failed",
		}
		return resp, err
	}
	err = s.RedisManager.Action(ctx, req.UserId, req.ToUserId, req.ActionType)
	if err != nil {
		klog.Errorf("sociality redis action failed", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "sociality redis action failed",
		}
		return resp, err
	}
	resp.BaseResp = &base.DouyinBaseResponse{}
	resp.BaseResp.StatusCode = 0
	resp.BaseResp.StatusMsg = "sociality action success"
	return resp, nil
}

// GetRelationIdList implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) GetRelationIdList(ctx context.Context, req *sociality.DouyinGetRelationIdListRequest) (resp *sociality.DouyinGetRelationIdListResponse, err error) {
	resp = new(sociality.DouyinGetRelationIdListResponse)

	resp.UserIdList, err = s.RedisManager.GetUserIdList(ctx, req.OwnerId, req.Option)
	if err != nil {
		klog.Errorf("sociality redis get relationIdList failed,", err)
		resp.UserIdList, err = s.MysqlManager.GetUserIdList(ctx, req.OwnerId, req.Option)
		if err != nil {
			klog.Errorf("sociality mysql get relationIdList failed,", err)
			resp.BaseResp = &base.DouyinBaseResponse{
				StatusCode: 500,
				StatusMsg:  "sociality get relationIdList failed",
			}
			return resp, err
		}
	}
	resp.BaseResp = &base.DouyinBaseResponse{
		StatusCode: 0,
		StatusMsg:  "get sociality relationIdList success",
	}
	return resp, nil
}

// GetSocialInfo implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) GetSocialInfo(ctx context.Context, req *sociality.DouyinGetSocialInfoRequest) (resp *sociality.DouyinGetSocialInfoResponse, err error) {
	resp = new(sociality.DouyinGetSocialInfoResponse)

	socialInfo, err := s.RedisManager.GetSocialInfo(ctx, req.OwnerId, req.ViewerId)
	if err != nil {
		klog.Errorf("sociality redis get socialInfo failed,", err)
		socialInfo, err = s.MysqlManager.GetSocialInfo(ctx, req.OwnerId, req.ViewerId)
		if err != nil {
			klog.Errorf("sociality mysql get socialInfo failed,", err)
			resp.BaseResp = &base.DouyinBaseResponse{
				StatusCode: 500,
				StatusMsg:  "sociality get socialInfo failed",
			}
			return resp, err
		}
	}
	resp.BaseResp = &base.DouyinBaseResponse{
		StatusCode: 0,
		StatusMsg:  "sociality get socialInfo success",
	}
	resp.SocialInfo = &base.SocialInfo{
		FollowCount:   socialInfo.FollowerCount,
		FollowerCount: socialInfo.FollowerCount,
		IsFollow:      socialInfo.IsFollow,
	}
	return resp, nil
}

// BatchGetSocialInfo implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) BatchGetSocialInfo(ctx context.Context, req *sociality.DouyinBatchGetSocialInfoRequest) (resp *sociality.DouyinBatchGetSocialInfoResponse, err error) {
	resp = new(sociality.DouyinBatchGetSocialInfoResponse)

	socialInfos, err := s.RedisManager.BatchGetSocialInfo(ctx, req.OwnerIdList, req.ViewerId)
	if err != nil {
		klog.Errorf("sociality redis batch get socialInfo failed,", err)
		socialInfos, err = s.MysqlManager.BatchGetSocialInfo(ctx, req.OwnerIdList, req.ViewerId)
		if err != nil {
			klog.Errorf("sociality mysql batch get socialInfo failed,", err)
			resp.BaseResp = &base.DouyinBaseResponse{
				StatusCode: 500,
				StatusMsg:  "sociality batch get socialInfo failed",
			}
			return resp, err
		}
	}
	resp.BaseResp = &base.DouyinBaseResponse{
		StatusCode: 0,
		StatusMsg:  "sociality batch get socialInfo success",
	}
	var Infos []*base.SocialInfo
	for _, v := range socialInfos {
		Infos = append(Infos, &base.SocialInfo{
			FollowCount:   v.FollowCount,
			FollowerCount: v.FollowerCount,
			IsFollow:      v.IsFollow,
		})
	}
	resp.SocialInfoList = Infos
	resp.BaseResp = &base.DouyinBaseResponse{
		StatusCode: 0,
		StatusMsg:  "sociality batch get socialInfo success",
	}
	return resp, nil
}
