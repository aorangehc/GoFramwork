// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"database/sql"
	"time"

	"gozeroservicedemo/user/api/internal/svc"
	"gozeroservicedemo/user/api/internal/types"
	"gozeroservicedemo/user/sql/cachemodel"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateLogic) UserUpdate(req *types.UserUpdateReq) (resp *types.UserUpdateResp, err error) {
	// 已经经过 token 的验证了
	// 先检查参数的合法性
	if req.Id == 0 || req.UserName == "" || req.Password == "" {
		return &types.UserUpdateResp{
			Success: false,
			Message: "参数错误",
		}, nil
	}

	// 检查用户是否存在
	_, err = l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		logx.Errorw("user_UserUpdate_UserModel.FindOne error", logx.Field("err", err))
		return &types.UserUpdateResp{
			Success: false,
			Message: "用户不存在",
		}, nil
	}

	// 构建结构体
	user := &cachemodel.User{
		Id:       req.Id,
		Name:     sql.NullString{String: req.UserName, Valid: true},
		Nickname: req.NickName,
		Gender:   int64(req.Gender),
		Mobile:   req.Mobile,
		UpdateAt: time.Now(),
	}

	// 更新参数

	err = l.svcCtx.UserModel.Update(l.ctx, user)
	if err != nil {
		logx.Errorw("user_UserUpdate_UserModel.Update error", logx.Field("err", err))
		return &types.UserUpdateResp{
			Success: false,
			Message: "用户信息更新失败",
		}, nil
	}

	// 返回
	return &types.UserUpdateResp{
		Success: true,
		Message: "用户信息更新成功",
	}, nil
}
