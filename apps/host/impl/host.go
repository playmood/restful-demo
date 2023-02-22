package impl

import (
	"context"
	"github.com/playmood/restful-demo/apps/host"
)

// controller层
func (i *HostServiceImpl) CreateHost(ctx context.Context, ins *host.Host) (*host.Host, error) {
	// 直接打印日志
	//i.l.Debug("create host")
	// 带format的日志打印，fmt.Sprintf()
	//i.l.Debugf("create host %s", ins.Name)
	// 携带额外meta数据，常用于Trace系统
	//i.l.With(logger.NewAny("request-id", "req01")).Debug("create host with meta kv")
	// 校验数据合法性
	if err := ins.Validate(); err != nil {
		return nil, err
	}

	// 默认值填充
	ins.InjectDefault()

	// dao负责对象入库
	if err := i.save(ctx, ins); err != nil {
		return nil, err
	}
	return ins, nil
}

func (i *HostServiceImpl) QueryHost(ctx context.Context, request *host.QueryHostRequest) (*host.HostSet, error) {

	return nil, nil
}

func (i *HostServiceImpl) DescribeHost(ctx context.Context, request *host.QueryHostRequest) (*host.Host, error) {

	return nil, nil
}

func (i *HostServiceImpl) UpdateHost(ctx context.Context, request *host.UpdateHostRequest) (*host.Host, error) {

	return nil, nil
}

func (i *HostServiceImpl) DeleteHost(ctx context.Context, request *host.DeleteHostRequest) (*host.Host, error) {

	return nil, nil
}
