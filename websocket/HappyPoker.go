package websocket

import (
	"demo/service"
	"demo/service/wshandler"
	"demo/utils"
	"demo/utils/wscode"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

type HappyPock struct {
}

func NewHappyPock() *HappyPock {
	return &HappyPock{}
}

func (HappyPock) Ws(r *ghttp.Request) {
	var ctx = gctx.New()
	ctx, span := gtrace.NewSpan(ctx, "Ws")
	defer span.End()
	//defer tp.

	ws, err := r.WebSocket()
	if err != nil {
		return
	}
	token := r.Header.Get("x-token")

	var player = service.LoginService.CheckLogin(ctx, token)

	if player.Id == 0 {
		player = service.LoginService.CheckPWD(ctx, ws)
		if player.Id != 0 {
			service.LoginService.PutLogin(ctx, &player)
		}
	}
	if player.Id == 0 {
		//连接已经断开
		ws.Close()
		return
	}

	var read = true
	ws.SetCloseHandler(func(code int, text string) error {
		read = false
		//客户端主动断开逻辑

		return nil
	})
	//循环读取
	for read {

		_, array, err := ws.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "close") {
				read = false
				continue
			}
			glog.Info(ctx, "读取消息错误", err)
			continue
		}
		glog.Info(ctx, "消息内容", string(array))
		if array[0] == byte(58) {
			ws.WriteMessage(1, array)
			continue
		}
		//分割消息
		mType, msg := utils.SplitMsg(array)
		//根据消息找到处理类
		switch mType {
		case wscode.USER_LOGIN:
		case wscode.UserRegister:
			utils.WriteWs(ws, gconv.String(wscode.USER_LOGIN), player)
			continue
		default:
			handler := wshandler.Factory.Build(ctx, mType, msg)
			(handler).Processed(ctx, player, ws, mType, msg)
			continue
		}
	}
}
