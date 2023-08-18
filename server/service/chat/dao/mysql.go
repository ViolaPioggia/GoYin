package dao

import (
	"GoYin/server/service/chat/model"
	"context"
	"errors"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type MysqlManager struct {
	db *gorm.DB
}

func (m MysqlManager) HandleMessage(ctx context.Context, msg string, userId, toUserId, time int64) error {
	if userId < 0 || toUserId < 0 {
		err := errors.New("invalid user_id or to_user_id")
		return err
	}
	if msg == "" {
		err := errors.New("msg nil")
		return err
	}
	tx := m.db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		message := model.Message{
			ToUserId:   toUserId,
			FromUserId: userId,
			Content:    msg,
			CreateTime: time,
		}
		if err := m.db.Create(&message).Error; err != nil {
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

func (m MysqlManager) GetHistoryMessage(ctx context.Context, userId, toUserId, time int64) ([]*model.Message, error) {
	if userId < 0 || toUserId < 0 {
		err := errors.New("invalid user_id or to_user_id")
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
		var messages []*model.Message
		if err := m.db.
			Order("create_time DESC").
			Where("from_user_id = ? AND to_user_id = ? AND create_time < ?", userId, toUserId, time).
			Or("to_user_id = ? AND from_user_id = ? AND create_time < ?", toUserId, userId, time).
			Find(&messages).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		return messages, nil
	}
}

func (m MysqlManager) GetLatestMessage(ctx context.Context, userId, toUserId int64) (*model.Message, error) {
	if userId < 0 || toUserId < 0 {
		err := errors.New("invalid user_id or to_user_id")
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
		var message model.Message
		if err := m.db.
			Order("create_time DESC").
			Where("from_user_id = ? AND to_user_id = ?", userId, toUserId).
			First(&message).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		return &message, nil
	}
}

func (m MysqlManager) BatchGetLatestMessage(ctx context.Context, userId int64, toUserIdList []int64) ([]*model.Message, error) {
	if userId < 0 {
		err := errors.New("invalid user_id ")
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
		var messages []*model.Message
		for _, v := range toUserIdList {
			var msg model.Message
			if err := m.db.
				Where("from_user_id = ? AND to_user = ?", userId, v).
				Order("create_time DESC").
				First(&msg).
				Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			messages = append(messages, &msg)
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		return messages, nil
	}
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
