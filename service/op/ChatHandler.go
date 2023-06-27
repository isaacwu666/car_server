package op

import (
	"demo/dto"
	"demo/utils"
	"demo/utils/opcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"strconv"
)

type chatHandler struct {
	baseOpHandler
}

var ChatHandler = &chatHandler{
	baseOpHandler: baseOpHandler{
		opCode: opcode.Chat,
	},
}

func (c *chatHandler) RequireLogin(ctx g.Ctx) bool {
	return true
}

func (c *chatHandler) Execute(ctx g.Ctx, context *dto.Context, msgArray []byte) interface{} {
	mType, msg := utils.SplitMsg(msgArray)
	switch mType {
	case opcode.ChatCode.CREQ:
		c.chatToRoom(ctx, context, msg)
		break
	}

	return false
}

func (c *chatHandler) GetOpCode() string {
	return c.opCode
}

func (c *chatHandler) chatToRoom(ctx g.Ctx, context *dto.Context, array []byte) {
	roomId, text := utils.SplitMsg(array)

	id, _ := strconv.ParseInt(roomId, 10, 64)
	if FightRoomMap[id] == nil {
		glog.Info(ctx, "房间号已经被释放，消息被丢弃!", text)
		return
	}
	room := FightRoomMap[id]
	toMsg := gconv.String(context.Id) + ":" + string(text)
	room.Broadcast(ctx, utils.EnSubCode(opcode.Chat, opcode.ChatCode.SRES, toMsg))
}
