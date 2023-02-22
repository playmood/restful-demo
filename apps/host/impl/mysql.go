package impl

import (
	"database/sql"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/playmood/restful-demo/apps"
	"github.com/playmood/restful-demo/apps/host"
	"github.com/playmood/restful-demo/conf"
)

// 接口实现的静态检查
// 这里如果直接写 var impl = new.... 会报错
var impl = &HostServiceImpl{}

// 保证调用该函数前全局conf对象已经初始化
func NewHostServiceImpl() *HostServiceImpl {
	return &HostServiceImpl{}
}

type HostServiceImpl struct {
	l  logger.Logger
	db *sql.DB
}

// 只需要保证 全局对象Config和全局Logger已经加载完成
func (i *HostServiceImpl) Config() {
	// Host Service服务的子logger
	// 使用封装的zap让其满足logger接口
	// 为什么要封装：
	// 1. Logger全局实例
	// 2. Logger Level的动态调整，Logrus不支持Level动态调整
	// 3. 加入日志轮转功能的集合
	i.l = zap.L().Named("Host")
	i.db = conf.C().MySQL.GetDB()
}

// 返回服务名称
func (i *HostServiceImpl) Name() string {
	return host.AppName
}

func init() {
	// 之前都是在start时候，手动把服务实现注册到IOC层
	// 注册到IOC
	// apps.HostService = impl.NewHostServiceImpl()

	// 借鉴mysql驱动加载的实现方式
	// sql这个库是一个框架，驱动是引入依赖时加载的
	// 我们把app模块比作一个驱动，ioc比作框架
	// _ import app 将app注册到IOC
	apps.Registry(impl)
}
