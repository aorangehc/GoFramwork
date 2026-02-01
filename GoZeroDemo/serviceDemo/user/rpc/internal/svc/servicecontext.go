package svc

import (
	"gozeroservicedemo/user/rpc/internal/config"
	"gozeroservicedemo/user/sql/cachemodel"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel cachemodel.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		UserModel: cachemodel.NewUserModel(conn, c.CacheRedis),
	}
}
