package ws

import (
	"context"
	"demo/dao/model/entity"
	"demo/dto"
	"demo/service/op"
	"demo/utils"
	"demo/utils/opcode"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gorilla/websocket"

	"net/http"
	"strings"
)

type HappyPock struct {
}

func NewHappyPock() *HappyPock {
	return &HappyPock{}
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WS2(c *gin.Context) {
	var ctx = gctx.New()
	ctx, span := gtrace.NewSpan(ctx, "Ws")
	defer span.End()
	doWS2(ctx, c)
}
func doWS2(ctx context.Context, c *gin.Context) {

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	x := c.GetHeader("x-token")
	context := dto.NewContext(x)
	//切换操作类型
	var onLine = true
	signal := make(chan bool, 1)

	defer close(signal)
	go doSend(ctx, ws, context, signal)
	for onLine {
		i, array, err := ws.ReadMessage()
		if err != nil && strings.Contains(err.Error(), "close") {
			onLine = false
			continue
		}
		opCode, msgArray := utils.SplitMsg(array)
		opHandler := op.Build.Build(opCode)
		if opHandler != nil {
			if (!(opHandler).RequireLogin(ctx)) ||
				((opHandler).RequireLogin(ctx) && context.IsLogin) {
				res := opHandler.Execute(ctx, context, msgArray)
				if res != nil && opHandler.GetOpCode() == opcode.Account {
					if val, ok := res.(*entity.Player); ok && val != nil {
						if context.SetPlayer(val, signal) {
							go doSend(ctx, ws, context, signal)
						}
					}

				}
				continue
			}
		}

		//如果没有处理，则将消息直接返回
		ws.WriteMessage(i, array)
	}
}
func doSend(ctx context.Context, ws *websocket.Conn, context *dto.Context, signal chan bool) {
	glog.Info(ctx, context.Id, "开启发送消息协程")
	for true {
		select {
		case msg := <-*context.SendChan:
			if msg == nil {
				return
			}
			ws.WriteMessage(1, msg)
		case <-signal:
			glog.Info(ctx, context.Id, "退出发送消息协程")
			return
		}
	}
}
