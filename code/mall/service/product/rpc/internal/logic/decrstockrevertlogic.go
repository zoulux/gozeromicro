package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"

	"mall/service/product/rpc/internal/svc"
	"mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DecrStockRevertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecrStockRevertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecrStockRevertLogic {
	return &DecrStockRevertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DecrStockRevertLogic) DecrStockRevert(in *product.DecrStockRequest) (*product.DecrStockResponse, error) {
	db, err := sqlx.NewMysql(l.svcCtx.Config.Mysql.DataSource).RawDB()
	if err != nil {
		l.Error(err)
		return nil, status.Error(500, err.Error())
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		l.Error(err)
		return nil, status.Error(500, err.Error())
	}

	l.Logger.Info("DecrStockRevert", "barrier.CallWithDB")
	err = barrier.CallWithDB(db, func(tx *sql.Tx) error {
		_, err := l.svcCtx.ProductModel.TxAdjustStock(tx, in.Id, 1)
		return err
	})
	if err != nil {
		l.Error(err)
		return nil, err
	}

	return &product.DecrStockResponse{}, nil
}
