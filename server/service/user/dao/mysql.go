package dao

import (
	"GoYin/server/service/user/model"
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func (u User) CreateUser(ctx context.Context, user *model.User) error {
	//TODO implement me
	return nil
}
func (u User) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return nil, nil
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
