package host

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
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

// Host???????????????
type Host struct {
	// ????????????????????????
	*Resource
	// ??????????????????
	*Describe
}

// ????????????
func (h *Host) Put(ins *Host) error {
	if ins.Id != h.Id {
		return fmt.Errorf("id not euqal")
	}
	*h.Describe = *ins.Describe
	*h.Resource = *ins.Resource
	return nil
}

// ????????????
func (h *Host) Patch(ins *Host) error {
	//if ins.Name != "" {
	//	h.Name = ins.Name
	//}
	//if ins.CPU != 0 {
	//	h.CPU = ins.CPU
	//}
	return mergo.MergeWithOverwrite(h, ins)
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
	// ???????????????
	PRIVATE_IDC Vendor = iota
	// ?????????
	ALIYUN
	// ?????????
	TXYUN
	// ?????????
	HUAWEIYUN
)

type Resource struct {
	Id          string            `json:"id"  validate:"required"`     // ????????????Id
	Vendor      Vendor            `json:"vendor"`                      // ??????
	Region      string            `json:"region"  validate:"required"` // ??????
	CreateAt    int64             `json:"create_at"`                   // ????????????
	ExpireAt    int64             `json:"expire_at"`                   // ????????????
	Type        string            `json:"type"  validate:"required"`   // ??????
	Name        string            `json:"name"  validate:"required"`   // ??????
	Description string            `json:"description"`                 // ??????
	Status      string            `json:"status"`                      // ?????????????????????
	Tags        map[string]string `json:"tags"`                        // ??????
	UpdateAt    int64             `json:"update_at"`                   // ????????????
	SyncAt      int64             `json:"sync_at"`                     // ????????????
	Account     string            `json:"accout"`                      // ?????????????????????
	PublicIP    string            `json:"public_ip"`                   // ??????IP
	PrivateIP   string            `json:"private_ip"`                  // ??????IP
}

type Describe struct {
	CPU          int    `json:"cpu" validate:"required"`    // ??????
	Memory       int    `json:"memory" validate:"required"` // ??????
	GPUAmount    int    `json:"gpu_amount"`                 // GPU??????
	GPUSpec      string `json:"gpu_spec"`                   // GPU??????
	OSType       string `json:"os_type"`                    // ???????????????????????????Windows???Linux
	OSName       string `json:"os_name"`                    // ??????????????????
	SerialNumber string `json:"serial_number"`              // ?????????
}

type UPDATE_MODE string

const (
	PUT   UPDATE_MODE = "PUT"
	PATCH UPDATE_MODE = "PATCH"
)
