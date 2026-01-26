// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"database/sql"
	"fmt"

	"gozeroservicedemo/user/api/internal/svc"
	"gozeroservicedemo/user/api/internal/types"
	"gozeroservicedemo/user/sql/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 实现注册，将用户信息写到数据库中，并进行唯一性检查
	if len(req.UserName) == 0 || len(req.Password) == 0 {
		return nil, fmt.Errorf("用户名和密码不能为空")
	}
	// 雪花算法生成用户ID
	// userId := l.svcCtx.Config.Snowflake.Start()
	// 检查用户名是否存在
	var id int64
	id = 111
	_, err = l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Id:       id,
		Name:     sql.NullString{String: req.UserName, Valid: true}, // 假设用户名不能为空
		Password: req.Password,                                      //不能存储明文密码，实际项目中需要加密
		Mobile:   req.Mobile,
		Gender:   req.Gender,
		Nickname: req.NickName,
	})
	if err != nil {
		return nil, err
	}
	return &types.RegisterResp{
		Success: true,
		Message: "注册成功",
		Id:      id,
	}, nil
}
