package logic

import (
	"context"
	"errors"

	"gozeroservicedemo/user/rpc/internal/svc"
	"gozeroservicedemo/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户信息
func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	// todo: add your logic here and delete this line
	// 根据 userId 获取用户信息
	userData, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)

	if err == sqlx.ErrNotFound {
		logx.Info("user_GetUserInfo_UserModel.FindOne not found", logx.Field("userId", in.UserId))
		return nil, err
	}

	if err != nil {
		logx.Errorw("user_GetUserInfo_UserModel.FindOne error", logx.Field("err", err))
		return nil, errors.New("获取用户信息失败")
	}

	return &user.GetUserInfoResponse{
		UserId:   userData.Id,
		UserName: userData.Name.String,
		Gender:   userData.Gender,
	}, nil
}
