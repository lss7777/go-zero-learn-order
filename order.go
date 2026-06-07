// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"flag"
	"fmt"
	"net/http"

	"order/internal/config"
	"order/internal/errors"
	"order/internal/handler"
	"order/internal/response"
	"order/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"

	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

var configFile = flag.String("f", "etc/order-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 全局错误处理器：AppError → ApiResponse，其余降级为参数错误
	httpx.SetErrorHandler(func(err error) (int, any) {
		if appErr, ok := err.(*errors.AppError); ok {
			return http.StatusOK, response.Error(appErr.Code, appErr.Msg)
		}
		return http.StatusBadRequest, response.Error(errors.CodeInvalidParam, err.Error())
	})

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
