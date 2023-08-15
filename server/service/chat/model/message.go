package model

type Message struct {
	ID         int64  `gorm:"primarykey"`
	FromUserId int64  `gorm:"not null;index:idx_from_to"`
	ToUserId   int64  `gorm:"not null;index:idx_from_to"`
	Content    string `gorm:"type:varchar(256);not null"`
	CreateTime int64  `gorm:"not null"`
}
