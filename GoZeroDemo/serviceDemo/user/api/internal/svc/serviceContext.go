// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"gozeroservicedemo/user/api/internal/config"
	"gozeroservicedemo/user/sql/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlxConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(sqlxConn),
	}
}
