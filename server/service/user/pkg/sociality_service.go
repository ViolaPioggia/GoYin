package pkg

import (
	"GoYin/server/kitex_gen/base"
	"GoYin/server/kitex_gen/sociality"
	"GoYin/server/kitex_gen/sociality/socialityservice"
	"context"
	"errors"
)

type SocialManager struct {
	client socialityservice.Client
}

func NewSocialManager(client socialityservice.Client) *SocialManager {
	return &SocialManager{client: client}
}

func (s *SocialManager) GetSocialInfo(ctx context.Context, viewerId, ownerId int64) (*base.SocialInfo, error) {
	resp, err := s.client.GetSocialInfo(ctx, &sociality.DouyinGetSocialInfoRequest{
		ViewerId: viewerId,
		OwnerId:  ownerId,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errors.New("use social rpc getSocialInfo failed")
	}
	return resp.SocialInfo, nil
}

func (s *SocialManager) BatchGetSocialInfo(ctx context.Context, viewerId int64, ownerIdList []int64) ([]*base.SocialInfo, error) {
	resp, err := s.client.BatchGetSocialInfo(ctx, &sociality.DouyinBatchGetSocialInfoRequest{
		ViewerId:    viewerId,
		OwnerIdList: ownerIdList,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errors.New("use social rpc batchGetSocialInfo failed")
	}
	return resp.SocialInfoList, nil
}

func (s *SocialManager) GetRelationList(ctx context.Context, viewerId, ownerId int64, option int8) ([]int64, error) {
	resp, err := s.client.GetRelationIdList(ctx, &sociality.DouyinGetRelationIdListRequest{
		ViewerId: viewerId,
		OwnerId:  ownerId,
		Option:   option,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errors.New("use social rpc getRelationList failed")
	}
	return resp.UserIdList, nil
}
