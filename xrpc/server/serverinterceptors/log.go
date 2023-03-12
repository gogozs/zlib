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
		reqStr := tools.ToJsonStringWithMaxLen(req, 1024)
		rspStr := tools.ToJsonStringWithMaxLen(resp, 1024)
		if err != nil {
			xlog.MsgItem("error", err).
				MsgItem("cost", cost).
				MsgItem("req", reqStr).
				MsgItem("rsp", rspStr).
				Error(ctx, "[REQUEST] %s", info.FullMethod)
		} else {
			xlog.MsgItem("cost", cost).
				MsgItem("req", reqStr).
				MsgItem("rsp", rspStr).
				Info(ctx, "[REQUEST] %s", info.FullMethod)

		}
	}(time.Now())

	return handler(ctx, req)
}
