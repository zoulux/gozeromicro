package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mall/service/product/rpc/internal/svc"
	"mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DecrStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecrStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecrStockLogic {
	return &DecrStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DecrStockLogic) DecrStock(in *product.DecrStockRequest) (*product.DecrStockResponse, error) {

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

	l.Logger.Info("DecrStock", "barrier.CallWithDB")

	err = barrier.CallWithDB(db, func(tx *sql.Tx) error {
		result, err := l.svcCtx.ProductModel.TxAdjustStock(tx, in.Id, -1)
		if err != nil {
			l.Error(err)
			return err
		}

		affected, err := result.RowsAffected()
		if err == nil && affected == 0 {
			l.Error(err)
			return dtmcli.ErrFailure
		}
		l.Error(err)

		return err
	})

	// 这种情况是库存不足，不再重试，走回滚
	if err == dtmcli.ErrFailure {
		l.Error(err)
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	if err != nil {
		l.Error(err)
		return nil, err
	}

	return &product.DecrStockResponse{}, nil
}
