package pkg

import (
	"GoYin/server/common/consts"
	"GoYin/server/kitex_gen/chat"
	"GoYin/server/service/chat/config"
	"GoYin/server/service/chat/dao"
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nsqio/go-nsq"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type PublisherManager struct {
	Publisher *nsq.Producer
}
type SubscriberManager struct {
	Subscriber *nsq.Consumer
}

func (s SubscriberManager) Subscribe(ctx context.Context, dao *dao.MysqlManager) (err error) {
	s.Subscriber.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		var req *chat.DouyinMessageActionRequest
		err = sonic.Unmarshal(message.Body, &req)
		if err != nil {
			klog.Error("subscriber unmarshal message failed,", err)
			return err
		}
		err = dao.HandleMessage(ctx, req.Content, req.UserId, req.ToUserId, time.Now().UnixNano())
		if err != nil {
			klog.Error("subscriber handleMessage failed,", err)
			return err
		}
		return nil
	}))

	err = s.Subscriber.ConnectToNSQD(config.GlobalServerConfig.NsqInfo.Host + ":" + config.GlobalServerConfig.NsqInfo.Port)
	if err != nil {
		klog.Error(err)
	}

	// 处理退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// 停止消费者连接
	s.Subscriber.Stop()
	<-s.Subscriber.StopChan

	return nil
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
	return &SubscriberManager{Subscriber: consumer}
}
