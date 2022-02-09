package logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
	"mall/service/user/rpc/user"

	"mall/service/order/rpc/internal/svc"
	"mall/service/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRevertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRevertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRevertLogic {
	return &CreateRevertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRevertLogic) CreateRevert(in *order.CreateRequest) (*order.CreateResponse, error) {
	l.Logger.Info("CreateRevert", "NewMysql")
	db, err := sqlx.NewMysql(l.svcCtx.Config.Mysql.DataSource).RawDB()
	if err != nil {
		l.Error(err)
		return nil, status.Error(500, err.Error())
	}

	// 获取子事务屏障对象
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		l.Error(err)
		return nil, status.Error(500, err.Error())
	}

	l.Logger.Info("CreateRevert", "barrier.CallWithDB")

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		_, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{Id: in.Uid})
		if err != nil {
			l.Error(err)
			return fmt.Errorf("用户不存在")
		}

		resOrder, err := l.svcCtx.OrderModel.FindOneByUid(in.Uid)
		if err != nil {
			l.Error(err)
			return fmt.Errorf("订单不存在")
		}

		resOrder.Status = 9
		err = l.svcCtx.OrderModel.TxUpdate(tx, resOrder)
		if err != nil {
			l.Error(err)
			return fmt.Errorf("订单更新失败")
		}
		return nil
	}); err != nil {
		l.Error(err)
		return nil, status.Error(500, err.Error())
	}

	return &order.CreateResponse{}, nil
}
