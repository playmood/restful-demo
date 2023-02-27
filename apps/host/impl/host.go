package impl

import (
	"context"
	"fmt"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/sqlbuilder"
	"github.com/playmood/restful-demo/apps/host"
)

// controller层
func (i *HostServiceImpl) CreateHost(ctx context.Context, ins *host.Host) (*host.Host, error) {
	// 直接打印日志
	i.l.Debug("create host")
	// 带format的日志打印，fmt.Sprintf()
	i.l.Debugf("create host %s", ins.Name)
	// 携带额外meta数据，常用于Trace系统
	i.l.With(logger.NewAny("request-id", "req01")).Debug("create host with meta kv")
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
	b := sqlbuilder.NewBuilder(QueryHostSQL)
	if request.KeyWords != "" {
		b.Where("r.`name` LIKE ? OR r.description LIKE ? OR r.private_ip LIKE ? OR r.public_ip LIKE ?", "%"+request.KeyWords+"%",
			"%"+request.KeyWords+"%",
			request.KeyWords+"%", request.KeyWords+"%")
	}
	b.Limit(request.GetOffset(), request.GetPageSize())
	querySQL, args := b.Build()
	i.l.Debugf("query sql: %s, ", querySQL, args)

	// prepare语句执行查询SQL
	stmt, err := i.db.PrepareContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	set := host.NewHostSet()
	for rows.Next() {
		// 每扫描一行，就要读取出来
		ins := host.NewHost()
		if err := rows.Scan(&ins.Id, &ins.Vendor, &ins.Region, &ins.CreateAt, &ins.ExpireAt, &ins.Type, &ins.Name,
			&ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt, &ins.Account, &ins.PublicIP, &ins.PrivateIP,
			&ins.CPU, &ins.Memory, &ins.GPUSpec, &ins.GPUAmount, &ins.OSType, &ins.OSName, &ins.SerialNumber); err != nil {
			return nil, err
		}
		set.Add(ins)
	}

	// total统计
	countSQL, args := b.BuildCount()
	i.l.Debugf("count sql: %s, args: %v", countSQL, args)
	countStmt, err := i.db.PrepareContext(ctx, countSQL)
	if err != nil {
		return nil, err
	}
	defer countStmt.Close()
	if err := countStmt.QueryRowContext(ctx, args...).Scan(&set.Total); err != nil {
		return nil, err
	}
	return set, nil
}

func (i *HostServiceImpl) DescribeHost(ctx context.Context, request *host.DescribeHostRequest) (*host.Host, error) {
	b := sqlbuilder.NewBuilder(QueryHostSQL)
	b.Where("r.id = ?", request.Id)
	describeSQL, args := b.Build()
	i.l.Debugf("describe sql: %s, ", describeSQL, args)
	// prepare语句执行查询SQL
	stmt, err := i.db.PrepareContext(ctx, describeSQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	ins := host.NewHost()
	err = stmt.QueryRowContext(ctx, args...).Scan(&ins.Id, &ins.Vendor, &ins.Region, &ins.CreateAt, &ins.ExpireAt, &ins.Type, &ins.Name,
		&ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt, &ins.Account, &ins.PublicIP, &ins.PrivateIP,
		&ins.CPU, &ins.Memory, &ins.GPUSpec, &ins.GPUAmount, &ins.OSType, &ins.OSName, &ins.SerialNumber)
	if err != nil {
		return nil, err
	}
	return ins, nil
}

func (i *HostServiceImpl) UpdateHost(ctx context.Context, request *host.UpdateHostRequest) (*host.Host, error) {
	// 获取已有对象
	ins, err := i.DescribeHost(ctx, host.NewDescribeHostRequestWithId(request.Id))
	if err != nil {
		return nil, err
	}
	// 根据更新模式更新对象
	switch request.Mode {
	case host.PATCH:
		if err := ins.Patch(request.Host); err != nil {
			return nil, err
		}
	case host.PUT:
		if err := ins.Put(request.Host); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("update mode only support PUT and PATCH")
	}
	// 更新数据库中的数据
	if err := ins.Validate(); err != nil {
		return nil, err
	}
	if err := i.update(ctx, ins); err != nil {
		return nil, err
	}
	// 返回更新后的对象
	return ins, nil
}

func (i *HostServiceImpl) DeleteHost(ctx context.Context, request *host.DeleteHostRequest) (*host.Host, error) {
	// 检测id是否存在
	ins, err := i.DescribeHost(ctx, host.NewDescribeHostRequestWithId(request.Id))
	if err != nil {
		return nil, fmt.Errorf("input Id is invalid")
	}
	if err := i.delete(ctx, ins); err != nil {
		return nil, err
	}

	return ins, nil
}
