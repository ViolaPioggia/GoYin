package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/common/test"
	"context"
	"github.com/bytedance/sonic"
	"testing"
	"time"
)

func TestFollowLifecycleRedis(t *testing.T) {
	ctx := context.Background()

	cleanUpFunc, client, err := test.RunRedisInDocker(consts.RedisSocialClientDB, t)
	defer cleanUpFunc()

	if err != nil {
		t.Fatal(err)
	}

	manager := NewRedisManager(client)

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
			name: "action",
			op: func() (string, error) {
				time.Sleep(1 * time.Second)
				err = manager.Action(ctx, aid3, aid1, consts.Follow)
				err = manager.Action(ctx, aid2, aid1, consts.Follow)
				err = manager.Action(ctx, aid3, aid2, consts.Follow)
				err = manager.Action(ctx, aid2, aid3, consts.Follow)
				if err != nil {
					return "", err
				}
				return "", nil
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "get following id list",
			op: func() (string, error) {
				list, err := manager.GetUserIdList(ctx, aid2, consts.FollowList)
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
			name: "get follower list",
			op: func() (string, error) {
				list, err := manager.GetUserIdList(ctx, aid2, consts.FollowerList)
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
			name: "get friend list",
			op: func() (string, error) {
				list, err := manager.GetUserIdList(ctx, aid2, consts.FriendsList)
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
				info, err := manager.GetSocialInfo(ctx, aid1, aid2)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(info)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `{"FollowCount":0,"FollowerCount":2,IsFollow:true}`,
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
