package impl_test

import (
	"context"
	"fmt"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/playmood/restful-demo/apps/host"
	"github.com/playmood/restful-demo/apps/host/impl"
	"github.com/playmood/restful-demo/conf"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	// 定义对象必须满足该接口的实现
	service host.Service
)

func TestCreate(t *testing.T) {
	should := assert.New(t)
	ins := host.NewHost()
	ins.Name = "test"
	ins.Id = "1"
	ins.Region = "cn-hangzhou"
	ins.Type = "1"
	ins.CPU = 1
	ins.Memory = 4096
	ins, err := service.CreateHost(context.Background(), ins)
	if should.NoError(err) {
		fmt.Println(ins)
	}
}

func TestQuery(t *testing.T) {
	should := assert.New(t)
	set, err := service.QueryHost(context.Background(), host.NewQueryHostRequest())
	if should.NoError(err) {
		for i := range set.Items {
			fmt.Print(set.Items[i].Id)
		}
	}
}

func TestDescribe(t *testing.T) {
	should := assert.New(t)
	ins, err := service.DescribeHost(context.Background(), host.NewDescribeHostRequestWithId("ins-09"))
	if should.NoError(err) {
		fmt.Println(ins.Id)
	}
}

//{
//"id": "ins-10",
//"vendor": 0,
//"region": "cn-beijing",
//"create_at": 1677141684103,
//"expire_at": 0,
//"type": "sm3",
//"name": "hahahaha",
//"description": "",
//"status": "",
//"tags": null,
//"update_at": 0,
//"sync_at": 0,
//"accout": "",
//"public_ip": "10.2.1.3",
//"private_ip": "",
//"cpu": 2,
//"memory": 4096,
//"gpu_amount": 0,
//"gpu_spec": "",
//"os_type": "",
//"os_name": "",
//"serial_number": ""
//}

func TestPatch(t *testing.T) {
	should := assert.New(t)
	req := host.NewPatchUpdateHostRequest("ins-10")
	req.Name = "ddtest"
	ins, err := service.UpdateHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(ins.Id)
	}
}

func TestPut(t *testing.T) {
	should := assert.New(t)
	req := host.NewPutUpdateHostRequest("ins-06")
	req.Name = "can can need"
	req.Region = "cn-wuhan"
	req.Type = "small"
	req.CPU = 1
	req.Memory = 16384
	ins, err := service.UpdateHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(ins.Id)
	}
}

func TestDelete(t *testing.T) {
	should := assert.New(t)
	req := host.NewDeleteHostRequest("ins-02")
	ins, err := service.DeleteHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(ins.Id)
	}
}

func init() {
	err := conf.LoadConfigFromToml("../../../etc/demo.toml")
	if err != nil {
		fmt.Println(err)
	}
	// 需要初始化全局Logger
	// 为什么不设计为默认打印，因为性能
	zap.DevelopmentSetup()
	// host service的具体实现
	service = impl.NewHostServiceImpl()

}
