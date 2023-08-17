package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/sociality/model"
	"context"
	"errors"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type MysqlManager struct {
	db *gorm.DB
}

// IfExist 用来检测是否存在对应表
func (m MysqlManager) isExist(userId int64, option int8) (bool, error) {
	var temp model.ConcernList
	if option == consts.FollowList { //关注的人
		err := m.db.Where("follower_id = ?", userId).First(&temp).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return false, nil
			}
			return false, err
		}

	} else if option == consts.FollowerList { //粉丝
		err := m.db.Where("user_id = ?", userId).First(&temp).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return false, nil
			}
			return false, err
		}

	} else if option == consts.FriendsList {
		isConcern := true
		isFollow := true
		flag := true
		err1 := m.db.Where("follower_id = ?", userId).First(&temp).Error
		if err1 != nil {
			if err1 == gorm.ErrRecordNotFound {
				isConcern = false //没有关注的人
			} else {
				flag = false
			}
		}

		err2 := m.db.Where("user_id = ?", userId).First(&temp).Error
		if err2 != nil {
			if err2 == gorm.ErrRecordNotFound {
				isFollow = false //没有粉丝
			} else {
				flag = false
			}
		}
		if isConcern && isFollow && flag == true { //
			return true, nil
		}
		return false, nil
	}
	return false, errors.New("something wrong")

}

func (m MysqlManager) GetUserIdList(ctx context.Context, userId int64, option int8) ([]int64, error) {
	flag, err := m.isExist(userId, option)
	if err != nil {
		return nil, err
	}

	tx := m.db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	default:
		if !flag {
			return nil, nil
		}

		if option == consts.FollowList {
			var concernList []*model.ConcernList
			if err = m.db.Where("follower_id = ?", userId).Find(&concernList).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			idList := make([]int64, len(concernList))
			for _, v := range concernList {
				idList = append(idList, v.UserId)
			}
			if err = tx.Commit().Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			return idList, nil

		} else if option == consts.FollowerList {
			var followerList []*model.ConcernList
			if err = m.db.Where("user_id = ?", userId).Find(&followerList).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			idList := make([]int64, len(followerList))
			for _, v := range followerList {
				idList = append(idList, v.FollowerId)
			}
			if err = tx.Commit().Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			return idList, nil

		} else if option == consts.FriendsList {
			var results []*model.ConcernList
			err = m.db.Distinct().Select("user_id, follower_id"). //复杂查询，查找互关数据
				Where("user_id IN (?) AND follower_id IN (?)",
					m.db.Table("concern_lists").Select("user_id").Where("follower_id = ?", userId),
					m.db.Table("concern_lists").Select("follower_id").Where("user_id = ?", userId).
						Or("user_id = ? AND follower_id = ?", userId, userId)).
				Find(&results).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			idList := make([]int64, len(results)+1)
			for _, v := range results {
				if v.UserId == userId {
					idList = append(idList, v.FollowerId)
				}
			}
			if err = tx.Commit().Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			return idList, nil

		}
		return nil, err
	}
}

func (m MysqlManager) GetSocialInfo(ctx context.Context, userId int64, viewerId int64) (*model.SocialInfo, error) {
	concernIdList, err := m.GetUserIdList(ctx, userId, consts.FollowList)
	if err != nil {
		klog.Errorf("get IdList wrong")
		return nil, err
	}
	followerIdList, err := m.GetUserIdList(ctx, userId, consts.FollowerList)
	if err != nil {
		klog.Errorf("get IdList wrong")
		return nil, err
	}
	var flag bool
	for _, v := range followerIdList {
		if v == viewerId {
			flag = true
		}
	}
	return &model.SocialInfo{
		FollowCount:   int64(len(concernIdList)),
		FollowerCount: int64(len(followerIdList)),
		IsFollow:      flag,
	}, nil
}

func (m MysqlManager) BatchGetSocialInfo(ctx context.Context, userId []int64, viewerId int64) ([]*model.SocialInfo, error) {
	var res []*model.SocialInfo
	for _, v := range userId {
		socialInfo, err := m.GetSocialInfo(ctx, v, viewerId)
		if err != nil {
			return nil, err
		}
		res = append(res, socialInfo)
	}
	return res, nil
}

func (m MysqlManager) HandleSocialInfo(ctx context.Context, userId int64, toUserId int64, actionType int8) error {
	tx := m.db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
		var temp model.ConcernList
		err := m.db.Where("user_id = ? AND follower_id = ? ", userId, toUserId).First(&temp).Error
		if actionType == consts.Follow { //关注（创建）
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			if err != nil && err == gorm.ErrRecordNotFound {
				err = m.db.Create(&model.ConcernList{
					UserId:     userId,
					FollowerId: toUserId,
				}).Error
				if err != nil {
					tx.Rollback()
					return err
				}
			}
			return errors.New("already concern before")
		} else if actionType == consts.UnFollow { //取关（删除）
			if err != nil && err != gorm.ErrRecordNotFound {
				tx.Rollback()
				return err
			}
			if err != nil && err == gorm.ErrRecordNotFound {
				return nil
			}
			err = m.db.Where("user_id = ? AND follower_id = ?", userId, toUserId).Delete(&model.ConcernList{}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			if err = tx.Commit().Error; err != nil {
				tx.Rollback()
				return err
			}

		}
		return errors.New("invalid action_type")

	}
}

func NewMysqlManager(db *gorm.DB) *MysqlManager {
	m := db.Migrator()
	if !m.HasTable(&model.ConcernList{}) {
		err := m.CreateTable(&model.ConcernList{})
		if err != nil {
			klog.Errorf("create mysql table failed,", err)
		}
	}
	return &MysqlManager{db: db}
}
