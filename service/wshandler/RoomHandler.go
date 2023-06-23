package wshandler

import (
	"context"
	"demo/dao/model/entity"
	"demo/service/game"
	"github.com/gogf/gf/v2/net/ghttp"
	"strings"
)

// 斗地主房间处理器
type happyPokerRoomHandler struct {
	WsHandler
}

func NewHappyPokerRoomHandler() *happyPokerRoomHandler {
	return &happyPokerRoomHandler{}
}

/*
*
这里的线程应该是负责 消息的读取和发布
*/
func (h *happyPokerRoomHandler) Processed(ctx context.Context, player entity.Player, ws *ghttp.WebSocket, mType int, msg any) {
	//判断是否是中途短线
	if h.isOffLine(player.Id) {
		//todo 中途加入房间
	}

	//加入

	//开启
	room := game.NewHappyPokerRoom()
	playerCtx := game.NewPokeRoomContext(ctx, player, ws)
	(room).AddPlayer(playerCtx)
	for playerCtx.InGame {
		//读取消息
		_, bu, err := ws.ReadMessage()
		if strings.Contains(err.Error(), "close") {
			//连接下线
			playerCtx.Online = false
			break
		}
		if bu[0] == byte(50) {
			//心跳信息
			continue
		}
		//将消息写入队列
		playerCtx.ToServer <- string(bu)
	}
	if !playerCtx.Online {
		playerCtx.Close()
	}

}

// isOffLine 判断是否是中途短线
func (h *happyPokerRoomHandler) isOffLine(id int64) bool {

	return true
}
