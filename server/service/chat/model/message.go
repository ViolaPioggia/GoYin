package model

type Message struct {
	ID         int64  `gorm:"primarykey"`
	ToUserId   int64  `gorm:"not null"`
	FromUserId int64  `gorm:"not null"`
	Content    string `gorm:"type:varchar(256);not null"`
	CreateTime int64  `gorm:"not null"`
}
