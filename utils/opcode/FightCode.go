package opcode

type fightCode struct {
	GRAB_LANDLORD_CREQ string "0" //客户端发起抢地主的请求
	GRAB_LANDLORD_BRO  string "1" //服务器广播抢地主的结果
	TURN_GRAB_BRO      string "2" //服务器广播下一个玩家抢地主的结果

	DEAL_CREQ string "3" //客户端发起出牌的请求
	DEAL_SRES string "4" //服务器给客户端出牌的响应
	DEAL_BRO  string "5" //服务器广播出牌的结果

	PASS_CREQ string "6" //客户端发起不出的请求
	PASS_SRES string "7" //服务器发给客户端不出的响应

	TURN_DEAL_BRO string "8" //服务器广播转换出牌的结果

	LEAVE_BRO string "9" //服务器广播有玩家退出游戏

	OVER_BRO string "10" //服务器广播游戏结束

	GET_CARD_SRES string "11" //服务器给客户端卡牌的响应

	//准备
	READY_CREQ  string "12"
	READY_BRO   string "13"
	STATUS_CREQ string "14" //客户端请求当前状态
	STATUS_BRO  string "15" //服务端请求广播当前状态
}

var FightCode = &fightCode{
	GRAB_LANDLORD_CREQ: "0",
	GRAB_LANDLORD_BRO:  "1",
	TURN_GRAB_BRO:      "2",
	DEAL_CREQ:          "3",
	DEAL_SRES:          "4",
	DEAL_BRO:           "5",
	PASS_CREQ:          "6",
	PASS_SRES:          "7",
	TURN_DEAL_BRO:      "8",
	LEAVE_BRO:          "9",
	OVER_BRO:           "10",
	GET_CARD_SRES:      "11",
	READY_CREQ:         "12",
	READY_BRO:          "13",
	STATUS_CREQ:        "14",
	STATUS_BRO:         "15",
}
