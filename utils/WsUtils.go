package utils

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

func WriteWs(ws *ghttp.WebSocket, bus string, data any) bool {
	if data == nil {
		return false
	}
	var out_data string
	switch value := data.(type) {
	case string:
		out_data = string(bus) + ":" + value
		break
	default:
		if o, err := json.Marshal(data); err != nil {
			return false
		} else {
			out_data = bus + ":" + string(o)
		}

		break
	}

	return ws.WriteMessage(1, []byte(out_data)) == nil

}
func ReadWS(ctx context.Context, ws *ghttp.WebSocket) (mType int, msg any) {
	_, array, err := ws.ReadMessage()
	if err != nil {
		glog.Info(ctx, "读取消息错误", err)
		return 0, nil
	}
	glog.Info(ctx, "消息内容", string(array))
	idx := strings.Index(string(array), ":")
	if idx < 1 {
		return -1, nil
	}
	t := string(array[0:idx])
	mType = gconv.Int(t)
	msg = array[idx+1:]
	return mType, msg
}

func SplitMsg(array []byte) (mType int, msg any) {

	idx := strings.Index(string(array), ":")
	if idx < 1 {
		return -1, nil
	}
	t := string(array[0:idx])
	mType = gconv.Int(t)
	msg = array[idx+1:]
	return mType, msg
}
