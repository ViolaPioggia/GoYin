package pkg

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/video/config"
	"GoYin/server/service/video/dao"
	"GoYin/server/service/video/model"
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

func (p PublisherManager) Publish(ctx context.Context, video *model.Video) error {
	body, err := sonic.Marshal(video)
	if err != nil {
		klog.Error("video publish marshal failed,", err)
		return err
	}
	err = p.Publisher.Publish(consts.NsqVideoTopic, body)
	if err != nil {
		klog.Error("video publish failed,", err)
		return err
	}
	return nil
}

func (s SubscriberManager) Subscribe(ctx context.Context, dao *dao.MysqlManager) (err error) {
	s.Subscriber.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		var req *model.Video
		err = sonic.Unmarshal(message.Body, &req)
		if err != nil {
			klog.Error("subscriber unmarshal message failed,", err)
			return err
		}
		err = dao.HandleVideo(ctx, req.ID, req.AuthorId, req.CoverUrl, req.PlayUrl, req.Title)
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
	return &PublisherManager{publisher}
}

func NewSubscriberManager(subscriber *nsq.Consumer) *SubscriberManager {
	return &SubscriberManager{subscriber}
}
