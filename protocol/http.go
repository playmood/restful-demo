package protocol

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/playmood/restful-demo/apps"
	"github.com/playmood/restful-demo/conf"
	"net/http"
	"time"
)

func NewHttpService() *HttpService {
	// 还未加载路由
	r := gin.Default()
	server := &http.Server{
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1M
		Addr:              conf.C().App.HTTPAddr(),
		Handler:           r,
	}
	return &HttpService{
		server: server,
		l:      zap.L().Named("HTTP Service"),
		r:      r,
	}
}

type HttpService struct {
	server *http.Server
	l      logger.Logger
	r      gin.IRouter
}

func (s *HttpService) Start() error {
	// 加载路由，把所有模块的Handler注册给了Gin
	apps.InitGin(s.r)

	// 已加载App的日志信息
	apps := apps.LoadGinApps()
	s.l.Infof("load gin apps: %v", apps)

	// 该操作是阻塞的 监听端口 等待请求
	// 如果服务的正常关闭
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Info("service stopped success")
			return nil
		}
		return fmt.Errorf("start service error, %s", err.Error())
	}

	return nil

}

func (s *HttpService) Stop() {
	s.l.Info("start graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Errorf("graceful shutdown timeout, force exit")
	}
}
