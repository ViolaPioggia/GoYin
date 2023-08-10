package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/user/model"
	"context"
	"errors"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func (u User) CreateUser(ctx context.Context, user *model.User) error {
	tx := u.db.Begin() //开启事务

	if tx.Error != nil { //检查事务是否正常开启
		tx.Rollback()
		return tx.Error
	}

	select {
	case <-ctx.Done(): //若收到取消信号则退出
		tx.Rollback()
		return ctx.Err()
	default:
		var temp model.User
		err := u.db.Where("username = ?", user.Username).First(&temp).Error
		if err != nil {
			klog.Errorf("mysql select failed,", err)
			tx.Rollback()
			return err
		}
		if temp.Username != "" {
			err = errors.New(consts.MysqlAlreadyExists)
			tx.Rollback()
			return err
		}
		err = u.db.Create(&user).Error
		if err != nil {
			klog.Errorf("mysql insert failed", err)
			tx.Rollback()
			return err
		}

		if err = tx.Commit().Error; err != nil { //提交事务并判断是否成功提交
			tx.Rollback()
			return err
		}
		return nil
	}
}
func (u User) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	tx := u.db.Begin() //开启事务

	if tx.Error != nil { //检查事务是否正常开启
		tx.Rollback()
		return nil, tx.Error
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	default:
		var user model.User
		if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Commit().Error; err != nil { //提交事务并判断是否成功提交
			tx.Rollback()
			return nil, err
		}
		return &user, nil
	}
}
func NewUser(db *gorm.DB) *User {
	m := db.Migrator()
	if !m.HasTable(&model.User{}) {
		err := m.CreateTable(&model.User{})
		if err != nil {
			klog.Errorf("create mysql table failed,", err)
		}
	}
	return &User{db: db}
}
