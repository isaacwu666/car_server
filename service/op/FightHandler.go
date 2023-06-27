package op

import (
	"context"
	"demo/dto"
	"demo/utils"
	"demo/utils/opcode"
	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
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
	fightContext = c.LoadFightContext(ctx, context)
	//如果不存在，则返回
	if fightContext == nil {
		return false
	}
	mType, msg := utils.SplitMsg(msgArray)
	defer glog.Info(ctx, "处理", mType, "完成!", msg)
	switch mType {
	//准备就绪 5:12
	case opcode.FightCode.READY_CREQ:
		c.ready(ctx, context, fightContext)
		break

	}

	return false
}

// ready 准备
func (c *fightHandler) ready(ctx g.Ctx, context *dto.Context, fightContext *fightContext) {
	if context.CtxMap == nil || context.CtxMap.Get("rIdx") == nil {
		return
	}
	rIdxS, ok := (context.CtxMap.Get("rIdx")).(string)
	if !ok {
		return
	}
	rIdx := gconv.Int(rIdxS)
	if fightContext.Ready[rIdx] {
		//已经准备好了，则不用再次处理
		//广播准备信息
		fightContext.Broadcast(ctx, utils.EnSubCode(c.GetOpCode(),
			opcode.FightCode.READY_BRO, fightContext.Ready))
		return
	}
	fightContext.RLock.RLock()
	fightContext.OnLine[rIdx] = true
	fightContext.Ready[rIdx] = true
	cmd := fightCommand{
		FightCode: opcode.FightCode.READY_CREQ,
		Idx:       fightContext.idxStep,
		PlayerId:  context.Id,
		PlayerIdx: rIdx,
		CTime:     time.Now().Unix(),
		Array:     nil,
	}
	//增加命令行索引
	fightContext.idxStep++
	fightContext.Command.PushBack(cmd)
	if AllTrue(fightContext.Ready) {
		//全部就绪
		if fightContext.Status < FightStatusLandowner {
			//开始进入抢地主流程
			fightContext.Status = FightStatusLandowner
			fightContext.CmdPIdx = 0
		}
	}
	fightContext.RLock.RUnlock()
	//广播准备信息
	fightContext.Broadcast(ctx, utils.EnSubCode(c.GetOpCode(),
		opcode.FightCode.READY_BRO, fightContext.Ready))
	//广播抢地主
	if fightContext.Status == FightStatusLandowner {
		fightContext.Broadcast(ctx, utils.EnSubCode(c.GetOpCode(),
			opcode.FightCode.STATUS_BRO, fightContext))
	}

}
func AllTrue(bool2 []bool) bool {

	for _, b := range bool2 {
		if !b {
			return false
		}
	}
	return true
}

// LoadFightContext 加载房间信息
func (c *fightHandler) LoadFightContext(ctx g.Ctx, context *dto.Context) *fightContext {
	var fightContext *fightContext
	if context.RoomId == 0 {
		if PlayerIdRoomIdMap[context.Player.Id] == 0 {
			return fightContext
		}
		roomId := PlayerIdRoomIdMap[context.Player.Id]
		//如果房间信息为空，去掉脏水数据
		if FightRoomMap[roomId] == nil {
			FightRoomLock.RLock()
			if FightRoomMap[roomId] == nil {
				delete(PlayerIdRoomIdMap, context.Player.Id)
			}
			defer FightRoomLock.RUnlock()
			return fightContext
		}
		fightContext = FightRoomMap[roomId]

	} else if FightRoomMap[context.RoomId] == nil {
		context.RoomId = 0
		return nil
	} else {
		fightContext = FightRoomMap[context.RoomId]
	}

	//设置下标
	if context.CtxMap == nil {
		context.CtxMap = gmap.NewHashMap(true)
	}

	rid := context.CtxMap.Get("rIdx")
	rIdx, ok := (rid).(int64)
	if !ok || rIdx == 0 {
		for i, player := range fightContext.Players {
			if context.Id == player.Id {
				rIdx = int64(i)
				context.CtxMap.Set("rIdx", gconv.String(rIdx))
				break
			}
		}
	}

	return fightContext
}

// FightContext 战斗游戏上下文
type fightContext struct {
	Players []*dto.Context `json:"players,omitempty"`
	OnLine  []bool         `json:"onLine,omitempty"`
	Ready   []bool         `json:"ready,omitempty"`
	RoomId  int64          `json:"roomId,omitempty,string"`
	//房间写锁
	RLock sync.RWMutex `json:"-"`
	//当前的步数，用于保证时序
	idxStep int `json:"idxStep,omitempty"`
	//命令记录
	Command *glist.List `json:"command,omitempty"`
	//游戏状态。1等待中，5抢地主，10翻倍中，15发牌中，20开始游戏，30结算中，40退出游戏
	Status int `json:"status"`
	//当前可以执行命令的玩家下标，在抢地主，翻倍，发牌 状态中有效
	CmdPIdx int `json:"cmdPIdx"`
}

var (
	FightStatusWait         = 1  //1等待中
	FightStatusLandowner    = 5  //5抢地主
	FightStatusDouble       = 10 //10翻倍中
	FightStatusDealing      = 15 //15发牌中
	FightStatusInGame       = 20 //20开始游戏
	FightStatusInSettlement = 30 //30结算中
	FightStatusExit         = 40 //40退出游戏
)

type fightCommand struct {
	//命令编码
	FightCode string `json:"fightCode,omitempty"`
	//命令标识
	Idx int `json:"idx,omitempty"`
	//玩家Id
	PlayerId int64 `json:"playerId,omitempty,string"`
	//玩家下标
	PlayerIdx int `json:"playerIdx,omitempty"`
	//命令执行时间戳
	CTime int64 `json:"CTime,omitempty,string"`
	//命令内容
	Array []byte `json:"array,omitempty"`
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
		Ready:   make([]bool, 3),
		Players: make([]*dto.Context, 3),
		Command: glist.New(true),
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
