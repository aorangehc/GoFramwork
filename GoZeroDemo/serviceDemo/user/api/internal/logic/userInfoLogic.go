// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"

	"gozeroservicedemo/user/api/internal/svc"
	"gozeroservicedemo/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// 参数校验
	if req.Id == 0 {
		return nil, fmt.Errorf("参数错误")
	}

	token := l.ctx.Value("payload")

	// 获取数据
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil && err != sqlx.ErrNotFound {
		logx.Errorw("user_UserInfo_UserModel.FindOne error", logx.Field("err", err))
		return nil, fmt.Errorf("错误: %v", err)
	}
	if err == sqlx.ErrNotFound {
		logx.Info("user_UserInfo_UserModel.FindOne not found", logx.Field("userId", req.Id))
		return nil, fmt.Errorf("用户不存在")
	}

	// 获取成功
	resp = &types.UserInfoResp{
		Id:       user.Id,
		UserName: user.Name.String,
		Gender:   int(user.Gender),
		Mobile:   user.Mobile,
		Token:    token.(string), // 这里输出的是 payload
		NickName: user.Nickname,
		CreateAt: user.CreateAt.Time.String(),
		UpdateAt: user.UpdateAt.String(),
	}

	return
}

// func (l *UserInfoLogic)
