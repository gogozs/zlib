package serverinterceptors

import (
	"context"
	"fmt"
	"time"

	"github.com/gogozs/zlib/tools"
	"github.com/gogozs/zlib/xlog"
	"google.golang.org/grpc"
)

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	resp interface{}, err error) {
	ctx = xlog.WrapTrace(ctx)
	defer func(start time.Time) {
		cost := fmt.Sprintf("%dms", time.Since(start).Milliseconds())
		if err != nil {
			xlog.MsgItem("error", err).
				MsgItem("cost", cost).
				MsgItem("req", tools.ToJsonString(req)).
				Error(ctx, "[REQUEST] %s", info.FullMethod)
		} else {
			xlog.MsgItem("cost", cost).
				MsgItem("req", tools.ToJsonString(req)).
				Info(ctx, "[REQUEST] %s", info.FullMethod)

		}
	}(time.Now())

	return handler(ctx, req)
}
