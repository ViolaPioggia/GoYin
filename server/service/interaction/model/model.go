package model

type Favorite struct {
	UserId     int64 `gorm:"not null;index:idx_user_video"`
	VideoId    int64 `gorm:"not null;index:idx_user_video"`
	ActionType int8  `gorm:"type:tinyint;not null"`
	CreateDate int64 `gorm:"not null"`
}

type Comment struct {
	ID          int64  `gorm:"primarykey"`
	UserId      int64  `gorm:"not null"`
	VideoId     int64  `gorm:"not null"`
	ActionType  int8   `gorm:"type:tinyint;not null"`
	CommentText string `gorm:"type:varchar(256);not null"`
	CreateDate  int64  `gorm:"not null"`
}
