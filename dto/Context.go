package dto

import (
	"context"
	"demo/dao/model/entity"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/guuid"
	"time"
)

// 用户上下文
type Context struct {
	Id       int64          `json:"id,omitempty,string"`
	RoomId   int64          `json:"roomId,omitempty,string"`
	Player   *entity.Player `json:"player,omitempty"`
	CtxMap   *gmap.Map      `json:"-"`
	IsLogin  bool           `json:"isLogin,omitempty"`
	SendChan *chan []byte   `json:"-"`
}

func CopyContent(from *Context, to *Context) {
	to.Id = from.Id
	to.RoomId = from.RoomId
	to.Player = from.Player
	to.CtxMap = from.CtxMap
	to.IsLogin = from.IsLogin
	to.SendChan = from.SendChan
}

// ContextHolder <token string,context *dto.Context>
var ContextHolder = gmap.NewHashMap(true)

// TokenHolderTime <token string, time int64>
var TokenHolderTime = gmap.NewHashMap(true)

// UserIdTokenMap <playerId int64 ，token string>
var UserIdTokenMap = gmap.NewHashMap(true)

func (c *Context) SendMsg(array []byte) bool {
	if len(*c.SendChan) > 31 {
		return false
	}
	*c.SendChan <- array
	return true
}

func NewToken(id int64) string {
	v := UserIdTokenMap.Get(id)

	if v != nil {
		value, _ := v.(string)
		return value
	} else {
		token := guuid.New().String()[0:32]
		UserIdTokenMap.Set(id, token)
		return token
	}

}
func NewContext(x string) *Context {
	var player *entity.Player
	var context *Context
	v := ContextHolder.Get(x)
	con, ok := (v).(*Context)
	if x != "" && v != nil && ok {

		context = con
		player = context.Player
		TokenHolderTime.Set(player.Token, time.Now().Unix())
	} else {
		player = &entity.Player{}
		context = new(Context)
		context.Player = player
		c := make(chan []byte, 32)
		context.SendChan = &(c)
	}

	return context
}
func (c *Context) SetPlayer(p *entity.Player, signal chan bool) bool {
	res := false
	if p.Token == "" {
		p.Token = NewToken(p.Id)
	}
	//如果用户信息没变则不需要更换
	if c.Player != nil && c.Player.Id == p.Id {
		return res
	}

	//判断是否有旧的消息
	v := ContextHolder.Get(p.Token)
	con, ok := (v).(*Context)
	if v != nil && ok {
		CopyContent(con, c)
		c.Id = p.Id
		c.Player = p
		c.IsLogin = true
		signal <- false
		//关闭当前使用的管道
		newChan := c.SendChan
		ContextHolder.Remove(p.Token)
		UserIdTokenMap.Remove(con.Id)
		TokenHolderTime.Remove(con.Id)
		if len(*newChan) > 0 {
			go func() {
				glog.Info(context.Background(), "更换发送消息管道，搬运管道内消息")
				for len(*newChan) > 0 {
					t := <-*newChan
					c.SendMsg(t)
				}
				glog.Info(context.Background(), "SetPlayer协程关闭newChan")
				//close(*newChan)
			}()
		}
		glog.Info(context.Background(), "SetPlayer关闭newChan")
		//close(*newChan)
		res = true
	} else {
		c.Id = p.Id
		c.Player = p
		c.IsLogin = true
	}

	ContextHolder.Remove(p.Token)
	UserIdTokenMap.Set(p.Id, p.Token)
	TokenHolderTime.Set(p.Token, time.Now().Unix())
	ContextHolder.Set(p.Token, c)
	return res
}
