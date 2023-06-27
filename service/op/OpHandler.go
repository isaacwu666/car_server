package op

import (
	"demo/dto"
	"demo/utils/opcode"
	"github.com/gogf/gf/v2/frame/g"
	_ "github.com/gogf/gf/v2/net/ghttp"
)

type (
	iOpHandler interface {
		//RequireLogin 是否需要登录
		RequireLogin(ctx g.Ctx) bool
		//Execute 执行消息进程，
		Execute(ctx g.Ctx, context *dto.Context, msgArray []byte) interface{}
		//GetOpCode 获取操作编码
		GetOpCode() string
	}
	//Handler 执行
	baseOpHandler struct {
		opCode string
	}
	iBuild interface {
		Build(op int) *iOpHandler
	}
	build struct {
		iBuild
	}
)

var (
	Build = new(build)
)

// Build 构建
func (b build) Build(op string) (res iOpHandler) {
	switch op {
	case opcode.Beat:
		return BeatHandler
	case opcode.Account:
		return AccountHandler
	case opcode.User:
		return UserHandler
	case opcode.Match:
		return MatchHandler()
		break
	case opcode.Chat:
		return ChatHandler
		break
	case opcode.HPFight:
		return FightHandler
	}
	return res
}
