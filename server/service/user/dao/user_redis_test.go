package dao

import (
	"GoYin/server/common/consts"
	"GoYin/server/common/test"
	"GoYin/server/service/user/model"
	"context"
	"github.com/bytedance/sonic"
	"testing"
)

func TestUserLifecycleInRedis(t *testing.T) {
	ctx := context.Background()

	cleanUpFunc, client, err := test.RunRedisInDocker(consts.RedisUserClientDB, t)
	defer cleanUpFunc()
	if err != nil {
		t.Fatal(err)
	}

	aid1 := int64(1024)
	aid2 := int64(2048)

	manager := NewRedisManager(client)

	cases := []struct {
		name       string
		op         func() (string, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "create account1",
			op: func() (string, error) {
				err := manager.CreateUser(ctx, &model.User{
					ID:              aid1,
					Username:        "account1",
					Password:        "12345",
					Avatar:          "avatar1-url",
					BackGroundImage: "backgroundImage-url1",
					Signature:       "signature1",
				})
				return "", err
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "duplicate create account1",
			op: func() (string, error) {
				err := manager.CreateUser(ctx, &model.User{
					ID:              aid1,
					Username:        "account1",
					Password:        "12345",
					Avatar:          "avatar1-url",
					BackGroundImage: "backgroundImage-url1",
					Signature:       "signature1",
				})
				return "", err
			},
			wantErr: true,
		},
		{
			name: "create account2",
			op: func() (string, error) {
				err := manager.CreateUser(ctx, &model.User{
					ID:              aid2,
					Username:        "account2",
					Password:        "666666",
					Avatar:          "avatar2-url",
					BackGroundImage: "backgroundImage-url2",
					Signature:       "signature2",
				})
				return "", err
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "get user by id",
			op: func() (string, error) {
				user, err := manager.GetUserById(ctx, aid1)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(user)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `{"ID":1024,"Username":"account1","Password":"12345","Avatar":"avatar1-url","BackGroundImage":"backgroundImage-url1","Signature":"signature1"}`,
		},
		{
			name: "batch get user by id",
			op: func() (string, error) {
				userList, err := manager.BatchGetUserById(ctx, []int64{aid1, aid2})
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(userList)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `[{"ID":1024,"Username":"account1","Password":"12345","Avatar":"avatar1-url","BackGroundImage":"backgroundImage-url1","Signature":"signature1"},{"ID":2048,"Username":"account2","Password":"666666","Avatar":"avatar2-url","BackGroundImage":"backgroundImage-url2","Signature":"signature2"}]`,
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
