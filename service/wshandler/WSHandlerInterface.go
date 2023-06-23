package wshandler

import (
	"context"
	"demo/dao/model/entity"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

type WsHandler interface {
	Processed(ctx context.Context, player entity.Player, ws *ghttp.WebSocket, mType int, msg any)
}

var EmptyHandler = &emptyHandler{}

type emptyHandler struct {
	WsHandler
}

func (receiver *emptyHandler) Processed(ctx context.Context, player entity.Player, ws *ghttp.WebSocket, mType int, msg any) {
	glog.Info(ctx, "未知数据类型")
}
