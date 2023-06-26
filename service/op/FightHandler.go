package op

import (
	"demo/dao/model/entity"
	"demo/dto"
	"demo/utils/opcode"
	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
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
func (c *fightHandler) Execute(ctx g.Ctx, context *dto.Context, ws *ghttp.WebSocket, msgArray []byte) interface{} {
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
	return false
}

// FightContext 战斗游戏上下文
type fightContext struct {
	Players []*playerContext
	OnLine  []bool
	RoomId  int64
}

// FightRoomList [roomId]
var FightRoomList = glist.New(true)

// FightRoomMap <RoomId,fightContext>
var FightRoomMap = make(map[int64]*fightContext)
var PlayerIdRoomIdMap = make(map[int64]int64)

var FightRoomLock sync.RWMutex

func NewFightContext(player1 *playerContext, player2 *playerContext, player3 *playerContext) *fightContext {
	FightRoomLock.RLock()
	defer FightRoomLock.Unlock()
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
			player1.MsgChan = make(chan string, 32)
		}
		PlayerIdRoomIdMap[player1.Id] = roomId
		i.Players[idx] = player1
		idx++
	}
	if player2 != nil {
		if player2.MsgChan == nil {
			player2.MsgChan = make(chan string, 32)
		}
		PlayerIdRoomIdMap[player2.Id] = roomId
		i.Players[idx] = player1
		idx++
	}
	if player3 != nil {
		if player3.MsgChan == nil {
			player3.MsgChan = make(chan string, 32)
		}
		PlayerIdRoomIdMap[player3.Id] = roomId
		i.Players[idx] = player1
		idx++

	}

	return i

}

type playerContext struct {
	Id      int64
	Player  *entity.Player
	WS      *ghttp.WebSocket
	MsgChan chan string
}

func NewPlayerContext(p *dto.Context, ws *ghttp.WebSocket) *playerContext {
	i := new(playerContext)
	i.Id = p.Id
	i.Player = p.Player
	i.WS = ws
	return i
}
