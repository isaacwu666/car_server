package dtofight

import "github.com/gogf/gf/v2/util/grand"

// Poker 扑克poker
type Poker struct {
	Id     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Color  int    `json:"color,omitempty"`
	Weight int    `json:"weight,omitempty"` //权值
}

type pokerColor struct {
	CLUB   int `json:"CLUB"`   // 1 梅花
	HEART  int `json:"HEART"`  // 2 红桃
	SPADE  int `json:"SPADE"`  // 3 黑桃
	SQUARE int `json:"SQUARE"` // 4 方片
	JOKER  int `json:"JOKER"`  // 5 王
}

func (c *pokerColor) GetColorStr(color int) string {
	switch color {
	case c.CLUB:
		return "Club"
	case c.HEART:
		return "Heart"
	case c.SPADE:
		return "Spade"
	case c.SQUARE:
		return "Square"
	case c.JOKER:
		return "Joker"
	}
	return ""
}

type pokerWeight struct {
	THREE  int `json:"THREE,omitempty"`  //3
	FOUR   int `json:"FOUR,omitempty"`   //4
	FIVE   int `json:"FIVE,omitempty"`   //5
	SIX    int `json:"SIX,omitempty"`    //6
	SEVEN  int `json:"SEVEN,omitempty"`  //7
	EIGHT  int `json:"EIGHT,omitempty"`  //8
	NINE   int `json:"NINE,omitempty"`   //9
	TEN    int `json:"TEN,omitempty"`    //10
	JACK   int `json:"JACK,omitempty"`   //11
	QUEEN  int `json:"QUEEN,omitempty"`  //12
	KING   int `json:"KING,omitempty"`   //13
	ONE    int `json:"ONE,omitempty"`    //14
	TWO    int `json:"TWO,omitempty"`    //15
	SJOKER int `json:"SJOKER,omitempty"` //16
	LJOKER int `json:"LJOKER,omitempty"` //17
}

func (c *pokerWeight) GetWeightStr(wight int) string {
	switch wight {
	case c.THREE:
		return "Three"
	case c.FOUR:
		return "Four"
	case c.FIVE:
		return "Five"
	case c.SIX:
		return "Six"
	case c.SEVEN:
		return "Seven"
	case c.EIGHT:
		return "Eight"
	case c.NINE:
		return "Nine"
	case c.TEN:
		return "Ten"
	case c.JACK:
		return "Jack"
	case c.QUEEN:
		return "Queen"
	case c.KING:
		return "King"
	case c.ONE:
		return "One"
	case c.TWO:
		return "Two"
	case c.SJOKER:
		return "Sjoker"
	case c.LJOKER:
		return "Ljoker"
	}
	return ""
}

var (
	PokerColor = pokerColor{
		CLUB:   1,
		HEART:  2,
		SPADE:  3,
		SQUARE: 4,
		JOKER:  5,
	}
	PokerWeight = pokerWeight{
		THREE:  3,
		FOUR:   4,
		FIVE:   5,
		SIX:    6,
		SEVEN:  7,
		EIGHT:  8,
		NINE:   9,
		TEN:    10,
		JACK:   11,
		QUEEN:  12,
		KING:   13,
		ONE:    14,
		TWO:    15,
		SJOKER: 16, //小王
		LJOKER: 17, //大王
	}
)

// NewPokerArray 创建牌
func NewPokerArray() []*Poker {
	var idx = 0
	var pokerArray []*Poker = make([]*Poker, 54)
	//外层先循环花色
	for Color := PokerColor.CLUB; Color < PokerColor.JOKER; Color++ {
		cStr := PokerColor.GetColorStr(Color)
		//循环，权值
		for Weight := PokerWeight.THREE; Weight < PokerWeight.SJOKER; Weight++ {

			pokerArray[idx] = &Poker{
				Id:     idx,
				Name:   cStr + PokerWeight.GetWeightStr(Weight),
				Weight: Weight,
				Color:  Color,
			}
			idx++
		}
	}
	pokerArray[idx] = &Poker{
		Id: idx,
		Name: PokerColor.GetColorStr(PokerColor.JOKER) +
			PokerWeight.GetWeightStr(PokerWeight.SJOKER),
		Color:  0,
		Weight: 0,
	}
	idx++
	pokerArray[idx] = &Poker{
		Id: idx,
		Name: PokerColor.GetColorStr(PokerColor.JOKER) +
			PokerWeight.GetWeightStr(PokerWeight.LJOKER),
		Color:  0,
		Weight: 0,
	}
	return pokerArray
}

// ShufflePoker 洗牌
func ShufflePoker(array []*Poker) {

	for i := len(array) - 1; i > 0; i-- {
		ran := grand.N(1, i)
		tmp := array[i]
		array[i] = array[ran]
		array[ran] = tmp
	}
}
