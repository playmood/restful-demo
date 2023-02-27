package host

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"time"
)

var (
	validate = validator.New()
)

func NewHostSet() *HostSet {
	return &HostSet{
		Items: []*Host{},
	}
}

type HostSet struct {
	Total int     `json:"total"`
	Items []*Host `json:"items"`
}

func (s *HostSet) Add(item *Host) {
	s.Items = append(s.Items, item)
}

func NewQueryHostFromHTTP(r *http.Request) *QueryHostRequest {
	req := NewQueryHostRequest()
	// query string
	qs := r.URL.Query()
	pss := qs.Get("page_size")
	if pss != "" {
		ps, _ := strconv.Atoi(pss)
		req.PageSize = uint64(ps)
	}
	pns := qs.Get("page_num")
	if pns != "" {
		pn, _ := strconv.Atoi(pns)
		req.PageNumber = uint64(pn)
	}
	req.KeyWords = qs.Get("kws")
	return req
}

func NewQueryHostRequest() *QueryHostRequest {
	return &QueryHostRequest{
		PageSize:   20,
		PageNumber: 1,
		KeyWords:   "aaaaa",
	}
}

type QueryHostRequest struct {
	PageSize   uint64 `json:"page_size,omitempty"`
	PageNumber uint64 `json:"page_number,omitempty"`
	KeyWords   string `json:"kws"`
}

func (req *QueryHostRequest) GetPageSize() uint {
	return uint(req.PageSize)
}

func (req *QueryHostRequest) GetOffset() int64 {
	return int64((req.PageNumber - 1) * req.PageSize)
}

func NewDescribeHostRequestWithId(id string) *DescribeHostRequest {
	return &DescribeHostRequest{
		Id: id,
	}
}

type DescribeHostRequest struct {
	Id string
}

func NewPutUpdateHostRequest(id string) *UpdateHostRequest {
	h := NewHost()
	h.Id = id
	return &UpdateHostRequest{
		Mode: PUT,
		Host: h,
	}
}

func NewPatchUpdateHostRequest(id string) *UpdateHostRequest {
	h := NewHost()
	h.Id = id
	return &UpdateHostRequest{
		Mode: PATCH,
		Host: h,
	}
}

type UpdateHostRequest struct {
	Mode UPDATE_MODE `json:"update_mode"`
	*Host
}

func NewDeleteHostRequest(id string) *DeleteHostRequest {
	h := NewHost()
	h.Id = id
	return &DeleteHostRequest{
		Host: h,
	}
}

type DeleteHostRequest struct {
	*Host
}

func NewHost() *Host {
	return &Host{
		Resource: &Resource{},
		Describe: &Describe{},
	}
}

// Host模型的定义
type Host struct {
	// 资源公共属性部分
	*Resource
	// 资源独有属性
	*Describe
}

// 全量更新
func (h *Host) Put(ins *Host) error {
	if ins.Id != h.Id {
		return fmt.Errorf("id not euqal")
	}
	*h.Describe = *ins.Describe
	*h.Resource = *ins.Resource
	return nil
}

// 局部更新
func (h *Host) Patch(ins *Host) error {
	if ins.Name != "" {
		h.Name = ins.Name
	}
	if ins.CPU != 0 {
		h.CPU = ins.CPU
	}
	return nil
}

func (h *Host) Validate() error {
	return validate.Struct(h)
}

func (h *Host) InjectDefault() {
	if h.CreateAt == 0 {
		h.CreateAt = time.Now().UnixMilli()
	}
}

type Vendor int

const (
	// 枚举默认值
	PRIVATE_IDC Vendor = iota
	// 阿里云
	ALIYUN
	// 腾讯云
	TXYUN
	// 华为云
	HUAWEIYUN
)

type Resource struct {
	Id          string            `json:"id"  validate:"required"`     // 全局唯一Id
	Vendor      Vendor            `json:"vendor"`                      // 厂商
	Region      string            `json:"region"  validate:"required"` // 地域
	CreateAt    int64             `json:"create_at"`                   // 创建时间
	ExpireAt    int64             `json:"expire_at"`                   // 过期时间
	Type        string            `json:"type"  validate:"required"`   // 规格
	Name        string            `json:"name"  validate:"required"`   // 名称
	Description string            `json:"description"`                 // 描述
	Status      string            `json:"status"`                      // 服务商中的状态
	Tags        map[string]string `json:"tags"`                        // 标签
	UpdateAt    int64             `json:"update_at"`                   // 更新时间
	SyncAt      int64             `json:"sync_at"`                     // 同步时间
	Account     string            `json:"accout"`                      // 资源的所属账号
	PublicIP    string            `json:"public_ip"`                   // 公网IP
	PrivateIP   string            `json:"private_ip"`                  // 内网IP
}

type Describe struct {
	CPU          int    `json:"cpu" validate:"required"`    // 核数
	Memory       int    `json:"memory" validate:"required"` // 内存
	GPUAmount    int    `json:"gpu_amount"`                 // GPU数量
	GPUSpec      string `json:"gpu_spec"`                   // GPU类型
	OSType       string `json:"os_type"`                    // 操作系统类型，分为Windows和Linux
	OSName       string `json:"os_name"`                    // 操作系统名称
	SerialNumber string `json:"serial_number"`              // 序列号
}

type UPDATE_MODE string

const (
	PUT   UPDATE_MODE = "PUT"
	PATCH UPDATE_MODE = "PATCH"
)
