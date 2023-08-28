package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/interaction/model"
	"context"
	"errors"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"time"
)

type MysqlManager struct {
	commentDb  *gorm.DB
	favoriteDb *gorm.DB
}

func (m MysqlManager) GetFavoriteCountByVideoId(videoId int64) (int64, error) {
	if videoId < 0 {
		err := errors.New("invalid user_id")
		return 0, err
	}

	var count int64
	if err := m.favoriteDb.
		Model(&model.Favorite{}).
		Select("count(*)").
		Where("video_id = ?", videoId).
		Group("video_id").
		Count(&count).Error; err != nil {

		return 0, err
	}

	return count, nil
}

func (m MysqlManager) GetFavoriteVideoCountByUserId(userId int64) (int64, error) {
	var count int64
	if err := m.favoriteDb.
		Model(&model.Favorite{}).
		Select("count(*)").
		Where("user_id = ?", userId).
		Group("user_id").
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (m MysqlManager) FavoriteAction(ctx context.Context, userId, videoId int64) error {
	if userId < 0 || videoId < 0 {
		err := errors.New("invalid user_id or video_id")
		return err
	}

	favorite := &model.Favorite{
		UserId:     userId,
		VideoId:    videoId,
		ActionType: consts.Like,
		CreateDate: time.Now().Unix(), //现在的时间戳
	}

	err := m.favoriteDb.Create(favorite).Error
	if err != nil {
		return err
	}

	return nil
}

func (m MysqlManager) UnFavoriteAction(ctx context.Context, userId, videoId int64) error {
	if userId < 0 || videoId < 0 {
		err := errors.New("invalid user_id or video_id")
		return err
	}

	var favorite model.Favorite
	err := m.favoriteDb.Where("user_id = ? AND video_id = ?", userId, videoId).First(&favorite).Error
	if err != nil {
		klog.Errorf("mysql select failed,", err)
		return err
	}
	favorite.ActionType = consts.UnLike
	err = m.favoriteDb.Model(&model.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).UpdateColumn("action_type", favorite.ActionType).Error
	if err != nil {
		klog.Errorf("mysql update failed: %v", err)
		return err
	}
	return nil
}

func (m MysqlManager) GetFavoriteVideoIdList(ctx context.Context, userId int64) ([]int64, error) {
	if userId < 0 {
		err := errors.New("invalid user_id")
		return nil, err
	}

	var favorites []model.Favorite
	err := m.favoriteDb.Where("user_id = ?", userId).Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	videoIDs := make([]int64, len(favorites))
	for i, favorite := range favorites {
		videoIDs[i] = favorite.VideoId
	}

	return videoIDs, nil
}

func (m MysqlManager) GetFavoriteCount(ctx context.Context, videoId int64) (int64, error) {
	if videoId < 0 {
		err := errors.New("invalid video_id")
		return 0, err
	}

	var count int64
	err := m.favoriteDb.
		Model(&model.Favorite{}).
		Select("count(*)").
		Where("video_id = ?", videoId).
		Group("video_id").
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m MysqlManager) JudgeIsFavoriteCount(ctx context.Context, videoId, userId int64) (bool, error) {
	if userId < 0 || videoId < 0 {
		err := errors.New("invalid user_id or video_id")
		return false, err
	}

	var favorite model.Favorite
	err := m.favoriteDb.Where("user_id = ? AND video_id = ? AND action_type = ?", userId, videoId, consts.Like).First(&favorite).Error
	if err != nil {
		return false, err
	}

	if favorite.ActionType == consts.Like {
		return true, nil
	} else {
		return false, nil
	}
}

func (m MysqlManager) Comment(ctx context.Context, comment *model.Comment) error {
	err := m.commentDb.Create(comment).Error
	if err != nil {
		return err
	}

	return nil
}

func (m MysqlManager) DeleteComment(ctx context.Context, commentId int64) error {
	if commentId < 0 {
		err := errors.New("invalid comment_id")
		return err
	}

	err := m.commentDb.Delete(&model.Comment{}, commentId).Error
	if err != nil {
		return err
	}

	return nil
}

func (m MysqlManager) GetComment(ctx context.Context, videoId int64) ([]*model.Comment, error) {
	if videoId < 0 {
		err := errors.New("invalid video_id")
		return nil, err
	}

	var comments []*model.Comment
	err := m.commentDb.Where("video_id = ?", videoId).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (m MysqlManager) GetCommentCount(ctx context.Context, videoId int64) (int64, error) {
	if videoId < 0 {
		err := errors.New("invalid video_id")
		return 0, err
	}

	var count int64
	err := m.commentDb.
		Model(&model.Comment{}).
		Select("count(*)").
		Where("video_id = ?", videoId).
		Group("video_id").
		Count(&count).Error
	if err != nil {
		klog.Errorf("mysql select failed,", err)
		return 0, err
	}

	return count, nil
}

func NewMysqlManager(db *gorm.DB) *MysqlManager {
	m := db.Migrator()
	if !m.HasTable(&model.Comment{}) {
		err := m.CreateTable(&model.Comment{})
		if err != nil {
			klog.Errorf("create mysql table failed,", err)
		}
	}
	if !m.HasTable(&model.Favorite{}) {
		err := m.CreateTable(&model.Favorite{})
		if err != nil {
			klog.Errorf("create mysql table failed,", err)
		}
	}
	return &MysqlManager{commentDb: db,
		favoriteDb: db}
}
