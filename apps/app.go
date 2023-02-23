package apps

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/playmood/restful-demo/apps/host"
)

// IOC 容器层： 管理所有的服务实例
// 1. HostService的实例需要注册过来，HostService才会有具体的实例，服务启动时注册
// 2. HTTP暴露模块，依赖于IOC里面的HostService
var (
	// 模块多起来，需要抽象，使用interface{} + 断言
	HostService host.Service
	// 维护当前所有的服务
	implApps = map[string]ImplService{}
	ginApps  = map[string]GinService{}
)

func RegistryImpl(svc ImplService) {
	// 服务实例注册到svcs map当中
	if _, ok := implApps[svc.Name()]; ok {
		panic(fmt.Sprintf("service %s has registried", svc.Name()))
	}
	implApps[svc.Name()] = svc
	// 根据满足的接口来实现具体的服务
	if v, ok := svc.(host.Service); ok {
		HostService = v
	}
}

// 如果指定了具体类型，就导致每增加一种类型就要多一个Get方法
// 返回空接口，使用时，由使用方进行断言
func GetImpl(name string) interface{} {
	for k, v := range implApps {
		if k == name {
			return v
		}
	}

	return nil
}

// 用于初始化 注册到IOC容器里面的所有服务
func InitImpl() {
	for _, v := range implApps {
		v.Config()
	}
}

// 已经加载的Gin有哪些
func LoadGinApps() (names []string) {
	for k := range ginApps {
		names = append(names, k)
	}

	return
}

type ImplService interface {
	Config()
	Name() string
}

// 注册由Gin编写的handler
// HTTP服务A，只需实现Registry方法
type GinService interface {
	Registry(r gin.IRouter)
	Name() string
	Config()
}

func RegistryGin(svc GinService) {
	// 服务实例注册到svcs map当中
	if _, ok := ginApps[svc.Name()]; ok {
		panic(fmt.Sprintf("service %s has registried", svc.Name()))
	}
	ginApps[svc.Name()] = svc
}

func InitGin(r gin.IRouter) {
	// 先初始化好所有对象
	for _, v := range ginApps {
		v.Config()
	}
	// 再完成handler的注册
	for _, v := range ginApps {
		v.Registry(r)
	}
}
