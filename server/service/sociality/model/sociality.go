package model

type SocialInfo struct {
	FollowCount   int64
	FollowerCount int64
	IsFollow      bool
}

type ConcernList struct {
	Id         int64
	UserId     int64
	FollowerId int64
}
