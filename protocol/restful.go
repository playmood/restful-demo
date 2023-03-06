package protocol

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/playmood/restful-demo/apps"
	"github.com/playmood/restful-demo/conf"
	"net/http"
	"time"
)

func NewRestfulService() *RestfulService {
	// 还未加载路由
	r := restful.DefaultContainer
	server := &http.Server{
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1M
		Addr:              conf.C().App.RestfulAddr(),
		Handler:           r,
	}
	return &RestfulService{
		server: server,
		l:      zap.L().Named("HTTP Service"),
		r:      r,
	}
}

type RestfulService struct {
	server *http.Server
	l      logger.Logger
	r      *restful.Container
}

func (s *RestfulService) Start() error {
	// 加载路由，把所有模块的Handler注册给了Restful router
	apps.InitRestful(s.r)

	// 已加载App的日志信息
	apps := apps.LoadRestfulApps()
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

func (s *RestfulService) Stop() {
	s.l.Info("start graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Errorf("graceful shutdown timeout, force exit")
	}
}
