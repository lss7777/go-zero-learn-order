package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Context key 类型，避免冲突（不能是空结构体，零大小类型指针可能相同）
type ctxKey int

// AppKey 和 TokenKey 是 gRPC metadata 中的字段名，对应 go-zero auth 的约定
const (
	AppKey   = "app"
	TokenKey = "token"
)

var (
	// CtxAppKey  context 中存储 app 的 key
	CtxAppKey = ctxKey(1)
	// CtxTokenKey context 中存储 token 的 key
	CtxTokenKey = ctxKey(2)
)

// AuthUnaryInterceptor 从 context 中取出 app 和 token，注入 gRPC metadata
func AuthUnaryInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	app, _ := ctx.Value(CtxAppKey).(string)
	token, _ := ctx.Value(CtxTokenKey).(string)

	if app != "" || token != "" {
		md := metadata.Pairs(AppKey, app, TokenKey, token)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}
