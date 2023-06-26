package op

import (
	"demo/dto"
	"demo/utils"
	"demo/utils/opcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type userHandler struct {
	baseOpHandler
}

var UserHandler = newUserHandler()

func newUserHandler() *userHandler {
	i := new(userHandler)
	i.opCode = opcode.User
	return i
}

// GetOpCode 获取操作编码
func (c *userHandler) GetOpCode() string {
	return c.opCode
}

// RequireLogin 是否需要登录
func (c *userHandler) RequireLogin(ctx g.Ctx) bool {
	return false
}

// Execute 执行消息进程，
func (c *userHandler) Execute(ctx g.Ctx, context *dto.Context, ws *ghttp.WebSocket, msgArray []byte) interface{} {
	ws.WriteMessage(1, utils.EnCode(c.opCode, msgArray))
	return false
}
