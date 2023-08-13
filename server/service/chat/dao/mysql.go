package dao

import (
	"GoYin/server/service/chat/model"
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type MysqlManager struct {
	db *gorm.DB
}

func (m MysqlManager) HandleMessage(ctx context.Context, msg string, userId, toUserId, time int64) error {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) GetHistoryMessage(ctx context.Context, userId, toUserId, time int64) ([]*model.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) GetLatestMessage(ctx context.Context, userId, toUserId int64) (*model.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (m MysqlManager) BatchGetLatestMessage(ctx context.Context, userId int64, toUserIdList []int64) ([]*model.Message, error) {
	//TODO implement me
	panic("implement me")
}

func NewMysqlManager(db *gorm.DB) *MysqlManager {
	m := db.Migrator()
	if !m.HasTable(&model.Message{}) {
		err := m.CreateTable(&model.Message{})
		if err != nil {
			klog.Errorf("create mysql table failed,", err)
		}
	}
	return &MysqlManager{db: db}
}
