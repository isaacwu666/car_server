package main

import (
	"demo/ws"
	"github.com/gin-gonic/gin"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/contrib/trace/jaeger/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"net/http"
)

func main() {
	ctx := gctx.New()
	tp, _ := jaeger.Init("hp_ws", "websocket")
	defer tp.Shutdown(ctx)
	s := g.Server()
	s.SetServerRoot(gfile.MainPkgPath())

	r := setupRouter()
	r.Run(":3014")

}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("./index.html")
	r.GET("hp_ws", ws.WS2)
	r.GET("index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	return r
}
