package initialize

import (
	serverConfig "GoYin/server/service/api/config"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func InitSentinel() {
	cfg := config.NewDefaultConfig()
	cfg.Sentinel.Log.Dir = "./tmp/sentinel/api"
	err := sentinel.InitWithConfig(cfg)
	if err != nil {
		hlog.Fatal("init sentinel failed", err)
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               serverConfig.GlobalServerConfig.FlowRule.Resource,
			Threshold:              float64(serverConfig.GlobalServerConfig.FlowRule.Threshold),
			TokenCalculateStrategy: flow.TokenCalculateStrategy(serverConfig.GlobalServerConfig.FlowRule.TokenCalculateStrategy),
			ControlBehavior:        flow.ControlBehavior(serverConfig.GlobalServerConfig.FlowRule.TokenCalculateStrategy),
			StatIntervalInMs:       serverConfig.GlobalServerConfig.FlowRule.StatIntervalInMs,
		},
	})
	if err != nil {
		hlog.Fatal("load sentinel failed", err)
	}
}
