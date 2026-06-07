package svc

import (
	"order/internal/config"
	"order/model"

	"user-rpc/userservice"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	OrdersModel  model.OrdersModel
	UserRpc      userservice.UserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	return &ServiceContext{
		Config: c,
		OrdersModel: model.NewOrdersModel(conn, cache.CacheConf{
			{
				RedisConf: c.BizRedis,
				Weight:    100,
			},
		}),
		UserRpc: userservice.NewUserService(zrpc.MustNewClient(c.UserRpcConf)),
	}
}
