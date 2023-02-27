package impl

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/playmood/restful-demo/apps/host"
)

// 对象与数据库之间的转换

// 把Host对象保存到数据库内
func (i *HostServiceImpl) save(ctx context.Context, ins *host.Host) error {
	var err error
	// 把数据入库到 resource表和host表
	// 一次需要往2个表录入数据, 我们需要2个操作 要么都成功，要么都失败, 事务的逻辑
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start tx error, %s", err)
	}

	// 通过Defer处理事务提交方式
	// 有报错则rollback
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				i.l.Error("rollback error, $s", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				i.l.Error("commit error, $s", err)
			}

		}
	}()

	// 插入resource数据
	rstmt, err := tx.PrepareContext(ctx, InsertResourceSQL)
	if err != nil {
		return err
	}
	defer rstmt.Close()
	_, err = rstmt.ExecContext(ctx,
		ins.Id, ins.Vendor, ins.Region, ins.CreateAt, ins.ExpireAt, ins.Type,
		ins.Name, ins.Description, ins.Status, ins.UpdateAt, ins.SyncAt, ins.Account, ins.PublicIP,
		ins.PrivateIP,
	)
	if err != nil {
		return err
	}

	// 插入describe数据
	dstmt, err := tx.PrepareContext(ctx, InsertDescribeSQL)
	if err != nil {
		return err
	}
	defer dstmt.Close()

	_, err = dstmt.ExecContext(ctx,
		ins.Id, ins.CPU, ins.Memory, ins.GPUAmount, ins.GPUSpec,
		ins.OSType, ins.OSName, ins.SerialNumber,
	)
	if err != nil {
		return err
	}
	return nil
}

func (i *HostServiceImpl) update(ctx context.Context, ins *host.Host) error {
	var (
		err          error
		resourceStmt *sql.Stmt
		hostStmt     *sql.Stmt
	)
	// 开启事务
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// 通过Defer处理事务提交方式
	// 有报错则rollback
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				i.l.Error("rollback error, $s", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				i.l.Error("commit error, $s", err)
			}

		}
	}()
	// 更新 resource表
	resourceStmt, err = tx.PrepareContext(ctx, UpdateResourceSQL)
	if err != nil {
		return err
	}
	// vendor=?,region=?,expire_at=?,name=?,description=?
	_, err = resourceStmt.Exec(ins.Vendor, ins.Region, ins.ExpireAt, ins.Name, ins.Description, ins.Id)
	if err != nil {
		return err
	}
	// 更新 host表
	hostStmt, err = tx.PrepareContext(ctx, UpdateHostSQL)
	if err != nil {
		return err
	}
	_, err = hostStmt.ExecContext(ctx, ins.CPU, ins.Memory, ins.Id)
	if err != nil {
		return err
	}
	return nil
}
