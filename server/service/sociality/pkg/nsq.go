package pkg

import (
	"GoYin/server/common/consts"
	"GoYin/server/kitex_gen/sociality"
	"GoYin/server/service/sociality/config"
	"GoYin/server/service/sociality/dao"
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nsqio/go-nsq"
	"os"
	"os/signal"
	"syscall"
)

type PublisherManager struct {
	Publisher *nsq.Producer
}

type SubscriberManager struct {
	Subscriber *nsq.Consumer
}

func (p PublisherManager) Publish(ctx context.Context, req *sociality.DouyinRelationActionRequest) error {
	body, err := sonic.Marshal(req)
	if err != nil {
		klog.Error("subscriber marshal req failed,", err)
		return err
	}
	return p.Publisher.Publish(consts.NsqSocialityTopic, body)
}

func (s SubscriberManager) Subscribe(ctx context.Context, dao *dao.MysqlManager) (err error) {
	s.Subscriber.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		var req *sociality.DouyinRelationActionRequest
		err = sonic.Unmarshal(message.Body, &req)
		if err != nil {
			klog.Error("subscriber unmarshal socialInfo failed,", err)
			return err
		}
		err = dao.HandleSocialInfo(ctx, req.UserId, req.ToUserId, req.ActionType)
		if err != nil {
			klog.Error("mysql handleSocialInfo failed,", err)
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

func NewPublishManager(publisher *nsq.Producer) *PublisherManager {
	return &PublisherManager{Publisher: publisher}
}

func NewSubscriberManager(consumer *nsq.Consumer) *SubscriberManager {
	return &SubscriberManager{Subscriber: consumer}
}
