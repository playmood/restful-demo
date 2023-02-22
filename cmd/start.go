package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/playmood/restful-demo/apps"
	"github.com/playmood/restful-demo/apps/host/http"
	"github.com/playmood/restful-demo/conf"
	"github.com/spf13/cobra"
	// 注册所有服务实例
	_ "github.com/playmood/restful-demo/apps/all"
)

var (
	confType string
	confFile string
	confETCD string
)

// 程序的启动时 组装都在这里进行
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "启动 demo 后端API",
	Long:  "启动 demo 后端API",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 读取配置文件
		err := conf.LoadConfigFromToml(confFile)
		if err != nil {
			panic(err)
		}

		// 加载host service实体类
		// service := impl.NewHostServiceImpl()
		// 注册到IOC
		// 采用 import _ ....impl完成注册
		// apps.HostService = impl.NewHostServiceImpl()
		// 通过Host Api对外提供HTTP Restful接口

		// apps的接口没有保存初始化Config的方法
		apps.Init()
		api := http.NewHostHTTPHandler()
		// 从IOC中获取依赖
		api.Config()
		// 提供一个GIN Router
		engine := gin.Default()
		api.Registry(engine)
		return engine.Run(conf.C().App.HTTPAddr())
	},
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&confFile, "config", "f", "etc/demo.toml", "demo-api配置文件路径")
	RootCmd.AddCommand(StartCmd)
}
