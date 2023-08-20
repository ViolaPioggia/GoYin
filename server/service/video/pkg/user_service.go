package pkg

import (
	"GoYin/server/kitex_gen/base"
	"GoYin/server/kitex_gen/user"
	"context"

	"GoYin/server/kitex_gen/user/userservice"
)

type UserManager struct {
	UserService userservice.Client
}

func NewUserManager(client userservice.Client) *UserManager {
	return &UserManager{UserService: client}
}

// BatchGetUser gets users info by list.
func (m *UserManager) BatchGetUser(ctx context.Context, list []int64, viewerId int64) ([]*base.User, error) {
	res, err := m.UserService.BatchGetUserInfo(ctx, &user.DouyinBatchGetUserRequest{
		ViewerId:    viewerId,
		OwnerIdList: list,
	})
	if err != nil {
		return nil, err
	}
	if res.BaseResp.StatusCode != 0 {
		return nil, err
	}
	return res.UserList, nil
}

// GetUser gets user info.
func (m *UserManager) GetUser(ctx context.Context, viewerId, ownerId int64) (*base.User, error) {
	resp, err := m.UserService.GetUserInfo(ctx, &user.DouyinGetUserRequest{ViewerId: viewerId, OwnerId: ownerId})
	if err != nil {
		return nil, err
	}
	if int64(resp.BaseResp.StatusCode) != 0 {
		return nil, err
	}
	return resp.User, nil
}
