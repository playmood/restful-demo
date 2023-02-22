package impl

import (
	"database/sql"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/playmood/restful-demo/apps/host"
	"github.com/playmood/restful-demo/conf"
)

// 接口实现的静态检查
var _ host.Service = (*HostServiceImpl)(nil)

// 保证调用该函数前全局conf对象已经初始化
func NewHostServiceImpl() *HostServiceImpl {
	return &HostServiceImpl{
		// Host Service服务的子logger
		// 使用封装的zap让其满足logger接口
		// 为什么要封装：
		// 1. Logger全局实例
		// 2. Logger Level的动态调整，Logrus不支持Level动态调整
		// 3. 加入日志轮转功能的集合
		l:  zap.L().Named("Host"),
		db: conf.C().MySQL.GetDB(),
	}
}

type HostServiceImpl struct {
	l  logger.Logger
	db *sql.DB
}
