package dao

import (
	"GoYin/server/common/test"
	"GoYin/server/service/chat/model"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"testing"
)

func TestChatLifecycleInMySQL(t *testing.T) {
	cleanUpFunc, db, err := test.RunMysqlInDocker(t)
	defer cleanUpFunc()

	if err != nil {
		t.Fatal(err)
	}

	dao := NewMysqlManager(db)

	ctx := context.Background()

	var msgList []*model.Message
	timeStamp := int64(1676323214)
	for i := int64(0); i < 10; i++ {
		uid1 := 1 - i%2 + 100000
		uid2 := 200001 - uid1
		msg := &model.Message{
			ID:         200000 + i,
			ToUserId:   uid2,
			FromUserId: uid1,
			Content:    fmt.Sprintf("User %d send message%d to %d", uid1, i, uid2),
			CreateTime: timeStamp + i,
		}
		msgList = append(msgList, msg)
	}

	cases := []struct {
		name       string
		op         func() (string, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "send message",
			op: func() (string, error) {
				for _, msg := range msgList {
					err = dao.HandleMessage(ctx, msg.Content, msg.FromUserId, msg.ToUserId, msg.CreateTime)
					if err != nil {
						return "", err
					}
				}
				return "", nil
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "get latest Message",
			op: func() (string, error) {
				m, err := dao.GetLatestMessage(ctx, msgList[0].FromUserId, msgList[0].ToUserId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(m)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `{"ID":200009,"ToUserId":100001,"FromUserId":100000,"Content":"User 100000 send message9 to 100001","CreateTime":1676323223}`,
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
