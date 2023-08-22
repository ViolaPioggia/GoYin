package main

import (
	"GoYin/server/common/consts"
	"GoYin/server/common/middleware"
	"GoYin/server/common/tools"
	"GoYin/server/kitex_gen/base"
	user "GoYin/server/kitex_gen/user"
	"GoYin/server/service/api/models"
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
type SocialManager interface {
	GetRelationList(ctx context.Context, viewerId, ownerId int64, option int8) ([]int64, error)
	GetSocialInfo(ctx context.Context, viewerId, ownerId int64) (*base.SocialInfo, error)
	BatchGetSocialInfo(ctx context.Context, viewerId int64, ownerIdList []int64) ([]*base.SocialInfo, error)
}
type InteractionManager interface {
	GetInteractInfo(ctx context.Context, userId int64) (*base.UserInteractInfo, error)
	BatchGetInteractInfo(ctx context.Context, userIdList []int64) ([]*base.UserInteractInfo, error)
}
type ChatManager interface {
	BatchGetLatestMessage(ctx context.Context, userId int64, toUserIdList []int64) ([]*base.LatestMsg, error)
}

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	Jwt *middleware.JWT
	MysqlManager
	RedisManager
	SocialManager
	InteractionManager
	ChatManager
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
	if err != nil {
		if err.Error() == consts.MysqlAlreadyExists {
			resp.BaseResp = &base.DouyinBaseResponse{
				StatusCode: 500,
				StatusMsg:  "user already exists",
			}
			return resp, err
		} else {
			klog.Errorf("mysql create user failed: %s", err.Error())
			resp.BaseResp = &base.DouyinBaseResponse{
				StatusCode: 500,
				StatusMsg:  fmt.Sprintf("mysql create user failed: %s", err.Error()),
			}
			return resp, err
		}
	}
	err = s.RedisManager.CreateUser(ctx, usr)
	if err != nil {
		klog.Errorf("redis create user failed: %s", err.Error())
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  fmt.Sprintf("mysql create user failed: %s", err.Error()),
		}
		return resp, err
	}
	resp.UserId = usr.ID
	resp.Token, err = s.Jwt.CreateToken(models.CustomClaims{
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
		StatusCode: 0,
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
			return resp, err
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
	resp.Token, err = s.Jwt.CreateToken(models.CustomClaims{
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
		StatusCode: 0,
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

	socialInfo, err := s.SocialManager.GetSocialInfo(ctx, req.ViewerId, req.OwnerId)
	if err != nil {
		klog.Errorf("socialManager get socialInfo failed,", err)
	}

	interactionInfo, err := s.InteractionManager.GetInteractInfo(ctx, req.OwnerId)
	if err != nil {
		klog.Errorf("interactionManager get interactionInfo failed,", err)
	}

	if err != nil {
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "get userInfo failed",
		}
		return resp, err
	}

	resp.BaseResp = &base.DouyinBaseResponse{
		StatusCode: 0,
		StatusMsg:  "get user by id success",
	}
	resp.User = &base.User{
		Id:              usr.ID,
		Name:            usr.Username,
		FollowCount:     socialInfo.FollowCount,
		FollowerCount:   socialInfo.FollowerCount,
		IsFollow:        socialInfo.IsFollow,
		Avatar:          usr.Avatar,
		BackgroundImage: usr.BackGroundImage,
		Signature:       usr.Signature,
		TotalFavorited:  interactionInfo.TotalFavorited,
		WorkCount:       interactionInfo.WorkCount,
		FavoriteCount:   interactionInfo.FavoriteCount,
	}
	return resp, nil
}

// BatchGetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) BatchGetUserInfo(ctx context.Context, req *user.DouyinBatchGetUserRequest) (resp *user.DouyinBatchGetUserResonse, err error) {
	resp = new(user.DouyinBatchGetUserResonse)

	userList, err := s.RedisManager.BatchGetUserById(ctx, req.OwnerIdList)
	if err != nil {
		klog.Errorf("redis batch get users by id failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  fmt.Sprintf("redis batch get users by id failed,%s", err),
		}
		return nil, err
	}
	socialList, err := s.SocialManager.BatchGetSocialInfo(ctx, req.ViewerId, req.OwnerIdList)
	if err != nil {
		klog.Errorf("user socialManager get socialList failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user socialManager get socialList failed",
		}
		return resp, err
	}
	interactionList, err := s.InteractionManager.BatchGetInteractInfo(ctx, req.OwnerIdList)
	if err != nil {
		klog.Errorf("user interactionManager get interactionList failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user interactionManager get interactionList failed",
		}
		return resp, err
	}
	for i := 0; i <= len(userList)-1; i++ {
		resp.UserList = append(resp.UserList, &base.User{
			Id:              userList[i].ID,
			Name:            userList[i].Username,
			FollowCount:     socialList[i].FollowCount,
			FollowerCount:   socialList[i].FollowerCount,
			IsFollow:        socialList[i].IsFollow,
			Avatar:          userList[i].Avatar,
			BackgroundImage: userList[i].BackGroundImage,
			Signature:       userList[i].Signature,
			TotalFavorited:  interactionList[i].TotalFavorited,
			WorkCount:       interactionList[i].WorkCount,
			FavoriteCount:   interactionList[i].FavoriteCount,
		})
	}
	resp.BaseResp = &base.DouyinBaseResponse{
		StatusCode: 0,
		StatusMsg:  "batch get userInfo success",
	}
	return resp, nil
}

// GetFollowList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowList(ctx context.Context, req *user.DouyinGetRelationFollowListRequest) (resp *user.DouyinGetRelationFollowListResponse, err error) {
	resp = new(user.DouyinGetRelationFollowListResponse)

	userIdlist, err := s.SocialManager.GetRelationList(ctx, req.ViewerId, req.OwnerId, consts.FollowList)
	if err != nil {
		klog.Errorf("user socialManager get follow list failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user socialManager get followList failed ",
		}
		return resp, err
	}
	userList, err := s.RedisManager.BatchGetUserById(ctx, userIdlist)
	if err != nil {
		klog.Errorf("user redis get user failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user redis get user failed ",
		}
		return resp, err
	}
	socialList, err := s.SocialManager.BatchGetSocialInfo(ctx, req.ViewerId, userIdlist)
	if err != nil {
		klog.Errorf("user socialManager get socialList failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user socialManager get socialList failed",
		}
		return resp, err
	}
	interactionList, err := s.InteractionManager.BatchGetInteractInfo(ctx, userIdlist)
	if err != nil {
		klog.Errorf("user interactionManager get interactionList failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user interactionManager get interactionList failed",
		}
		return resp, err
	}
	for i := 0; i <= len(userList)-1; i++ {
		resp.UserList = append(resp.UserList, &base.User{
			Id:              userList[i].ID,
			Name:            userList[i].Username,
			FollowCount:     socialList[i].FollowCount,
			FollowerCount:   socialList[i].FollowerCount,
			IsFollow:        socialList[i].IsFollow,
			Avatar:          userList[i].Avatar,
			BackgroundImage: userList[i].BackGroundImage,
			Signature:       userList[i].Signature,
			TotalFavorited:  interactionList[i].TotalFavorited,
			WorkCount:       interactionList[i].WorkCount,
			FavoriteCount:   interactionList[i].FavoriteCount,
		})
	}
	resp.BaseResp = &base.DouyinBaseResponse{
		StatusCode: 0,
		StatusMsg:  "batch get followList success",
	}
	return resp, nil
}

