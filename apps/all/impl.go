package all

// 所有模块的注册 IOC层
import (
	_ "github.com/playmood/restful-demo/apps/book/api"
	_ "github.com/playmood/restful-demo/apps/book/impl"
	_ "github.com/playmood/restful-demo/apps/host/http"
	_ "github.com/playmood/restful-demo/apps/host/impl"
)
