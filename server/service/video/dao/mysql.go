package dao

import (
	"GoYin/server/service/video/model"
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type MysqlManager struct {
	db *gorm.DB
}

func (m MysqlManager) GetBasicVideoListByLatestTime(ctx context.Context, userId, latestTime int64) ([]*model.Video, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) GetPublishedVideoListByUserId(ctx context.Context, userId int64) ([]*model.Video, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) GetFavoriteVideoListByUserId(ctx context.Context, userId int64) ([]*model.Video, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) GetPublishedVideoIdListByUserId(ctx context.Context, userId int64) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) PublishVideo(ctx context.Context, video *model.Video) error {
	//TODO implement me
	panic("implement me")
}
func (m MysqlManager) HandleVideo(ctx context.Context, userId int64, playUrl, coverUrl, title string) error {
	return nil
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
