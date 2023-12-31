// Code generated by hertz generator.

package api

import (
	"GoYin/server/common/middleware"
	"GoYin/server/common/tools"
	"GoYin/server/kitex_gen/chat"
	"GoYin/server/kitex_gen/interaction"
	"GoYin/server/kitex_gen/sociality"
	"GoYin/server/kitex_gen/user"
	"GoYin/server/kitex_gen/video"
	"GoYin/server/service/api/biz/model/base"
	"GoYin/server/service/api/config"
	"GoYin/server/service/api/pkg"
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"os"
	"strings"
	"time"

	consts2 "GoYin/server/common/consts"
	api "GoYin/server/service/api/biz/model/api"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Register .
// @router /douyin/user/register/ [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinUserRegisterRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api register bindAndValidate failed", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.DouyinUserRegisterResponse)
	cg := config.GlobalUserClient
	res, err := cg.Register(ctx, &user.DouyinUserRegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		hlog.Error("rpc call user_srv failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp = &api.DouyinUserRegisterResponse{
		StatusCode: res.BaseResp.StatusCode,
		StatusMsg:  res.BaseResp.StatusMsg,
		UserID:     res.UserId,
		Token:      res.Token,
	}
	c.JSON(consts.StatusOK, resp)
}

// Login .
// @router /douyin/user/login/ [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinUserLoginRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api login bindAndValidate failed", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	cg := config.GlobalUserClient
	res, err := cg.Login(ctx, &user.DouyinUserLoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		hlog.Error("call user_srv failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	resp := new(api.DouyinUserLoginResponse)
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	resp.Token = res.Token
	resp.UserID = res.UserId
	c.JSON(consts.StatusOK, resp)
}

// GetUserInfo .
// @router /douyin/user/ [GET]
func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinUserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api getUserInfo bindAndValidate failed", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	viewerId, flag := c.Get("userId")
	if !flag {
		hlog.Error("api get viewerId failed,", err)
		c.String(consts.StatusBadRequest, errors.New("api context get viewerId failed").Error())
		return
	}
	res, err := config.GlobalUserClient.GetUserInfo(ctx, &user.DouyinGetUserRequest{
		ViewerId: viewerId.(int64),
		OwnerId:  req.UserID,
	})
	if err != nil {
		hlog.Error("rpc call user_srv failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := new(api.DouyinUserResponse)
	resp.User = &base.User{
		ID:              res.User.Id,
		Name:            res.User.Name,
		FollowCount:     res.User.FollowCount,
		FollowerCount:   res.User.FollowerCount,
		IsFollow:        res.User.IsFollow,
		Avatar:          res.User.Avatar,
		BackgroundImage: res.User.BackgroundImage,
		Signature:       res.User.Signature,
		TotalFavorited:  res.User.TotalFavorited,
		WorkCount:       res.User.WorkCount,
		FavoriteCount:   res.User.FavoriteCount,
	}
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	c.JSON(consts.StatusOK, resp)
}

// Feed .
// @router /douyin/feed/ [GET]
func Feed(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinFeedRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api feed bindAndValidate failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := new(api.DouyinFeedResponse)
	var viewerId int64
	token := req.Token
	if token != "" {
		j := middleware.NewJWT(config.GlobalServerConfig.JWTInfo.SigningKey)
		claims, err := j.ParseToken(req.Token)
		if err != nil {
			resp.StatusCode = 400
			resp.StatusMsg = "bad token"
			c.String(consts.StatusBadRequest, err.Error())
			return
		}
		viewerId = claims.ID
	}
	res, err := config.GlobalVideoClient.Feed(ctx, &video.DouyinFeedRequest{
		LatestTime: req.LatestTime,
		ViewerId:   viewerId,
	})
	if err != nil {
		hlog.Error("api call video_srv failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	resp.NextTime = res.NextTime
	for _, v := range res.VideoList {
		resp.VideoList = append(resp.VideoList, &base.Video{
			ID: v.Id,
			Author: &base.User{
				ID:              v.Author.Id,
				Name:            v.Author.Name,
				FollowCount:     v.Author.FollowCount,
				FollowerCount:   v.Author.FollowerCount,
				IsFollow:        v.Author.IsFollow,
				Avatar:          v.Author.Avatar,
				BackgroundImage: v.Author.BackgroundImage,
				Signature:       v.Author.Signature,
				TotalFavorited:  v.Author.TotalFavorited,
				WorkCount:       v.Author.WorkCount,
				FavoriteCount:   v.Author.FavoriteCount,
			},
			PlayURL:       v.PlayUrl,
			CoverURL:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
		})
	}
	c.JSON(consts.StatusOK, resp)
}

// PublishVideo .
// @router /douyin/publish/action/ [POST]
func PublishVideo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinPublishActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api publishVideo bindAndValidate failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	userId, flag := c.Get("userId")
	if !flag {
		hlog.Error("api get viewerId failed,", err)
		c.String(consts.StatusBadRequest, errors.New("api context get viewerId failed").Error())
		return
	}
	resp := new(api.DouyinPublishActionResponse)
	_, flag = c.GetQuery("data")
	if !flag {
		hlog.Info("get data success")
	} else {
		hlog.Info("get data failed")
	}
	fileHeader, err := c.Request.FormFile("data")
	if err != nil {
		hlog.Error("api read video file failed,err", err)
		resp.StatusCode = 500
		resp.StatusMsg = "get publish video formFile failed"
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	s := strings.Split(fileHeader.Filename, ".")
	s2 := s[len(s)-1:]
	suffix := strings.Join(s2, "")
	sf, err := snowflake.NewNode(consts2.MinioSnowFlakeNode)
	if err != nil {
		hlog.Error("minio snowFlake generate failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
	}
	id := sf.Generate().String()
	uploadPathBase := time.Now().Format("2006/01/02/") + id
	VTmpPath := "./tmp/video/" + id + "." + suffix
	CTmpPath := "./tmp/cover/" + id + ".png"
	VUpPath := uploadPathBase + "." + suffix
	CUpPath := uploadPathBase + ".png"
	videoFile, err := os.Create("./tmp/video/" + id + "." + suffix)
	if err != nil {
		hlog.Error("tmp create video failed")
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	defer videoFile.Close()

	mpFile, err := fileHeader.Open()
	if err != nil {
		hlog.Error("fileHeader open failed")
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	defer mpFile.Close()

	_, err = videoFile.ReadFrom(mpFile)
	if err != nil {
		hlog.Error("readFrom from mpFile failed", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	err = pkg.MinioVideoUpgrade(suffix, VTmpPath, VUpPath)
	if err != nil {
		hlog.Error("api_srv upgrade video failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	err = tools.GetVideoCover(VTmpPath, CTmpPath)
	if err != nil {
		hlog.Error("api_srv upgrade minio object failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	err = pkg.MinioCoverUpgrade(CTmpPath, CUpPath)
	if err != nil {
		hlog.Error("api_srv upgrade cover failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	res, err := config.GlobalVideoClient.PublishVideo(ctx, &video.DouyinPublishActionRequest{
		UserId:   userId.(int64),
		PlayUrl:  config.GlobalServerConfig.MinioInfo.UrlPrefix + VUpPath,
		CoverUrl: config.GlobalServerConfig.MinioInfo.UrlPrefix + CUpPath,
		Title:    req.Title,
	})
	if err != nil {
		hlog.Error("api call video_srv failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	c.JSON(consts.StatusOK, resp)
}

// VideoList .
// @router /douyin/publish/list/ [GET]
func VideoList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinPublishListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api videoList bindAndValidate failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	viewId, flag := c.Get("userId")
	if !flag {
		hlog.Error("api get viewerId failed,", err)
		c.String(consts.StatusBadRequest, errors.New("api context get viewerId failed").Error())
		return
	}
	resp := new(api.DouyinPublishListResponse)
	res, err := config.GlobalVideoClient.GetPublishedVideoList(ctx, &video.DouyinGetPublishedListRequest{
		ViewerId: viewId.(int64),
		OwnerId:  req.UserID,
	})
	if err != nil {
		hlog.Error("api get videoList failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	for _, v := range res.VideoList {
		resp.VideoList = append(resp.VideoList, &base.Video{
			ID: v.Id,
			Author: &base.User{
				ID:              v.Author.Id,
				Name:            v.Author.Name,
				FollowCount:     v.Author.FollowCount,
				FollowerCount:   v.Author.FollowerCount,
				IsFollow:        v.Author.IsFollow,
				Avatar:          v.Author.Avatar,
				BackgroundImage: v.Author.BackgroundImage,
				Signature:       v.Author.Signature,
				TotalFavorited:  v.Author.TotalFavorited,
				WorkCount:       v.Author.WorkCount,
				FavoriteCount:   v.Author.FavoriteCount,
			},
			PlayURL:       v.PlayUrl,
			CoverURL:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
		})
	}
	c.JSON(consts.StatusOK, resp)
}

// Favorite .
// @router /douyin/favorite/action/ [POST]
func Favorite(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinFavoriteActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api favorite bindAndValidate failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	userId, flag := c.Get("userId")
	if !flag {
		hlog.Error("api get viewerId failed,", err)
		c.String(consts.StatusBadRequest, errors.New("api context get viewerId failed").Error())
		return
	}
	resp := new(api.DouyinFavoriteActionResponse)
	res, err := config.GlobalInteractionClient.Favorite(ctx, &interaction.DouyinFavoriteActionRequest{
		UserId:     userId.(int64),
		VideoId:    req.VideoID,
		ActionType: req.ActionType,
	})
	if err != nil {
		hlog.Error("api call interaction rpc failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	c.JSON(consts.StatusOK, resp)
}

// FavoriteList .
// @router /douyin/favorite/list/ [GET]
func FavoriteList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinFavoriteListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api favoriteList bindAndValidate failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	viewerId, flag := c.Get("userId")
	if !flag {
		hlog.Error("api get viewerId failed,", err)
		c.String(consts.StatusBadRequest, errors.New("api context get viewerId failed").Error())
		return
	}
	res, err := config.GlobalVideoClient.GetFavoriteVideoList(ctx, &video.DouyinGetFavoriteListRequest{
		ViewerId: viewerId.(int64),
		OwnerId:  req.UserID,
	})
	if err != nil {
		hlog.Error("api call video_srv failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	resp := new(api.DouyinFavoriteListResponse)
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	for _, v := range res.VideoList {
		resp.VideoList = append(resp.VideoList, &base.Video{
			ID: v.Id,
			Author: &base.User{
				ID:              v.Author.Id,
				Name:            v.Author.Name,
				FollowCount:     v.Author.FollowCount,
				FollowerCount:   v.Author.FollowerCount,
				IsFollow:        v.Author.IsFollow,
				Avatar:          v.Author.Avatar,
				BackgroundImage: v.Author.BackgroundImage,
				Signature:       v.Author.Signature,
				TotalFavorited:  v.Author.TotalFavorited,
				WorkCount:       v.Author.WorkCount,
				FavoriteCount:   v.Author.FavoriteCount,
			},
			PlayURL:       v.PlayUrl,
			CoverURL:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
		})
	}

	c.JSON(consts.StatusOK, resp)
}

// Comment .
// @router /douyin/comment/action/ [POST]
func Comment(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinCommentActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api comment bindAndValidate failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	userId, flag := c.Get("userId")
	if !flag {
		hlog.Error("api get viewerId failed,", err)
		c.String(consts.StatusBadRequest, errors.New("api context get viewerId failed").Error())
		return
	}
	res, err := config.GlobalInteractionClient.Comment(ctx, &interaction.DouyinCommentActionRequest{
		UserId:      userId.(int64),
		VideoId:     req.VideoID,
		ActionType:  req.ActionType,
		CommentText: req.CommentText,
		CommentId:   req.CommentID,
	})
	if err != nil {
		hlog.Error("api call interaction_srv failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	resp := new(api.DouyinCommentActionResponse)
	resp.Comment = &base.Comment{
		ID: res.Comment.Id,
		User: &base.User{
			ID:              res.Comment.User.Id,
			Name:            res.Comment.User.Name,
			FollowCount:     res.Comment.User.FollowerCount,
			FollowerCount:   res.Comment.User.FollowerCount,
			IsFollow:        res.Comment.User.IsFollow,
			Avatar:          res.Comment.User.Avatar,
			BackgroundImage: res.Comment.User.BackgroundImage,
			Signature:       res.Comment.User.Signature,
			TotalFavorited:  res.Comment.User.TotalFavorited,
			WorkCount:       res.Comment.User.WorkCount,
			FavoriteCount:   res.Comment.User.FavoriteCount,
		},
		Content:    res.Comment.Content,
		CreateDate: res.Comment.CreateDate,
	}
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	c.JSON(consts.StatusOK, resp)
}

// CommentList .
// @router /douyin/comment/list/ [GET]
func CommentList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinCommentListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api commentList bindAndValidate failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	res, err := config.GlobalInteractionClient.GetCommentList(ctx, &interaction.DouyinGetCommentListRequest{VideoId: req.VideoID})
	if err != nil {
		hlog.Error("api call interaction_srv failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	resp := new(api.DouyinCommentListResponse)
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	for _, v := range res.CommentList {
		resp.CommentList = append(resp.CommentList, &base.Comment{
			ID: v.Id,
			User: &base.User{
				ID:              v.User.Id,
				Name:            v.User.Name,
				FollowCount:     v.User.FollowerCount,
				FollowerCount:   v.User.FollowerCount,
				IsFollow:        v.User.IsFollow,
				Avatar:          v.User.Avatar,
				BackgroundImage: v.User.BackgroundImage,
				Signature:       v.User.Signature,
				TotalFavorited:  v.User.TotalFavorited,
				WorkCount:       v.User.WorkCount,
				FavoriteCount:   v.User.FavoriteCount,
			},
			Content:    v.Content,
			CreateDate: v.CreateDate,
		})
	}
	c.JSON(consts.StatusOK, resp)
}

// Action .
// @router /douyin/relation/action/ [POST]
func Action(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinRelationActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api action bindAndValidate failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	userId, flag := c.Get("userId")
	if !flag {
		hlog.Error("api get viewerId failed,", err)
		c.String(consts.StatusBadRequest, errors.New("api context get viewerId failed").Error())
		return
	}
	res, err := config.GlobalSocialClient.Action(ctx, &sociality.DouyinRelationActionRequest{
		UserId:     userId.(int64),
		ToUserId:   req.ToUserID,
		ActionType: req.ActionType,
	})
	if err != nil {
		hlog.Error("api call sociality rpc failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	resp := new(api.DouyinRelationActionResponse)
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode

	c.JSON(consts.StatusOK, resp)
}

// FollowingList .
// @router /douyin/relation/follow/list/ [GET]
func FollowingList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinRelationFollowListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api followingList bindAndValidate failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	viewId, flag := c.Get("userId")
	if !flag {
		hlog.Error("api get viewerId failed,", err)
		c.String(consts.StatusBadRequest, errors.New("api context get viewerId failed").Error())
		return
	}
	res, err := config.GlobalUserClient.GetFollowList(ctx, &user.DouyinGetRelationFollowListRequest{
		ViewerId: viewId.(int64),
		OwnerId:  req.UserID,
	})
	if err != nil {
		hlog.Error("api call user_srv failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	resp := new(api.DouyinRelationFollowListResponse)
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	for _, v := range res.UserList {
		resp.UserList = append(resp.UserList, &base.User{
			ID:              v.Id,
			Name:            v.Name,
			FollowCount:     v.FollowCount,
			FollowerCount:   v.FollowerCount,
			IsFollow:        v.IsFollow,
			Avatar:          v.Avatar,
			BackgroundImage: v.BackgroundImage,
			Signature:       v.Signature,
			TotalFavorited:  v.TotalFavorited,
			WorkCount:       v.WorkCount,
			FavoriteCount:   v.FavoriteCount,
		})
	}
	c.JSON(consts.StatusOK, resp)
}

// FollowerList .
// @router /douyin/relation/follower/list/ [GET]
func FollowerList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinRelationFollowerListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api followerList bindAndValidate failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	viewId, flag := c.Get("userId")
	if !flag {
		hlog.Error("api get viewerId failed,", err)
		c.String(consts.StatusBadRequest, errors.New("api context get viewerId failed").Error())
		return
	}
	res, err := config.GlobalUserClient.GetFollowerList(ctx, &user.DouyinGetRelationFollowerListRequest{
		ViewerId: viewId.(int64),
		OwnerId:  req.UserID,
	})
	if err != nil {
		hlog.Error("api call user_srv failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	resp := new(api.DouyinRelationFollowerListResponse)
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	for _, v := range res.UserList {
		resp.UserList = append(resp.UserList, &base.User{
			ID:              v.Id,
			Name:            v.Name,
			FollowCount:     v.FollowCount,
			FollowerCount:   v.FollowerCount,
			IsFollow:        v.IsFollow,
			Avatar:          v.Avatar,
			BackgroundImage: v.BackgroundImage,
			Signature:       v.Signature,
			TotalFavorited:  v.TotalFavorited,
			WorkCount:       v.WorkCount,
			FavoriteCount:   v.FavoriteCount,
		})
	}

	c.JSON(consts.StatusOK, resp)
}

// FriendList .
// @router /douyin/relation/friend/list/ [GET]
func FriendList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinRelationFriendListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api friendList bindAndValidate failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	viewId, flag := c.Get("userId")
	if !flag {
		hlog.Error("api get viewerId failed,", err)
		c.String(consts.StatusBadRequest, errors.New("api context get viewerId failed").Error())
		return
	}
	res, err := config.GlobalUserClient.GetFriendList(ctx, &user.DouyinGetRelationFriendListRequest{
		ViewerId: viewId.(int64),
		OwnerId:  req.UserID,
	})
	if err != nil {
		hlog.Error("api call user_srv failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	resp := new(api.DouyinRelationFriendListResponse)
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	for _, v := range res.UserList {
		resp.UserList = append(resp.UserList, &base.FriendUser{
			ID:              v.Id,
			Name:            v.Name,
			FollowCount:     v.FollowCount,
			FollowerCount:   v.FollowerCount,
			IsFollow:        v.IsFollow,
			Avatar:          v.Avatar,
			BackgroundImage: v.BackgroundImage,
			Signature:       v.Signature,
			TotalFavorited:  v.TotalFavorited,
			WorkCount:       v.WorkCount,
			FavoriteCount:   v.FavoriteCount,
			Message:         v.Message,
			MsgType:         v.MsgType,
		})
	}

	c.JSON(consts.StatusOK, resp)
}

// ChatHistory .
// @router /douyin/message/chat/ [GET]
func ChatHistory(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinMessageChatRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api chatHistory bindAndValidate failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	viewId, flag := c.Get("userId")
	if !flag {
		hlog.Error("api get viewerId failed,", err)
		c.String(consts.StatusBadRequest, errors.New("api context get viewerId failed").Error())
		return
	}
	res, err := config.GlobalChatClient.GetChatHistory(ctx, &chat.DouyinMessageGetChatHistoryRequest{
		UserId:     viewId.(int64),
		ToUserId:   req.ToUserID,
		PreMsgTime: req.PreMsgTime,
	})
	if err != nil {
		hlog.Error("api call chat rpc failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	resp := new(api.DouyinMessageChatResponse)
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	for _, v := range res.MessageList {
		resp.MessageList = append(resp.MessageList, &base.Message{
			ID:         v.Id,
			ToUserID:   v.ToUserId,
			FromUserID: v.FromUserId,
			Content:    v.Content,
			CreateTime: v.CreateTime,
		})
	}
	c.JSON(consts.StatusOK, resp)
}

// SentMessage .
// @router /douyin/message/action/ [POST]
func SentMessage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinMessageActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		hlog.Error("api sentMessage failed,", err)
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	viewId, flag := c.Get("userId")
	if !flag {
		hlog.Error("api get viewerId failed,", err)
		c.String(consts.StatusBadRequest, errors.New("api context get viewerId failed").Error())
		return
	}
	res, err := config.GlobalChatClient.SentMessage(ctx, &chat.DouyinMessageActionRequest{
		UserId:     viewId.(int64),
		ToUserId:   req.ToUserID,
		ActionType: req.ActionType,
		Content:    req.Content,
	})
	if err != nil {
		hlog.Error("api call chat rpc failed,", err)
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	resp := new(api.DouyinMessageActionResponse)
	resp.StatusMsg = res.BaseResp.StatusMsg
	resp.StatusCode = res.BaseResp.StatusCode
	c.JSON(consts.StatusOK, resp)
}
