package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"order/internal/svc"
	"order/internal/types"

	"user-rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	// 1. 根据 order_no 查询订单
	order, err := l.svcCtx.OrdersModel.FindOneByOrderNo(l.ctx, req.OrderNo)
	if err != nil {
		return nil, fmt.Errorf("订单不存在: %w", err)
	}

	// 2. 调用 user-rpc 获取用户名
	userResp, err := l.svcCtx.UserRpc.GetProfile(l.ctx, &user.GetProfileRequest{
		UserId: int64(order.UserId),
	})
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 3. 组装响应
	return &types.SearchResponse{
		OrderInfo: types.OrderInfo{
			Id:              int64(order.Id),
			OrderNo:         order.OrderNo,
			UserId:          int64(order.UserId),
			Status:          order.Status,
			TotalAmount:     order.TotalAmount,
			DiscountAmount:  order.DiscountAmount,
			PayAmount:       order.PayAmount,
			PaymentMethod:   order.PaymentMethod,
			TransactionNo:   order.TransactionNo,
			PaymentTime:     nullTimeToStr(order.PaymentTime),
			ReceiverName:    order.ReceiverName,
			ReceiverPhone:   order.ReceiverPhone,
			ReceiverAddress: order.ReceiverAddress,
			ExpressCompany:  order.ExpressCompany,
			ExpressNo:       order.ExpressNo,
			Remark:          order.Remark,
			CreatedAt:       order.CreatedAt.Format(time.DateTime),
			UpdatedAt:       order.UpdatedAt.Format(time.DateTime),
		},
		UserName: userResp.Username,
	}, nil
}

func nullTimeToStr(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format(time.DateTime)
	}
	return ""
}
