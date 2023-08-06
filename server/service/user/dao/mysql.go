package dao

import (
	"GoYin/server/service/user/models"
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func (u User) CreateUser(ctx context.Context, user *models.User) error {
	//TODO implement me
	return nil
}
func (u User) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return nil, nil
}
func NewUser(db *gorm.DB) *User {
	m := db.Migrator()
	if !m.HasTable(&models.User{}) {
		err := m.CreateTable(&models.User{})
		if err != nil {
			klog.Errorf("create mysql table failed,", err)
		}
	}
	return &User{db: db}
}
