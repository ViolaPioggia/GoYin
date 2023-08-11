package pkg

import (
	"GoYin/server/common/consts"
	"GoYin/server/kitex_gen/chat"
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nsqio/go-nsq"
)

type PublisherManager struct {
	Publisher *nsq.Producer
}
type SubscriberManager struct {
	Publisher *nsq.Consumer
}

func (s SubscriberManager) Subscribe(ctx context.Context) (request *chat.DouyinMessageActionRequest, err error) {
	return nil, nil
}

func (p PublisherManager) Publish(ctx context.Context, request *chat.DouyinMessageActionRequest) error {
	body, err := sonic.Marshal(request)
	if err != nil {
		klog.Error("subscriber marshal req failed,", err)
		return err
	}
	return p.Publisher.Publish(consts.NsqChatTopic, body)
}

func NewPublishManager(publisher *nsq.Producer) *PublisherManager {
	return &PublisherManager{Publisher: publisher}
}

func NewSubscriberManager(consumer *nsq.Consumer) *SubscriberManager {
	return &SubscriberManager{Publisher: consumer}
}
