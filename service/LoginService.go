package service

import (
	"context"
	"demo/DTO"
	"demo/dao/model/entity"
	"demo/utils"
	"demo/utils/wscode"
	"encoding/json"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/guuid"
	"strings"
)

type loginService struct {
}

var LoginService = newService()

func newService() *loginService {
	return &loginService{}
}

var tokenMap = make(map[string]*entity.Player, 4096)

func (c loginService) CheckLogin(ctx context.Context, token string) (res entity.Player) {
	ctx, span := gtrace.NewSpan(ctx, "main")
	defer span.End()

	if token == "" {
		return res
	}
	if tokenMap[token] != nil {
		return *tokenMap[token]
	}
	return res
}
func (c loginService) PutLogin(ctx context.Context, p *entity.Player) {
	ctx, span := gtrace.NewSpan(ctx, "main")
	defer span.End()
	tokenMap[p.Token] = p
}

// CheckPWD 检查密码登录
func (c *loginService) CheckPWD(ctx context.Context, ws *ghttp.WebSocket) (res entity.Player) {
	ctx, span := gtrace.NewSpan(ctx, "CheckPWD")
	defer span.End()
	var read = true
	for read {

		_, array, err := ws.ReadMessage()
		if err != nil {
			read = false
			glog.Info(ctx, "loginService读取消息错误", err)
			return res
		}
		if array[0] == byte(48) {
			utils.WriteWs(ws, gconv.String(wscode.CONNET), "")
			continue
		}
		glog.Info(ctx, "消息内容", string(array))
		if array[0] != byte(49) && byte(50) != array[0] {
			utils.WriteWs(ws, gconv.String(wscode.LOGIN_FAIL), "")
			continue
		}

		idx := strings.Index(string(array), ":")
		if idx != 1 {
			continue
		}

		//处理登录逻辑
		glog.Info(ctx, "处理登录逻辑", string(array[idx+1:]))
		var dto DTO.WsLoginDTO
		json.Unmarshal(array[idx+1:], &dto)
		if dto.Pwd == "" || dto.Phone == "" {
			continue
		}
		var player = doLogin(ctx, dto)
		if player.Phone == "" && array[0] == byte(49) {
			player = doRegister(ctx, dto)
		}
		if player.Id == 0 {
			utils.WriteWs(ws, gconv.String(wscode.LOGIN_FAIL), "")
			continue
		}
		player.Pwd = ""
		player.Token = guuid.New().String()
		utils.WriteWs(ws, gconv.String(wscode.USER_LOGIN), player)
		return player

	}
	return res
}

func doLogin(ctx context.Context, dto DTO.WsLoginDTO) (player entity.Player) {
	ctx, span := gtrace.NewSpan(ctx, "doLogin")
	defer span.End()
	//do.Player

	res, err := g.Model("player").Ctx(ctx).With().Where("phone = ", dto.Phone).Limit(0, 1).One()
	if err != nil {
		glog.Info(ctx, "查询失败", err)
		return player
	}
	res.Struct(&player)
	if player.Phone == "" {
		return player
	}
	md, err := gmd5.Encrypt(player.Salt + ":" + dto.Pwd)
	if md != player.Pwd {
		return entity.Player{}
	}

	return player
}
func doRegister(ctx context.Context, dto DTO.WsLoginDTO) (player entity.Player) {
	ctx, span := gtrace.NewSpan(ctx, "doRegister")
	defer span.End()

	res, err := g.Model("player").Ctx(ctx).With().Where("phone = ", dto.Phone).Limit(0, 1).Count()
	if err != nil {
		glog.Info(ctx, "查询失败", err)
		return player
	}
	if res > 0 {
		return player
	}
	md, err := gmd5.Encrypt(player.Salt + ":" + dto.Pwd)
	player.Pwd = md
	player.Phone = dto.Phone
	player.NickName = "玩家" + dto.Phone
	res2, err := g.Model("player").Ctx(ctx).Save(player)
	if err != nil {
		glog.Errorf(ctx, "保存错误", player, res2, err)
	}
	if id, err := res2.LastInsertId(); err == nil {
		player.Id = id
	}
	return player
}
