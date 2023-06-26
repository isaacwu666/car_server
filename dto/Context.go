package dto

import (
	"demo/dao/model/entity"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/guuid"
	"github.com/google/uuid"
	"time"
)

type Context struct {
	Id      int64
	RoomId  int64
	Player  *entity.Player
	CtxMap  map[string]any
	IsLogin bool
}

var ContextHolder = gmap.NewHashMap(true)
var ContextHolderTime = gmap.NewHashMap(true)
var UserIdTokenMap = gmap.NewHashMap(true)

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
		ContextHolderTime.Set(player.Token, time.Now().Unix())
	} else {
		player = &entity.Player{}
		context = new(Context)
		context.Player = player
	}
	return context
}
func (c *Context) SetPlayer(p *entity.Player) {
	if p.Token == "" {
		p.Token = uuid.New().String()[0:32]
	}
	c.Id = p.Id
	c.Player = p
	c.IsLogin = true
	ContextHolderTime.Set(p.Token, time.Now().Unix())
	ContextHolder.Set(p.Token, c)
}
