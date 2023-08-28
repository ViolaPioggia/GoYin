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
	var temp model.User
	err := u.db.Where("username = ?", user.Username).First(&temp).Error
	if err != gorm.ErrRecordNotFound && err != nil {
		klog.Errorf("mysql select failed,", err)

		return err
	}
	if temp.Username != "" {
		err = errors.New(consts.MysqlAlreadyExists)
		return err
	}
	err = u.db.Create(&user).Error
	if err != nil {
		klog.Errorf("mysql insert failed", err)

		return err
	}

	return nil
}

func (u User) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {

		return nil, err
	}

	return &user, nil
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
