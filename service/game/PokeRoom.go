package game

import (
	"github.com/gogf/gf/v2/os/gtime"
)

type happyPokerRoom struct {
	running  bool
	rootName string "欢乐斗地主"
	rootId   int64
	player1  *PokeRoomContext
	player2  *PokeRoomContext
	player3  *PokeRoomContext
	pCount   int
}

func (r *happyPokerRoom) AddPlayer(player *PokeRoomContext) {
	if r.player1 == nil {
		r.player1 = player
		r.pCount = 1
	} else if r.player2 == nil {
		r.player2 = player
		r.pCount = 2
	} else if r.player3 == nil {
		r.player3 = player
		r.pCount = 3
	}
}
func NewHappyPokerRoom() *happyPokerRoom {
	i := new(happyPokerRoom)

	i.running = true
	i.rootId = gtime.Now().Timestamp()
	//开启房间线程
	go i.startGame()
	return i
}

func (r *happyPokerRoom) startGame() {

	for r.running {

	}
}
