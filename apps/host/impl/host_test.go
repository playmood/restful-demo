package impl_test

import (
	"context"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/playmood/restful-demo/apps/host"
	"github.com/playmood/restful-demo/apps/host/impl"
	"testing"
)

var (
	// 定义对象必须满足该接口的实现
	service host.Service
)

func TestCreate(t *testing.T) {
	ins := host.NewHost()
	ins.Name = "Test"
	service.CreateHost(context.Background(), ins)
}

func init() {
	// 需要初始化全局Logger
	// 为什么不设计为默认打印，因为性能
	zap.DevelopmentSetup()
	// host service的具体实现
	service = impl.NewHostServiceImpl()
}
