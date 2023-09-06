package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/common/test"
	"context"
	"github.com/bytedance/sonic"
	"testing"
)

func TestFollowLifecycleInMySQL(t *testing.T) {
	cleanUpFunc, db, err := test.RunMysqlInDocker(t)

	defer cleanUpFunc()

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	dao := NewMysqlManager(db)

	aid1 := int64(100001)
	aid2 := int64(100002)
	aid3 := int64(100003)

	cases := []struct {
		name       string
		op         func() (string, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "create follow",
			op: func() (string, error) {
				err = dao.HandleSocialInfo(ctx, aid1, aid2, 1)
				err = dao.HandleSocialInfo(ctx, aid1, aid3, 1)
				err = dao.HandleSocialInfo(ctx, aid2, aid3, 1)
				err = dao.HandleSocialInfo(ctx, aid3, aid2, 1)
				if err != nil {
					return "", err
				}
				return "", nil
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "get follower id list",
			op: func() (string, error) {
				list, err := dao.GetUserIdList(ctx, aid2, consts.FollowerList)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(list)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: "[100003]",
		},
		{
			name: "get following id list",
			op: func() (string, error) {
				list, err := dao.GetUserIdList(ctx, aid2, consts.FollowList)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(list)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: "[100001,100003]",
		},
		{
			name: "get friend id list",
			op: func() (string, error) {
				list, err := dao.GetUserIdList(ctx, aid2, consts.FriendsList)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(list)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: "[100003]",
		},
		{
			name: "get social info",
			op: func() (string, error) {
				info, err := dao.GetSocialInfo(ctx, aid1, aid2)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(info)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `{"FollowCount":0,"FollowerCount":2,IsFollow:true}`,
		},
		{
			name: "unfollow",
			op: func() (string, error) {
				err = dao.HandleSocialInfo(ctx, aid1, aid2, 2)
				if err != nil {
					return "", err
				}
				return "", nil
			},
			wantErr:    false,
			wantResult: "",
		},
	}

	for _, cc := range cases {
		result, err := cc.op()
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s:want error;got none", cc.name)
			} else {
				continue
			}
		}
		if err != nil {
			t.Errorf("%s:operation failed: %v", cc.name, err)
		}
		if result != cc.wantResult {
			t.Errorf("%s:result err: want %s,got %s", cc.name, cc.wantResult, result)
		}
	}
}
