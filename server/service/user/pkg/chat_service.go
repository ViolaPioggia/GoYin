package pkg

import (
	"GoYin/server/kitex_gen/base"
	"GoYin/server/kitex_gen/chat"
	"GoYin/server/kitex_gen/chat/chatservice"
	"context"
)

type ChatManager struct {
	client chatservice.Client
}

func NewChatManager(client chatservice.Client) *ChatManager {
	return &ChatManager{client: client}
}

func (m *ChatManager) BatchGetLatestMessage(ctx context.Context, userId int64, toUserIdList []int64) ([]*base.LatestMsg, error) {
	resp, err := m.client.BatchGetLatestMessage(ctx, &chat.DouyinMessageBatchGetLatestRequest{
		UserId:       userId,
		ToUserIdList: toUserIdList,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, err
	}
	return resp.LatestMsgList, nil
}
