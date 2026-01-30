// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gozeroservicedemo/user/api/internal/svc"
	"gozeroservicedemo/user/api/internal/types"
	"gozeroservicedemo/user/sql/cachemodel"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var secret = []byte("I Like Kabocha")

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
	// 验证参数合法性
	// 雪花算法生成用户ID
	// 加密加盐

	// 参数校验--简单
	if len(req.UserName) == 0 || len(req.Password) == 0 || len(req.Repassword) == 0 {
		logx.Info("注册失败")
		return nil, fmt.Errorf("用户名和密码不能为空")
	}
	if len(req.Mobile) == 0 {
		return nil, fmt.Errorf("手机号不能为空")
	} // 实际业务场景下手机号的校验还是很复杂的，需要进行正则表达式校验
	if req.Password != req.Repassword {
		return nil, fmt.Errorf("两次输入的密码不一致")
	}
	logx.Infof("req: %#v", req)

	// 判断用户是否存在
	u, err := l.svcCtx.UserModel.FindOneByName(l.ctx, sql.NullString{String: req.UserName, Valid: true})

	// 在这里使用缓存模型进行查询

	// 查询失败
	if err != nil && err != sqlx.ErrNotFound {
		logx.Errorw("user_Redister_UserModel.FindOneByName error", logx.Field("err", err))
		return nil, errors.New("内部错误")
	}
	// 查到了
	if u != nil {
		return nil, errors.New("用户已存在")
	}
	// 没查到，可以创建

	// 加密加盐
	h := md5.New()
	h.Write([]byte(req.Password))
	h.Write(secret)
	passwordStr := hex.EncodeToString(h.Sum(nil))

	user := &cachemodel.User{
		Id:       time.Now().Unix(),
		Name:     sql.NullString{String: req.UserName, Valid: true}, // 假设用户名不能为空
		Password: passwordStr,                                       //不能存储明文密码，实际项目中需要加密
		Mobile:   req.Mobile,
		Gender:   int64(req.Gender),
		Nickname: req.NickName,
	}

	_, err = l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		return nil, err
	}
	return &types.RegisterResp{
		Success: true,
		Message: "注册成功",
		Id:      user.Id,
	}, nil
}
