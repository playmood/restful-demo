package http

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/http/response"
	"github.com/playmood/restful-demo/apps/host"
)

// 用于暴露Host Service接口
func (h *Handler) CreateHost(c *gin.Context) {
	ins := host.NewHost()
	// 用户传递过来的参数进行解析，传递给Host对象实例
	if err := c.Bind(ins); err != nil {
		response.Failed(c.Writer, err)
		return
	}
	// 接口调用
	ins, err := h.svc.CreateHost(c.Request.Context(), ins)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
	response.Success(c.Writer, ins)
}

// 暴露查询接口
func (h *Handler) QueryHost(c *gin.Context) {
	req := host.NewQueryHostFromHTTP(c.Request)
	// 接口调用
	set, err := h.svc.QueryHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
	response.Success(c.Writer, set)
}
