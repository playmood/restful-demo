package http

import (
	"github.com/gin-gonic/gin"
	"github.com/playmood/restful-demo/apps/host"
)

func NewHostHTTPHandler(svc host.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// 通过写一个实例类，把内部的接口通过HTTP协议暴露出去
// 基于gin的handler
type Handler struct {
	svc host.Service
}

// 完成了HTTP Handler的注册
func (h *Handler) Registry(r gin.IRouter) {
	r.POST("/hosts", h.CreateHost)
}
