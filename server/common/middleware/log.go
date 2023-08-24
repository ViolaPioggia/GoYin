package middleware

import (
	"GoYin/server/service/api/config"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func LogFeedInfo() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ip := c.ClientIP()
		db := config.GlobalServerConfig.IpInfo
		results, err := db.Get_all(ip)
		if err != nil {
			hlog.Error("get ip_info failed,", err)
		}
		hlog.Infof("query feed from country:%s,region:%s,city:%s", results.Country_long, results.Region, results.City)
	}
}