// GetFollowerList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowerList(ctx context.Context, req *user.DouyinGetRelationFollowerListRequest) (resp *user.DouyinGetRelationFollowerListResponse, err error) {
	resp = new(user.DouyinGetRelationFollowerListResponse)

	userIdlist, err := s.SocialManager.GetRelationList(ctx, req.ViewerId, req.OwnerId, consts.FollowerList)
	if err != nil {
		klog.Errorf("user socialManager get follower list failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user socialManager get followerList failed ",
		}
		return resp, err
	}
	userList, err := s.RedisManager.BatchGetUserById(ctx, userIdlist)
	if err != nil {
		klog.Errorf("user redis get user failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user redis get user failed ",
		}
		return resp, err
	}
	socialList, err := s.SocialManager.BatchGetSocialInfo(ctx, req.ViewerId, userIdlist)
	if err != nil {
		klog.Errorf("user socialManager get socialList failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user socialManager get socialList failed",
		}
		return resp, err
	}
	interactionList, err := s.InteractionManager.BatchGetInteractInfo(ctx, userIdlist)
	if err != nil {
		klog.Errorf("user interactionManager get interactionList failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user interactionManager get interactionList failed",
		}
		return resp, err
	}
	for i := 0; i <= len(userList)-1; i++ {
		resp.UserList = append(resp.UserList, &base.User{
			Id:              userList[i].ID,
			Name:            userList[i].Username,
			FollowCount:     socialList[i].FollowCount,
			FollowerCount:   socialList[i].FollowerCount,
			IsFollow:        socialList[i].IsFollow,
			Avatar:          userList[i].Avatar,
			BackgroundImage: userList[i].BackGroundImage,
			Signature:       userList[i].Signature,
			TotalFavorited:  interactionList[i].TotalFavorited,
			WorkCount:       interactionList[i].WorkCount,
			FavoriteCount:   interactionList[i].FavoriteCount,
		})
	}
	resp.BaseResp = &base.DouyinBaseResponse{
		StatusCode: 0,
		StatusMsg:  "batch get followList success",
	}
	return resp, nil

}

