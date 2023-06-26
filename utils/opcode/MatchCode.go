package opcode

type matchCode struct {

	//进入匹配队列
	//EnterCreq
	EnterCreq string "0"
	//EnterSres
	EnterSres string "1"
	//EnterBro
	EnterBro string "10"

	//离开匹配队列
	//LeaveCreq
	LeaveCreq string "2"
	//LeaveSres
	// LeaveSres string  "3";
	//LeaveBro
	LeaveBro string "3"

	//准备
	//ReadyCreq
	ReadyCreq string "4"
	//ReadySres
	// ReadySres string  "5";
	//ReadyBro
	ReadyBro string "5"

	//开始游戏
	//StartCreq
	// StartCreq string  "6";
	//StartSres
	// StartSres string  "7";
	//StartBro
	StartBro string "6"
}

var MatchCode = matchCode{

	EnterCreq: "0",
	EnterSres: "1",
	EnterBro:  "10",
	LeaveCreq: "2",
	//LeaveSres :"3",
	LeaveBro:  "3",
	ReadyCreq: "4",
	//ReadySres :"5",
	ReadyBro: "5",
	//StartCreq :"6",
	//StartSres :"7",
	StartBro: "6",
}
