package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"mall/common/cryptx"
	"mall/service/user/model"
	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	_, err := l.svcCtx.UserModel.FindOneByMobile(in.Mobile)
	if err == nil {
		return nil, status.Error(100, "该用户已存在")
	}

	if err != model.ErrNotFound {
		return nil, status.Error(500, err.Error())
	}

	newUser := model.User{
		Name:     in.Name,
		Gender:   in.Gender,
		Mobile:   in.Mobile,
		Password: cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password),
	}

	res, err := l.svcCtx.UserModel.Insert(&newUser)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	newUser.Id, err = res.LastInsertId()
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	return &user.RegisterResponse{
		Id:     newUser.Id,
		Name:   newUser.Name,
		Gender: newUser.Gender,
		Mobile: newUser.Mobile,
	}, nil
}
