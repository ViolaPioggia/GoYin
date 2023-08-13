package model

type Video struct {
	ID         int64  `gorm:"primarykey"`
	AuthorId   int64  `gorm:"column:author_id; not null"`
	PlayUrl    string `gorm:"not null; type: varchar(255)"`
	CoverUrl   string `gorm:"not null; type: varchar(255)"`
	Title      string `gorm:"not null; type: varchar(255)"`
	CreateTime int64  `gorm:"not null;"`
}
