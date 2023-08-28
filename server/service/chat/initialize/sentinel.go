package initialize

import (
	"GoYin/server/common/consts"
	"GoYin/server/service/user/config"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	SentinelConfig "github.com/alibaba/sentinel-golang/core/config"
	"github.com/cloudwego/kitex/pkg/klog"
)

func Sentinel() {
	conf := SentinelConfig.NewDefaultConfig()
	conf.Sentinel.Log.Dir = consts.ChatSentinelFilePath
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		klog.Fatal(err)
	}
	// 注册状态变化监听器，用于观察内部断路器的状态变化
	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})
	c := config.GlobalServerConfig.CbRule
	// 加载断路器规则
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// 统计时间窗口=5秒，恢复时间=3秒，慢请求上限=50毫秒，最大慢请求比例=50%
		{
			Resource:                     c.Resource,
			Strategy:                     circuitbreaker.Strategy(c.Strategy),
			RetryTimeoutMs:               c.RetryTimeoutMs,
			MinRequestAmount:             c.MinRequestAmount,
			StatIntervalMs:               c.StatIntervalMs,
			StatSlidingWindowBucketCount: c.StatSlidingWindowBucketCount,
			MaxAllowedRtMs:               c.MaxAllowedRtMs,
			Threshold:                    c.Threshold,
		},
	})
	if err != nil {
		klog.Fatal(err)
	}
}

// stateChangeTestListener 是一个用于监视断路器状态变化的监听器
type stateChangeTestListener struct {
}

// OnTransformToClosed 当从其他状态转换为关闭状态时触发
func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	klog.Info("rule.strategy: %+v, 从 %s 转换为关闭状态", rule.Strategy, prev.String())
}

// OnTransformToOpen 当从其他状态转换为打开状态时触发
func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	klog.Info("rule.strategy: %+v, 从 %s 转换为打开状态", rule.Strategy, prev.String())
}

// OnTransformToHalfOpen 当从其他状态转换为半开状态时触发
func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	klog.Info("rule.strategy: %+v, 从 %s 转换为半开状态", rule.Strategy, prev.String())
}
