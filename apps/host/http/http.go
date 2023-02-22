package http

import (
	"github.com/gin-gonic/gin"
	"github.com/playmood/restful-demo/apps"
	"github.com/playmood/restful-demo/apps/host"
)

func NewHostHTTPHandler() *Handler {
	return &Handler{}
}

// 通过写一个实例类，把内部的接口通过HTTP协议暴露出去
// 基于gin的handler
type Handler struct {
	svc host.Service
}

func (h *Handler) Config() {
	if apps.HostService == nil {
		panic("dependent host service required")
	}
	// 从IOC里面获取service实例对象
	h.svc = apps.HostService
}

// 完成了HTTP Handler的注册
func (h *Handler) Registry(r gin.IRouter) {
	r.POST("/hosts", h.CreateHost)
}
