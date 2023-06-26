package op

import (
	"demo/dto"
	"demo/utils/opcode"
	"github.com/gogf/gf/v2/frame/g"
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
func (c *userHandler) Execute(ctx g.Ctx, context *dto.Context, msgArray []byte) interface{} {

	return false
}
