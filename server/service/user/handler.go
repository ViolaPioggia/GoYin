package main

import (
	"GoYin/server/common/consts"
	"GoYin/server/common/middleware"
	"GoYin/server/common/tools"
	"GoYin/server/kitex_gen/base"
	user "GoYin/server/kitex_gen/user"
	models2 "GoYin/server/service/api/models"
	"GoYin/server/service/user/config"
	"GoYin/server/service/user/model"
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"time"
)

type MysqlManager interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
}
type RedisManager interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserById(ctx context.Context, id int64) (*model.User, error)
	BatchGetUserById(ctx context.Context, id []int64) ([]*model.User, error)
}

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	Jwt *middleware.JWT
	MysqlManager
	RedisManager
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	resp = new(user.DouyinUserRegisterResponse)

	sf, err := snowflake.NewNode(consts.UserSnowflakeNode)
	if err != nil {
		klog.Errorf("generate user snowflake id failed: %s", err.Error())
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  err.Error(),
		}
		return resp, nil
	}
	usr := &model.User{
		ID:              sf.Generate().Int64(),
		Username:        req.Username,
		Password:        tools.Md5Crypt(req.Password, config.GlobalServerConfig.MysqlInfo.Salt),
		Avatar:          "",
		BackGroundImage: "",
		Signature:       "default signature",
	}
	err = s.MysqlManager.CreateUser(ctx, usr)
	if err.Error() == consts.MysqlAlreadyExists {
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user already exists",
		}
		return resp, err
	} else if err != nil {
		klog.Errorf("mysql create user failed: %s", err.Error())
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  fmt.Sprintf("mysql create user failed: %s", err.Error()),
		}
		return resp, err
	}
	err = s.RedisManager.CreateUser(ctx, usr)
	if err != nil {
		klog.Errorf("mysql create user failed: %s", err.Error())
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  fmt.Sprintf("mysql create user failed: %s", err.Error()),
		}
		return resp, err
	}
	resp.UserId = usr.ID
	resp.Token, err = s.Jwt.CreateToken(models2.CustomClaims{
		ID: usr.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 60*60*24*30,
			Issuer:    "GoYin",
			NotBefore: time.Now().Unix(),
		},
	})
	if err != nil {
		klog.Errorf("register create jwt failed", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  fmt.Sprintf("register create jwt failed,%s", err),
		}
		return resp, err
	}
	resp.BaseResp = &base.DouyinBaseResponse{
		StatusCode: 200,
		StatusMsg:  "user register success",
	}
	return resp, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.DouyinUserLoginRequest) (resp *user.DouyinUserLoginResponse, err error) {
	resp = new(user.DouyinUserLoginResponse)

	usr, err := s.MysqlManager.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.BaseResp = &base.DouyinBaseResponse{
				StatusCode: 500,
				StatusMsg:  "no such user",
			}
			return resp, nil
		} else {
			klog.Errorf("mysql get user by username failed", err)
			resp.BaseResp = &base.DouyinBaseResponse{
				StatusCode: 500,
				StatusMsg:  err.Error(),
			}
			return resp, err
		}
	}

	if usr.Password != tools.Md5Crypt(req.Password, config.GlobalServerConfig.MysqlInfo.Salt) {
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "wrong password",
		}
		return resp, nil
	}

	resp.UserId = usr.ID
	resp.Token, err = s.Jwt.CreateToken(models2.CustomClaims{
		ID: usr.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 60*60*24*30,
			Issuer:    "GoYin",
			NotBefore: time.Now().Unix(),
		},
	})
	if err != nil {
		klog.Errorf("register create jwt failed", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  fmt.Sprintf("register create jwt failed,%s", err),
		}
		return resp, err
	}

	resp.BaseResp = &base.DouyinBaseResponse{
		StatusCode: 200,
		StatusMsg:  "login success",
	}
	return resp, nil
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.DouyinGetUserRequest) (resp *user.DouyinGetUserResponse, err error) {
	resp = new(user.DouyinGetUserResponse)

	usr, err := s.RedisManager.GetUserById(ctx, req.OwnerId)
	if err != nil {
		klog.Errorf("redis get user by id failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  fmt.Sprintf("redis get user by id failed,%s", err),
		}
		return nil, err
	}
	resp.BaseResp = &base.DouyinBaseResponse{
		StatusCode: 200,
		StatusMsg:  "get user by id success",
	}
	resp.User = &base.User{
		Id:              usr.ID,
		Name:            usr.Username,
		FollowCount:     0,
		FollowerCount:   0,
		IsFollow:        false,
		Avatar:          usr.Avatar,
		BackgroundImage: usr.BackGroundImage,
		Signature:       usr.Signature,
		TotalFavorited:  0,
		WorkCount:       0,
		FavoriteCount:   0,
	}
	return
}

// BatchGetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) BatchGetUserInfo(ctx context.Context, req *user.DouyinBatchGetUserRequest) (resp *user.DouyinBatchGetUserResonse, err error) {
	resp = new(user.DouyinBatchGetUserResonse)

	usrs, err := s.BatchGetUserById(ctx, req.OwnerIdList)
	if err != nil {
		klog.Errorf("redis batch get users by id failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  fmt.Sprintf("redis batch get users by id failed,%s", err),
		}
		return nil, err
	}
	resp.BaseResp = &base.DouyinBaseResponse{
		StatusCode: 200,
		StatusMsg:  "batch get users by id success",
	}
	var u []*base.User
	for _, v := range usrs {
		u = append(u, &base.User{
			Id:              v.ID,
			Name:            v.Username,
			FollowCount:     0,
			FollowerCount:   0,
			IsFollow:        false,
			Avatar:          v.Avatar,
			BackgroundImage: v.BackGroundImage,
			Signature:       v.Signature,
			TotalFavorited:  0,
			WorkCount:       0,
			FavoriteCount:   0,
		})
	}
	resp.UserList = u
	return resp, nil
}

// GetFollowList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowList(ctx context.Context, req *user.DouyinGetRelationFollowListRequest) (resp *user.DouyinGetRelationFollowListResponse, err error) {
	resp = new(user.DouyinGetRelationFollowListResponse)

	return
}

// GetFollowerList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowerList(ctx context.Context, req *user.DouyinGetRelationFollowerListRequest) (resp *user.DouyinGetRelationFollowerListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFriendList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFriendList(ctx context.Context, req *user.DouyinGetRelationFriendListRequest) (resp *user.DouyinGetRelationFriendListResponse, err error) {
	// TODO: Your code here...
	return
}
