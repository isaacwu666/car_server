// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Player is the golang structure for table player.
type Player struct {
	Id       int64  `json:"id,string"       description:"id"`
	Phone    string `json:"phone"    description:"电话"`
	Pwd      string `json:"pwd"      description:"密码"`
	Salt     string `json:"salt"     description:"pwd盐"`
	NickName string `json:"nickName" description:"昵称"`
	Token    string `json:"token"  description:"临时登录token"`
}
