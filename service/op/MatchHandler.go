package op

import (
	"demo/dto"
	"demo/utils"
	"demo/utils/opcode"
	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"sync"
)

type matchHandler struct {
	baseOpHandler
}

func MatchHandler() *matchHandler {

	i := new(matchHandler)
	//3
	i.opCode = opcode.Match
	return i
}

// GetOpCode 获取操作编码
func (c *matchHandler) GetOpCode() string {
	return c.opCode
}

// RequireLogin 是否需要登录
func (c *matchHandler) RequireLogin(ctx g.Ctx) bool {
	return true
}

// Execute 执行消息进程，
func (c *matchHandler) Execute(ctx g.Ctx, context *dto.Context, ws *ghttp.WebSocket, msgArray []byte) interface{} {
	mtyp, array := utils.SplitMsg(msgArray)
	switch mtyp {
	//0
	case opcode.MatchCode.EnterCreq:
		c.enterCreq(ctx, context, ws, array)
		break
	//2
	case opcode.MatchCode.LeaveCreq:
		c.leaveCreq(ctx, context, ws, array)
		break
	}
	return true
}

// enterCreq 进入匹配队列
func (c *matchHandler) enterCreq(ctx g.Ctx, context *dto.Context, ws *ghttp.WebSocket, array []byte) {
	var success int = 1

	//查询是否在队列中
	for i := MatchPlayerList.Front(); i != nil; i = i.Next() {
		if i.Value.(*playerContext).Id == context.Id {
			//更新ws对象，有可能已经断开
			i.Value.(*playerContext).WS = ws
			glog.Info(ctx, "已经在等待队列中，直接返回")
			ws.WriteMessage(1, utils.EnSubCode(c.opCode, opcode.MatchCode.EnterSres, success))
			return
		}

	}

	matchRWLock.RLock()
	defer matchRWLock.RUnlock()
	item := NewPlayerContext(context, ws)
	MatchPlayerList.PushBack(item)
	if MatchPlayerList.Size() > 2 {
		//开启新房间
		i1 := MatchPlayerList.Back()
		i2 := MatchPlayerList.Back()
		i3 := MatchPlayerList.Back()
		room := NewFightContext(i1.Value.(*playerContext), i2.Value.(*playerContext), i3.Value.(*playerContext))
		context.RoomId = room.RoomId
	}
	ws.WriteMessage(1, utils.EnSubCode(c.opCode, opcode.MatchCode.EnterSres, success))
}

// leaveCreq 离开匹配队列
func (c *matchHandler) leaveCreq(ctx g.Ctx, context *dto.Context, ws *ghttp.WebSocket, array []byte) {
	var success int = 1
	//查询是否在队列中
	for i := MatchPlayerList.Front(); i != nil; i = i.Next() {
		if i.Value.(*playerContext).Id == context.Id {
			matchRWLock.RLock()
			MatchPlayerList.Remove(i)
			matchRWLock.RUnlock()
		}

	}
	ws.WriteMessage(1, utils.EnSubCode(c.opCode, opcode.MatchCode.LeaveBro, success))
}

// matchRWLock MatchPlayerList读写锁
var matchRWLock sync.RWMutex

// MatchPlayerList 匹配队列
var MatchPlayerList = glist.New(true)
