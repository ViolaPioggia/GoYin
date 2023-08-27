package initialize

import (
	"GoYin/server/service/api/config"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/ip2location/ip2location-go/v9"
)

func InitIP() {
	ipDb, err := ip2location.OpenDB("./deployment/ip_info/IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		klog.Fatal("initialize ipDb failed,", err)
		return
	}

	config.GlobalServerConfig.IpInfo = ipDb
}
