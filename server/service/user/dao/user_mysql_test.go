package dao

import (
	"GoYin/server/common/test"
	"GoYin/server/service/user/model"
	"context"
	"github.com/bytedance/sonic"
	"testing"
)

func TestUserLifecycleInMySQL(t *testing.T) {
	cleanUpFunc, db, err := test.RunMysqlInDocker(t)

	defer cleanUpFunc()

	if err != nil {
		t.Fatal(err)
	}

	dao := NewUser(db)

	ctx := context.Background()

	aid1 := int64(1024)
	aid2 := int64(2048)

	cases := []struct {
		name       string
		op         func() (string, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "create account1",
			op: func() (string, error) {
				err := dao.CreateUser(ctx, &model.User{
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
				err := dao.CreateUser(ctx, &model.User{
					ID:              aid1,
					Username:        "account1",
					Password:        "12345",
					Avatar:          "avatar2-url",
					BackGroundImage: "backgroundImage-url2",
					Signature:       "signature2",
				})
				return "", err
			},
			wantErr: true,
		},
		{
			name: "create account2",
			op: func() (string, error) {
				err := dao.CreateUser(ctx, &model.User{
					ID:       aid2,
					Username: "account2",
					Password: "654321",
				})
				return "", err
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "get user by username",
			op: func() (string, error) {
				user, err := dao.GetUserByUsername(ctx, "account1")
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
