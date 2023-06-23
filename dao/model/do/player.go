// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Player is the golang structure of table player for DAO operations like Where/Data.
type Player struct {
	g.Meta   `orm:"table:player, do:true"`
	Id       interface{} // id
	Phone    interface{} // 电话
	Pwd      interface{} // 密码
	Salt     interface{} // pwd盐
	NickName interface{} // 昵称
}
