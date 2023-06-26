package op

import (
	"demo/dto"
	"demo/utils/opcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type beatHandler struct {
	baseOpHandler
}

var BeatHandler = NewBeatHandler()

func NewBeatHandler() *beatHandler {
	i := new(beatHandler)
	i.opCode = opcode.Beat
	return i
}

// RequireLogin 是否需要登录
func (c *beatHandler) RequireLogin(ctx g.Ctx) bool {
	//glog.Info(ctx, "BeatHandler是否需要登录")
	return false
}

// Execute 执行消息进程，
func (c *beatHandler) Execute(ctx g.Ctx, context *dto.Context, ws *ghttp.WebSocket, msgArray []byte) interface{} {
	//glog.Info(ctx, "BeatHandler执行消息进程")
	//ws.WriteMessage(1, utils.EnCode(c.opCode, msgArray))
	return false
}
func (c *beatHandler) GetOpCode() string {
	return c.opCode
}