// GetFriendList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFriendList(ctx context.Context, req *user.DouyinGetRelationFriendListRequest) (resp *user.DouyinGetRelationFriendListResponse, err error) {
	resp = new(user.DouyinGetRelationFriendListResponse)

	userIdlist, err := s.SocialManager.GetRelationList(ctx, req.ViewerId, req.OwnerId, consts.FriendsList)
	if err != nil {
		klog.Errorf("user socialManager get follow list failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user socialManager get followList failed ",
		}
		return resp, err
	}
	userList, err := s.RedisManager.BatchGetUserById(ctx, userIdlist)
	if err != nil {
		klog.Errorf("user redis get user failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user redis get user failed ",
		}
		return resp, err
	}
	socialList, err := s.SocialManager.BatchGetSocialInfo(ctx, req.ViewerId, userIdlist)
	if err != nil {
		klog.Errorf("user socialManager get socialList failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user socialManager get socialList failed",
		}
		return resp, err
	}
	interactionList, err := s.InteractionManager.BatchGetInteractInfo(ctx, userIdlist)
	if err != nil {
		klog.Errorf("user interactionManager get interactionList failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user interactionManager get interactionList failed",
		}
		return resp, err
	}
	chatList, err := s.ChatManager.BatchGetLatestMessage(ctx, req.ViewerId, userIdlist)
	if err != nil {
		klog.Errorf("user chatManager get chatList failed,", err)
		resp.BaseResp = &base.DouyinBaseResponse{
			StatusCode: 500,
			StatusMsg:  "user chatManager get chatList failed",
		}
		return resp, err
	}
	for i := 0; i <= len(userList)-1; i++ {
		resp.UserList = append(resp.UserList, &base.FriendUser{
			Id:              userList[i].ID,
			Name:            userList[i].Username,
			FollowCount:     socialList[i].FollowCount,
			FollowerCount:   socialList[i].FollowerCount,
			IsFollow:        socialList[i].IsFollow,
			Avatar:          userList[i].Avatar,
			BackgroundImage: userList[i].BackGroundImage,
			Signature:       userList[i].Signature,
			TotalFavorited:  interactionList[i].TotalFavorited,
			WorkCount:       interactionList[i].WorkCount,
			FavoriteCount:   interactionList[i].FavoriteCount,
			MsgType:         chatList[i].MsgType,
			Message:         chatList[i].Message,
		})
	}
	resp.BaseResp = &base.DouyinBaseResponse{
		StatusCode: 0,
		StatusMsg:  "batch get followList success",
	}
	return resp, nil
}
