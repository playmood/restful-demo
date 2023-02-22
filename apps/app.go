package apps

import "github.com/playmood/restful-demo/apps/host"

// IOC 容器层： 管理所有的服务实例
// 1. HostService的实例需要注册过来，HostService才会有具体的实例，服务启动时注册
// 2. HTTP暴露模块，依赖于IOC里面的HostService
var (
	HostService host.Service
)
