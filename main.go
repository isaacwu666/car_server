package main

import (
	"demo/websocket"
	"github.com/gogf/gf/contrib/trace/jaeger/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)
import _ "github.com/gogf/gf/contrib/drivers/mysql/v2"

func main() {
	ctx := gctx.New()
	tp, _ := jaeger.Init("hp_ws", "websocket")
	defer tp.Shutdown(ctx)
	s := g.Server()

	var hp = websocket.NewHappyPock()
	s.BindHandler("/hp_ws", hp.Ws)
	s.SetServerRoot(gfile.MainPkgPath())
	s.SetPort(3014)
	s.Run()

}
