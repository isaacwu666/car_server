package main

import (
	"demo/ws"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/contrib/trace/jaeger/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"time"
)

func main() {
	ctx := gctx.New()
	tp, _ := jaeger.Init("hp_ws", "websocket")
	defer tp.Shutdown(ctx)
	s := g.Server()

	var hp = ws.NewHappyPock()
	//var sc ghttp.ServerConfig=ghttp.NewConfig()
	//s.SetConfig(sc)
	s.BindHandler("/hp_ws", hp.Ws)
	s.SetServerRoot(gfile.MainPkgPath())
	s.SetIdleTimeout(time.Duration.Abs(72))
	s.SetPort(3014)
	s.Run()
	TestDefault()
}
func TestDefault() {

}
