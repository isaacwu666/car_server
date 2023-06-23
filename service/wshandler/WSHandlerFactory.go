package wshandler

import (
	"context"
	"demo/utils/wscode"
)

type factory struct {
}

var Factory = &factory{}

func (receiver factory) Build(ctx context.Context, mType int, msg any) WsHandler {
	switch mType {
	case wscode.JOIN_ROOM:
		return NewHappyPokerRoomHandler()
		break
	default:
		return EmptyHandler
		break
	}

	//返回默认的执行方案
	return EmptyHandler
}
