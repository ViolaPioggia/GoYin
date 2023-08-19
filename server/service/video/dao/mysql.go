package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/video/model"
	"context"
	"errors"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"time"
)

type MysqlManager struct {
	db *gorm.DB
}

func (m MysqlManager) GetBasicVideoListByLatestTime(ctx context.Context, userId, latestTime int64) ([]*model.Video, error) {
	if userId < 0 {
		err := errors.New("invalid user_id")
		return nil, err
	}
	if latestTime < 0 {
		err := errors.New("invalid time")
		return nil, err
	}
	tx := m.db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	default:
		var videos []*model.Video
		if err := m.db.
			Order("create_time DESC"). //根据时间倒序排列视频
			Find(&videos).
			Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		return videos, nil
	}
}

func (m MysqlManager) GetPublishedVideoListByUserId(ctx context.Context, userId int64) ([]*model.Video, error) {
	if userId < 0 {
		err := errors.New("invalid user_id")
		return nil, err
	}
	tx := m.db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	default:
		var videos []*model.Video
		if err := m.db.Where("author_id = ?", userId).Find(&videos).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		return videos, nil
	}
}

func (m MysqlManager) GetFavoriteVideoListByUserId(ctx context.Context, userId int64) ([]*model.Video, error) {
	if userId < 0 {
		err := errors.New("invalid user_id")
		return nil, err
	}
	tx := m.db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	default:
		var videos []*model.Video
		if err := m.db.
			Joins("JOIN favorites ON videos.id = favorites.video_id").
			Where("favorites.user_id = ? AND favorites.action_type = ?", userId, consts.Like).
			Find(&videos).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		return videos, nil
	}
}

func (m MysqlManager) GetPublishedVideoIdListByUserId(ctx context.Context, userId int64) ([]int64, error) {
	if userId < 0 {
		err := errors.New("invalid user_id")
		return nil, err
	}
	tx := m.db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	default:
		var videos []*model.Video
		if err := m.db.Where("author_id = ?", userId).Find(&videos).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		idList := make([]int64, len(videos))
		for i, v := range videos {
			idList[i] = v.ID
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		return idList, nil
	}
}

func (m MysqlManager) PublishVideo(ctx context.Context, video *model.Video) error {
	tx := m.db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
		if err := m.db.Model(&model.Video{}).Create(&video).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}
}
func (m MysqlManager) HandleVideo(ctx context.Context, videoId, userId int64, playUrl, coverUrl, title string) error {
	tx := m.db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
		video := model.Video{
			ID:         videoId,
			AuthorId:   userId,
			PlayUrl:    playUrl,
			CoverUrl:   coverUrl,
			Title:      title,
			CreateTime: time.Now().Unix(),
		}
		if err := m.db.Model(&model.Video{}).Create(&video).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}
}

func NewMysqlManager(db *gorm.DB) *MysqlManager {
	m := db.Migrator()
	if !m.HasTable(&model.Video{}) {
		err := m.CreateTable(&model.Video{})
		if err != nil {
			klog.Errorf("create mysql table failed,", err)
		}
	}
	return &MysqlManager{db: db}
}
