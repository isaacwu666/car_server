package op

import (
	"demo/dao/model/entity"
	"demo/dto"
	"demo/utils"
	"demo/utils/opcode"
	"demo/utils/wscode"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/guuid"
)

type accountHandler struct {
	baseOpHandler
}

var AccountHandler = newAccountHandler()

func newAccountHandler() *accountHandler {
	i := new(accountHandler)
	i.opCode = opcode.Account
	return i
}

// GetOpCode 获取操作编码
func (c *accountHandler) GetOpCode() string {
	return c.opCode
}

// RequireLogin 是否需要登录
func (c *accountHandler) RequireLogin(ctx g.Ctx) bool {
	return false
}

// Execute 执行消息进程，
func (c *accountHandler) Execute(ctx g.Ctx, context *dto.Context, msgArray []byte) interface{} {
	mtype, array := utils.SplitMsg(msgArray)
	switch mtype {

	case wscode.AccountCode.Login:
		//1:2:{"phone":"1","pwd":"1"}
		return doLogin(ctx, context, array)
		break
	case wscode.AccountCode.RegistCreq:
		//1:2:{"phone":"1","pwd":"1"}
		return doRegister(ctx, context, array)
		break

	case wscode.AccountCode.NICK_NAME_CREQ:
		return doReNickName(ctx, context, array)
		break
	}
	return false
}

func doReNickName(ctx g.Ctx, context *dto.Context, array []byte) interface{} {
	if !context.IsLogin {
		return nil
	}
	d := context.Player
	toUpdate := entity.Player{NickName: string(array)}

	dbRes, err := g.Model("player").Ctx(ctx).Fields("nick_name").Where("id = ", d.Id).Update(toUpdate)
	if err != nil || dbRes == nil {

		context.SendMsg(utils.EnSubCode(AccountHandler.opCode,
			wscode.AccountCode.NICK_NAME_SRES,
			-1))
		return false
	}
	d.NickName = toUpdate.NickName
	context.SendMsg(utils.EnSubCode(AccountHandler.opCode,
		wscode.AccountCode.NICK_NAME_SRES,
		d))
	return true
}
func doLogin(ctx g.Ctx, context *dto.Context, msgArray []byte) (res *entity.Player) {
	d := &entity.Player{}
	err := gconv.Scan(msgArray, d)
	if err != nil {
		return res
	}
	dbRes, err := g.Model("player").Ctx(ctx).With("id").Where("phone = ", d.Phone).Limit(0, 1).One()
	if err != nil || dbRes == nil {
		glog.Info(ctx, "查询失败", err)
		context.SendMsg(utils.EnSubCode(AccountHandler.opCode,
			wscode.AccountCode.LoginFail,
			-1))
		return res
	}
	dbRes.Struct(&res)

	md, err := gmd5.Encrypt(res.Salt + ":" + d.Pwd)
	if md != res.Pwd {
		context.SendMsg(utils.EnSubCode(AccountHandler.opCode,
			wscode.AccountCode.LoginFail,
			-1))
		return res
	}
	res.Token = dto.NewToken(res.Id)
	res.Pwd = ""
	context.SendMsg(utils.EnSubCode(AccountHandler.opCode,
		wscode.AccountCode.Login,
		res))

	return res
}
func doRegister(ctx g.Ctx, context *dto.Context, msgArray []byte) (res *entity.Player) {

	err := gconv.Scan(msgArray, &res)
	res.Id = 0
	if err != nil {
		return res
	}
	dbRes, err := g.Model("player").Ctx(ctx).With("id").Where("phone = ", res.Phone).Limit(0, 1).Count()
	if err != nil {
		glog.Info(ctx, "查询失败", err)

		context.SendMsg(utils.EnSubCode(AccountHandler.opCode,
			wscode.AccountCode.RegistSres,
			-1))
		return res
	}
	if dbRes > 0 {

		context.SendMsg(utils.EnSubCode(AccountHandler.opCode,
			wscode.AccountCode.RegistSres,
			-1))
		return res
	}

	res.Salt = guuid.New().String()[0:32]
	md, err := gmd5.Encrypt(res.Salt + ":" + res.Pwd)
	res.Pwd = md
	r, err := g.Model("player").Ctx(ctx).Save(&res)
	if err != nil {
		context.SendMsg(utils.EnSubCode(AccountHandler.opCode,
			wscode.AccountCode.RegistSres,
			-1))
		return res
	}
	id, _ := r.LastInsertId()
	res.Id = id
	res.Token = dto.NewToken(res.Id)
	res.Pwd = ""
	context.SendMsg(utils.EnSubCode(AccountHandler.opCode,
		wscode.AccountCode.RegistSres,
		res))

	return res
}
