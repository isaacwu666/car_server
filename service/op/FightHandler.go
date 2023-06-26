package op

import (
	"context"
	"demo/dto"
	"demo/utils/opcode"
	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"sync"
	"time"
)

type fightHandler struct {
	baseOpHandler
}

var FightHandler = newFightHandler()

func newFightHandler() *fightHandler {
	i := new(fightHandler)
	i.opCode = opcode.HPFight
	return i
}

// GetOpCode 获取操作编码
func (c *fightHandler) GetOpCode() string {
	return c.opCode
}

// RequireLogin 是否需要登录
func (c *fightHandler) RequireLogin(ctx g.Ctx) bool {
	return true
}

// Execute 执行消息进程，
func (c *fightHandler) Execute(ctx g.Ctx, context *dto.Context, msgArray []byte) interface{} {
	var fightContext *fightContext
	if context.RoomId == 0 {
		if PlayerIdRoomIdMap[context.Player.Id] == 0 {
			return false
		}
		roomId := PlayerIdRoomIdMap[context.Player.Id]
		//如果房间信息为空，去掉脏水数据
		if FightRoomMap[roomId] == nil {
			FightRoomLock.RLock()
			if FightRoomMap[roomId] == nil {
				delete(PlayerIdRoomIdMap, context.Player.Id)
			}
			defer FightRoomLock.RUnlock()
			return false
		}
		fightContext = FightRoomMap[roomId]

	}
	//如果不存在，则返回
	if fightContext == nil {
		return false
	}
	return false
}

// FightContext 战斗游戏上下文
type fightContext struct {
	Players []*dto.Context
	OnLine  []bool
	RoomId  int64
}

var (
	//FightRoomList [roomId]
	FightRoomList = glist.New(true)
	//FightRoomMap <RoomId,fightContext>
	FightRoomMap = make(map[int64]*fightContext)
	//PlayerIdRoomIdMap <玩家ID，房间ID>
	PlayerIdRoomIdMap = make(map[int64]int64)
	//FightRoomLock 锁ID
	FightRoomLock sync.RWMutex
)

// Broadcast 向房间内玩家广播消息
func (c *fightContext) Broadcast(ctx g.Ctx, array []byte) {
	for _, player := range c.Players {
		glog.Info(ctx, "服务端主动给客户端", player.Player.Id, "发消息", string(array))
		player.SendMsg(array)
	}
}

// BroadcastIdx 向指定玩家广播信息
func (c *fightContext) BroadcastIdx(ctx g.Ctx, array []byte, idx int) {
	player := c.Players[idx]
	glog.Info(context.Background(), "服务端主动给客户端", player.Player.Id, "发消息", string(array))
	player.SendMsg(array)
}

func NewFightContext(player1 *dto.Context, player2 *dto.Context, player3 *dto.Context) *fightContext {
	FightRoomLock.RLock()
	defer FightRoomLock.RUnlock()
	var roomId int64 = time.Now().Unix()

	var i *fightContext = &fightContext{
		RoomId:  roomId,
		OnLine:  make([]bool, 3),
		Players: make([]*dto.Context, 3),
	}
	i.RoomId = roomId

	FightRoomList.PushBack(i)
	//放到Map中
	FightRoomMap[i.RoomId] = i
	var idx int = 0
	if player1 != nil {
		PlayerIdRoomIdMap[player1.Id] = roomId
		i.Players[idx] = player1
		idx++
	}
	if player2 != nil {
		PlayerIdRoomIdMap[player2.Id] = roomId
		i.Players[idx] = player2
		idx++
	}
	if player3 != nil {
		PlayerIdRoomIdMap[player3.Id] = roomId
		i.Players[idx] = player3
		idx++
	}

	return i

}
