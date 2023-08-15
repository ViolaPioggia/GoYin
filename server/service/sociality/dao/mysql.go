package dao

import (
	"GoYin/server/service/sociality/model"
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type MysqlManager struct {
	db *gorm.DB
}

func (m MysqlManager) Action(ctx context.Context, userId, toUserId int64, actionType int8) error {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) GetUserIdList(ctx context.Context, userId int64, option int8) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) GetSocialInfo(ctx context.Context, userId int64) (*model.SocialInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) BatchGetSocialInfo(ctx context.Context, userId []int64) ([]*model.SocialInfo, error) {
	//TODO implement me
	panic("implement me")
}

func NewMysqlManager(db *gorm.DB) *MysqlManager {
	m := db.Migrator()
	if !m.HasTable(&model.SocialInfo{}) {
		err := m.CreateTable(&model.SocialInfo{})
		if err != nil {
			klog.Errorf("create mysql table failed,", err)
		}
	}
	return &MysqlManager{db: db}
}
