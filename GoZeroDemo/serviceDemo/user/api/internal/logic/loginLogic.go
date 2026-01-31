// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"gozeroservicedemo/user/api/internal/svc"
	"gozeroservicedemo/user/api/internal/types"

	"github.com/golang-jwt/jwt/v4"
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
		return &types.LoginResp{
			Message: "用户名或密码不能为空",
		}, fmt.Errorf("用户名或密码不能为空")
	}
	// 现在检查正确性
	// 检查用户并取出密码
	user, err := l.svcCtx.UserModel.FindOneByName(l.ctx, sql.NullString{String: req.UserName, Valid: true})
	if err != nil {
		logx.Errorw("user_Login_UserModel.FindOneByName error", logx.Field("err", err))
		return &types.LoginResp{
			Message: "用户不存在",
		}, fmt.Errorf("用户不存在")
	}
	// 检查密码是否正确
	h := md5.New()
	h.Write([]byte(req.Password))
	h.Write(secret)
	passwordStr := hex.EncodeToString(h.Sum(nil))
	if user.Password != passwordStr {
		logx.Infof("user_Login_password not match, userId: %d", user.Id)
		return &types.LoginResp{
			Message: "密码错误",
		}, fmt.Errorf("密码错误")
	}

	// 将 iat设置为当前时间戳
	iat := time.Now().Unix()
	token, err := getJwtToken(l.svcCtx.Config.Auth.AccessSecret, iat, l.svcCtx.Config.Auth.AccessExpire, fmt.Sprintf("%d", user.Id))
	// 登录成功，生成 token 返回
	if err != nil {
		logx.Errorw("user_Login_getJwtToken error", logx.Field("err", err))
		return nil, fmt.Errorf("生成 token 失败")
	}

	return &types.LoginResp{
		Message:      "登录成功",
		AccessToken:  token,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		RefreshAfter: l.svcCtx.Config.Auth.AccessExpire / 2,
	}, nil
}

// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @payload: 数据载体
func getJwtToken(secretKey string, iat, seconds int64, payload string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["payload"] = payload
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
