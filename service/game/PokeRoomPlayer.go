package game

import (
	"context"
	"demo/dao/model/entity"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

// PokeRoomContext 斗地主玩家信息上下文
type PokeRoomContext struct {
	ctx    context.Context
	player entity.Player
	ws     *ghttp.WebSocket
	//服务端发送到客户端的管道
	ToPlayer chan string
	//
	ToServer chan string
	InGame   bool
	Online   bool
}

func (p *PokeRoomContext) Close() {
	p.Online = false
	close(p.ToServer)
	close(p.ToPlayer)
	p.ws.Close()
}
func (p *PokeRoomContext) GameEnd() {
	close(p.ToServer)
	close(p.ToPlayer)
}

func NewPokeRoomContext(ctx context.Context, player entity.Player, ws *ghttp.WebSocket) *PokeRoomContext {
	var d = &PokeRoomContext{ctx: ctx, player: player, ws: ws}
	d.ToPlayer = make(chan string, 32)
	d.ToServer = make(chan string, 32)
	d.Online = true
	go d.doSendMsg()
	return d
}
func (d *PokeRoomContext) doSendMsg() {
	for d.Online {
		d.ws.WriteMessage(1, []byte(<-d.ToPlayer))
	}
	glog.Info(d.ctx, "用户掉线，发送消息协程结束", d.Online)
}
