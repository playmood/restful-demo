package cmd

import (
	"fmt"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/playmood/restful-demo/apps"
	"github.com/playmood/restful-demo/conf"
	"github.com/playmood/restful-demo/protocol"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"

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

		// 初始化全局对象loggger
		if err := loadGlobalLogger(); err != nil {
			return err
		}

		// 加载host service实体类
		// service := impl.NewHostServiceImpl()
		// 注册到IOC
		// 采用 import _ ....impl完成注册
		// apps.HostService = impl.NewHostServiceImpl()
		// 通过Host Api对外提供HTTP Restful接口

		// apps的接口没有保存初始化Config的方法
		apps.InitImpl()
		// 提供一个GIN Router
		//engine := gin.Default()
		// 注册Ioc的所有http handler
		//apps.InitGin(engine)
		//engine.Run(conf.C().App.HTTPAddr())
		svc := NewManager()

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGINT)
		go svc.WaitStop(ch)
		return svc.Start()
	},
}

func NewManager() *manager {
	return &manager{
		http: protocol.NewHttpService(),
		l:    zap.L().Named("CLI"),
	}
}

// 管理所有需要启动的服务
// 1. HTTP服务的启动
// 2.
type manager struct {
	http *protocol.HttpService
	l    logger.Logger
}

func (m *manager) Start() error {
	return m.http.Start()
}

// 处理来自外部的中断信号，比如ternimal
func (m *manager) WaitStop(ch <-chan os.Signal) {
	for v := range ch {
		switch v {
		default:
			m.l.Infof("received signal %s", v)
			m.http.Stop()
		}
	}
}

// log 为全局变量, 只需要load 即可全局可用户, 依赖全局配置先初始化
func loadGlobalLogger() error {
	var (
		logInitMsg string
		level      zap.Level
	)

	// 根据Config里面的设置来配置全局Logger对象
	lc := conf.C().Log
	// 设置日志级别
	lv, err := zap.NewLevel(lc.Level)
	if err != nil {
		logInitMsg = fmt.Sprintf("%s, use default level INFO", err)
		level = zap.InfoLevel
	} else {
		level = lv
		logInitMsg = fmt.Sprintf("log level: %s", lv)
	}

	// 使用默认配置初始化logger的全局配置
	zapConfig := zap.DefaultConfig()
	// 配置日志的level级别
	zapConfig.Level = level
	// 程序每启动一次，不必都生成一个新的文件
	zapConfig.Files.RotateOnStartup = false

	// 配置日志输出方式
	switch lc.To {
	case conf.ToStdout:
		// 打印日志到标准输出
		zapConfig.ToStderr = true
		zapConfig.ToFiles = false
	case conf.ToFile:
		zapConfig.Files.Name = "api.log"
		zapConfig.Files.Path = lc.PathDir
	}
	// 配置日志输出格式
	switch lc.Format {
	case conf.JSONFormat:
		zapConfig.JSON = true
	}
	// 把配置应用到全局Logger
	if err := zap.Configure(zapConfig); err != nil {
		return err
	}
	zap.L().Named("INIT").Info(logInitMsg)
	return nil
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&confFile, "config", "f", "etc/demo.toml", "demo-api配置文件路径")
	RootCmd.AddCommand(StartCmd)
}
