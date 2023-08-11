package initialize

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/chat/config"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nsqio/go-nsq"
)

func InitProducer() *nsq.Producer {
	producer, err := nsq.NewProducer(config.GlobalServerConfig.NsqInfo.Host+":"+config.GlobalServerConfig.NsqInfo.Port, nsq.NewConfig())
	if err != nil {
		klog.Error("initialize nsq producer failed,", err)
		return nil
	}
	return producer
}

func InitSubscriber() *nsq.Consumer {
	subscriber, err := nsq.NewConsumer(consts.NsqChatTopic, consts.NsqChatChannel, nsq.NewConfig())
	if err != nil {
		klog.Error("nsq initialize subscriber failed,", err)
		return nil
	}
	return subscriber
}
