package impl

import (
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/playmood/restful-demo/apps/host"
)

// 接口实现的静态检查
var _ host.Service = (*HostServiceImpl)(nil)

func NewHostServiceImpl() *HostServiceImpl {
	return &HostServiceImpl{
		// Host Service服务的子logger
		// 使用封装的zap让其满足logger接口
		// 为什么要封装：
		// 1. Logger全局实例
		// 2. Logger Level的动态调整，Logrus不支持Level动态调整
		// 3. 加入日志轮转功能的集合
		l: zap.L().Named("Host"),
	}
}

type HostServiceImpl struct {
	l logger.Logger
}
