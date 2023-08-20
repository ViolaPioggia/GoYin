package pkg

import (
	"GoYin/server/kitex_gen/base"
	"GoYin/server/kitex_gen/interaction"
	"GoYin/server/kitex_gen/interaction/interactionserver"
	"context"
	"errors"
)

type InteractionManager struct {
	client interactionserver.Client
}

func NewInteractionManager(client interactionserver.Client) *InteractionManager {
	return &InteractionManager{client: client}
}

func (i *InteractionManager) GetInteractInfo(ctx context.Context, userId int64) (*base.UserInteractInfo, error) {
	resp, err := i.client.GetUserInteractInfo(ctx, &interaction.DouyinGetUserInteractInfoRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errors.New("user use interaction rpc getInteractInfo failed")
	}
	return resp.InteractInfo, nil
}

func (i *InteractionManager) BatchGetInteractInfo(ctx context.Context, userIdList []int64) ([]*base.UserInteractInfo, error) {
	resp, err := i.client.BatchGetUserInteractInfo(ctx, &interaction.DouyinBatchGetUserInteractInfoRequest{
		UserIdList: userIdList,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errors.New("user use interaction rpc BatchGetInteractInfo failed")
	}
	return resp.InteractInfoList, nil
}
