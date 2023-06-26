package op

import (
	"demo/dto"
	"demo/utils"
	"demo/utils/opcode"
	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gorilla/websocket"

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
func (c *matchHandler) Execute(ctx g.Ctx, context *dto.Context, ws *websocket.Conn, msgArray []byte) interface{} {
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
func (c *matchHandler) enterCreq(ctx g.Ctx, context *dto.Context, ws *websocket.Conn, array []byte) {
	var success int = 0

	//查询是否在队列中
	for i := MatchPlayerList.Front(); i != nil; i = i.Next() {
		if i.Value.(*playerContext).Id == context.Id {
			//更新ws对象，有可能已经断开
			glog.Info(ctx, "已经在等待队列中，直接返回")
			ws.WriteMessage(1, utils.EnSubCode(c.opCode, opcode.MatchCode.EnterSres, success))
			return
		}

	}

	matchRWLock.RLock()
	defer matchRWLock.RUnlock()

	//查询是否已经开始游戏
	if PlayerIdRoomIdMap[context.Id] > 0 {
		//已经在匹配队列中
		context.RoomId = PlayerIdRoomIdMap[context.Id]
		//开始游戏
		ws.WriteMessage(1, utils.EnSubCode(c.opCode, opcode.MatchCode.StartBro, PlayerIdRoomIdMap[context.Id]))
		return
	}
	//确认没有开始游戏，进入匹配队列

	item := NewPlayerContext(context, ws)
	MatchPlayerList.PushBack(item)
	if MatchPlayerList.Size() > 2 {
		//开启新房间
		array := MatchPlayerList.PopFronts(3)
		i1 := array[0]
		i2 := array[1]
		i3 := array[2]
		room := NewFightContext(i1.(*playerContext), i2.(*playerContext), i3.(*playerContext))
		context.RoomId = room.RoomId
		//广播进入游戏信息
		room.Broadcast(ctx, utils.EnSubCode(c.opCode, opcode.MatchCode.StartBro, room.RoomId))
	}
	ws.WriteMessage(1, utils.EnSubCode(c.opCode, opcode.MatchCode.EnterSres, success))
}

// leaveCreq 离开匹配队列
func (c *matchHandler) leaveCreq(ctx g.Ctx, context *dto.Context, ws *websocket.Conn, array []byte) {
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
