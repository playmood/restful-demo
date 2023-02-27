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

// 暴露Describe接口
func (h *Handler) DescribeHost(c *gin.Context) {
	req := host.NewDescribeHostRequestWithId(c.Param("id"))
	// 接口调用
	set, err := h.svc.DescribeHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
	response.Success(c.Writer, set)
}

// 暴露Patch接口
func (h *Handler) PatchHost(c *gin.Context) {
	req := host.NewPatchUpdateHostRequest(c.Param("id"))
	// 解析body数据
	if err := c.Bind(req.Host); err != nil {
		response.Failed(c.Writer, err)
		return
	}
	req.Id = c.Param("id")
	// 接口调用
	set, err := h.svc.UpdateHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
	response.Success(c.Writer, set)
}

// 暴露Put接口
func (h *Handler) PutHost(c *gin.Context) {
	req := host.NewPutUpdateHostRequest(c.Param("id"))
	// 解析body数据
	if err := c.Bind(req.Host); err != nil {
		response.Failed(c.Writer, err)
		return
	}
	req.Id = c.Param("id")
	// 接口调用
	set, err := h.svc.UpdateHost(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
	response.Success(c.Writer, set)
}
