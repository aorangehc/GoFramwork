// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"

	"gozeroservicedemo/user/api/internal/svc"
	"gozeroservicedemo/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	//先检查用户和密码是否为空
	if req.Password == "" || req.UserName == "" {
		return nil, fmt.Errorf("用户名或密码不能为空")
	}
	// 现在检查正确性
	// 检查用户并取出密码
	user, err := l.svcCtx.UserModel.FindOneByName(l.ctx, sql.NullString{String: req.UserName, Valid: true})
	if err != nil {
		logx.Errorw("user_Login_UserModel.FindOneByName error", logx.Field("err", err))
		return nil, fmt.Errorf("用户不存在")
	}
	// 检查密码是否正确
	h := md5.New()
	h.Write([]byte(req.Password))
	h.Write(secret)
	passwordStr := hex.EncodeToString(h.Sum(nil))
	if user.Password != passwordStr {
		logx.Infof("user_Login_password not match, userId: %d", user.Id)
		return nil, fmt.Errorf("密码错误")
	}
	return &types.LoginResp{
		Token: "114514",
	}, nil
}
