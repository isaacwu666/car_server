package op

import (
	"context"
	"demo/dao/model/entity"
	"demo/dto"
	"demo/utils/opcode"
	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gorilla/websocket"

	"strings"
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
func (c *fightHandler) Execute(ctx g.Ctx, context *dto.Context, ws *websocket.Conn, msgArray []byte) interface{} {
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
	Players []*playerContext
	OnLine  []bool
	RoomId  int64
}
type playerContext struct {
	Id      int64
	Player  *entity.Player
	MsgChan *glist.List
	WS      *websocket.Conn
	OnLine  bool
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
	for i, player := range c.Players {

		glog.Info(context.Background(), "服务端主动给客户端", player.Player.Id, "发消息", string(array))
		if player.OnLine {
			error := player.WS.WriteMessage(1, array)
			if error != nil && strings.Contains(error.Error(), "close") {
				c.OnLine[i] = false
				c.Players[i].OnLine = false
			} else if error == nil {
				c.OnLine[i] = true
				c.Players[i].OnLine = true
			}
		} else {
			//将消息发送到管道中
			c.Players[i].MsgChan.PushBack(array)
		}

	}
}

func NewFightContext(player1 *playerContext, player2 *playerContext, player3 *playerContext) *fightContext {
	FightRoomLock.RLock()
	defer FightRoomLock.RUnlock()
	var roomId int64 = time.Now().Unix()

	var i *fightContext = &fightContext{
		RoomId:  roomId,
		OnLine:  make([]bool, 3),
		Players: make([]*playerContext, 3),
	}
	i.RoomId = roomId

	FightRoomList.PushBack(i)
	//放到Map中
	FightRoomMap[i.RoomId] = i
	var idx int = 0
	if player1 != nil {
		if player1.MsgChan == nil {
			player1.MsgChan = glist.New(true)
		}
		PlayerIdRoomIdMap[player1.Id] = roomId
		i.Players[idx] = player1
		idx++
	}
	if player2 != nil {
		if player2.MsgChan == nil {
			player2.MsgChan = glist.New(true)
		}
		PlayerIdRoomIdMap[player2.Id] = roomId
		i.Players[idx] = player2
		idx++
	}
	if player3 != nil {
		if player3.MsgChan == nil {
			player3.MsgChan = glist.New(true)
		}
		PlayerIdRoomIdMap[player3.Id] = roomId
		i.Players[idx] = player3
		idx++

	}
	//发送消息

	return i

}

func NewPlayerContext(p *dto.Context, ws *websocket.Conn) *playerContext {
	i := new(playerContext)
	i.Id = p.Id
	i.Player = p.Player
	i.WS = ws
	i.OnLine = true
	return i
}
