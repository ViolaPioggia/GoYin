package dao

import (
	"GoYin/server/service/interaction/model"
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type MysqlManager struct {
	commentDb  *gorm.DB
	favoriteDb *gorm.DB
}

func (m MysqlManager) FavoriteAction(ctx context.Context, userId, videoId int64) error {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) UnFavoriteAction(ctx context.Context, userId, videoId int64) error {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) GetFavoriteVideoIdList(ctx context.Context, userId int64) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) GetFavoriteCount(ctx context.Context, videoId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) JudgeIsFavoriteCount(ctx context.Context, videoId, userId int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) Comment(ctx context.Context, comment *model.Comment) error {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) DeleteComment(ctx context.Context, commentId int64) error {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) GetComment(ctx context.Context, videoId int64) ([]*model.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) GetCommentCount(ctx context.Context, videoId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
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
