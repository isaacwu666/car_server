package opcode

type chatCode struct {
	CREQ string "0" //聊天
	SRES string "1"
}

var ChatCode = &chatCode{
	CREQ: "0",
	SRES: "1",
}
