package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/common/test"
	"GoYin/server/service/interaction/model"
	"context"
	"github.com/bytedance/sonic"
	"testing"
	"time"
)

func TestFavoriteLifecycleInRedis(t *testing.T) {
	ctx := context.Background()

	cleanFunc, client, err := test.RunRedisInDocker(consts.RedisFavoriteClientDB, t)
	defer cleanFunc()

	if err != nil {
		t.Fatal(err)
	}

	manager := NewRedisManager(client)

	favList := make([]*model.Favorite, 0)
	timeStamp := int64(1676323214)
	for i := int64(0); i < 10; i++ {
		f := &model.Favorite{
			UserId:     100000 + i%3,
			VideoId:    200000 + i%5,
			ActionType: consts.IsLike,
			CreateDate: timeStamp + i,
		}
		favList = append(favList, f)

	}

	cases := []struct {
		name       string
		op         func() (string, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "create favorite",
			op: func() (string, error) {
				time.Sleep(1 * time.Second) // wait redis docker completely start
				for _, f := range favList {
					if err = manager.FavoriteAction(ctx, f.UserId, f.VideoId); err != nil {
						return "", err
					}
				}
				return "", nil
			},
		},
		{
			name: "get favorite video id list by userid",
			op: func() (string, error) {
				list, err := manager.GetFavoriteVideoIdList(ctx, favList[0].UserId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(list)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `[200004,200001,200003,200000]`,
		},
		{
			name: "cancel favorite",
			op: func() (string, error) {
				return "", manager.UnFavoriteAction(ctx, favList[0].UserId, favList[0].VideoId)
			},
			wantErr: false,
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
