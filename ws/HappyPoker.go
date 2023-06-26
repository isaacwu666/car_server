package ws

import (
	"demo/dao/model/entity"
	"demo/dto"
	"demo/service/op"
	"demo/utils"
	"demo/utils/opcode"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gorilla/websocket"

	"net/http"
	"strings"
)

type HappyPock struct {
}

func NewHappyPock() *HappyPock {
	return &HappyPock{}
}

//	func (HappyPock) Ws(r *ghttp.Request) {
//		var ctx = gctx.New()
//		ctx, span := gtrace.NewSpan(ctx, "Ws")
//		defer span.End()
//		//defer tp.
//
//		ws, err := r.WebSocket()
//		if err != nil {
//			return
//		}
//		defer ws.Close()
//
//		x := r.Header.Get("x-token")
//		context := dto.NewContext(x)
//
//		//切换操作类型
//		var onLine = true
//		for onLine {
//			i, array, err := ws.ReadMessage()
//			if err != nil && strings.Contains(err.Error(), "close") {
//				onLine = false
//				continue
//			}
//			opCode, msgArray := utils.SplitMsg(array)
//			opHandler := op.Build.Build(opCode)
//			if opHandler != nil {
//				if (!(opHandler).RequireLogin(ctx)) ||
//					((opHandler).RequireLogin(ctx) && context.IsLogin) {
//					res := opHandler.Execute(ctx, context, ws, msgArray)
//					if res != nil && opHandler.GetOpCode() == opcode.Account {
//						if val, ok := res.(*entity.Player); ok && val != nil {
//							val.Token = dto.NewToken(val.Id)
//							context.SetPlayer(val)
//						}
//
//					}
//					continue
//				}
//			}
//
//			//如果没有处理，则将消息直接返回
//			ws.WriteMessage(i, array)
//		}
//		ws.Close()
//	}
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WS2(c *gin.Context) {
	go doWS2(c)
}
func doWS2(c *gin.Context) {
	var ctx = gctx.New()
	ctx, span := gtrace.NewSpan(ctx, "Ws")
	defer span.End()

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	x := c.GetHeader("x-token")
	context := dto.NewContext(x)
	//切换操作类型
	var onLine = true
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
				res := opHandler.Execute(ctx, context, ws, msgArray)
				if res != nil && opHandler.GetOpCode() == opcode.Account {
					if val, ok := res.(*entity.Player); ok && val != nil {
						val.Token = dto.NewToken(val.Id)
						context.SetPlayer(val)
					}

				}
				continue
			}
		}

		//如果没有处理，则将消息直接返回
		ws.WriteMessage(i, array)
	}
	ws.Close()
}
