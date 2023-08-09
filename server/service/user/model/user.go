package model

type User struct {
	ID              int64
	Username        string `gorm:"index:idx_username,unique"`
	Password        string
	Avatar          string
	BackGroundImage string
	Signature       string
}
