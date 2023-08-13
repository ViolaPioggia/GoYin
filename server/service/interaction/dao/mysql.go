package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/interaction/model"
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"time"
)

type MysqlManager struct {
	commentDb  *gorm.DB
	favoriteDb *gorm.DB
}

func (m MysqlManager) GetFavoriteCountByVideoId(videoId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) GetFavoriteVideoCountByUserId(userId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) FavoriteAction(ctx context.Context, userId, videoId int64) error {
	tx := m.favoriteDb.Begin()

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:

		favorite := &model.Favorite{
			UserId:     userId,
			VideoId:    videoId,
			ActionType: consts.Like,
			CreateDate: time.Now().Unix(), //现在的时间戳
		}

		err := m.favoriteDb.Create(favorite).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		if err = tx.Commit().Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}
}

func (m MysqlManager) UnFavoriteAction(ctx context.Context, userId, videoId int64) error {
	tx := m.favoriteDb.Begin()

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
		var favorite model.Favorite
		err := m.favoriteDb.Where("user_id = ? AND video_id = ?", userId, videoId).First(&favorite).Error
		if err != nil {
			klog.Errorf("mysql select failed,", err)
			tx.Rollback()
			return err
		}
		favorite.ActionType = consts.UnLike
		err = m.favoriteDb.Model(&model.Favorite{}).Updates(favorite).Error
		if err != nil {
			klog.Errorf("mysql insert failed,", err)
			tx.Rollback()
			return err
		}
		if err = tx.Commit().Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}
}

func (m MysqlManager) GetFavoriteVideoIdList(ctx context.Context, userId int64) ([]int64, error) {
	tx := m.favoriteDb.Begin()

	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	default:
		var favorites []model.Favorite
		err := m.favoriteDb.Where("user_id = ?", userId).Find(&favorites).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		videoIDs := make([]int64, len(favorites))
		for i, favorite := range favorites {
			videoIDs[i] = favorite.VideoId
		}
		if err = tx.Commit().Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		return videoIDs, nil
	}
}

func (m MysqlManager) GetFavoriteCount(ctx context.Context, videoId int64) (int64, error) {
	tx := m.favoriteDb.Begin()

	if tx.Error != nil {
		tx.Rollback()
		return 0, tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return 0, ctx.Err()
	default:
		var count int64
		err := m.favoriteDb.Model(&model.Favorite{}).Select("count(*)").Where("video_id = ?", videoId).Group("video_id").Count(&count).Error
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if err = m.favoriteDb.Commit().Error; err != nil {
			tx.Rollback()
			return 0, err
		}
		return count, nil
	}
}

func (m MysqlManager) JudgeIsFavoriteCount(ctx context.Context, videoId, userId int64) (bool, error) {
	tx := m.favoriteDb.Begin()

	if tx.Error != nil {
		tx.Rollback()
		return false, tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return false, ctx.Err()
	default:
		var favorite model.Favorite
		err := m.favoriteDb.Where("user_id = ? AND video_id = ?", userId, videoId).First(&favorite).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		if err = m.favoriteDb.Commit().Error; err != nil {
			tx.Rollback()
			return false, err
		}

		if favorite.ActionType == consts.Like {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func (m MysqlManager) Comment(ctx context.Context, comment *model.Comment) error {
	tx := m.commentDb.Begin()

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
		err := m.commentDb.Create(comment).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		if err = tx.Commit().Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	}
}

func (m MysqlManager) DeleteComment(ctx context.Context, commentId int64) error {
	tx := m.commentDb.Begin()

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
		err := m.commentDb.Delete(&model.Comment{}, commentId).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		if err = tx.Commit().Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	}
}

func (m MysqlManager) GetComment(ctx context.Context, videoId int64) ([]*model.Comment, error) {
	tx := m.commentDb.Begin()

	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	default:
		var comments []*model.Comment
		err := m.commentDb.Where("video_id = ?", videoId).Find(&comments).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if err = tx.Commit().Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		return comments, nil
	}
}

func (m MysqlManager) GetCommentCount(ctx context.Context, videoId int64) (int64, error) {
	tx := m.commentDb.Begin()

	if tx.Error != nil {
		tx.Rollback()
		return 0, tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return 0, ctx.Err()
	default:
		var count int64
		err := m.commentDb.Model(&model.Comment{}).Select("count(*)").Where("video_id = ?", videoId).Group("video_id").Count(&count).Error
		if err != nil {
			klog.Errorf("mysql select failed,", err)
			tx.Rollback()
			return 0, err
		}

		if err = tx.Commit().Error; err != nil {
			tx.Rollback()
			return 0, err
		}

		return count, nil
	}

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
