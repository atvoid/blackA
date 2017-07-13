package cards

import "fmt"

const (
	CARDTYPE_DIAMOND = 0
	CARDTYPE_CLUB = 1
	CARDTYPE_HEART = 2
	CARDTYPE_SPADE = 3
	CARDTYPE_JOKER_S = 4
	CARDTYPE_JOKER_L = 5
)

type Card struct {
	CardType	int8
	CardNumber	int32
}

func (card *Card) ToString() string {
	var num string
	switch {
		case card.CardNumber >= 2 && card.CardNumber <= 9:
			num = fmt.Sprintf("%d", card.CardNumber)
		case card.CardNumber == 1:
			num = "A"
		case card.CardNumber == 10:
			num = "0"
		case card.CardNumber == 11:
			num = "J"
		case card.CardNumber == 12:
			num = "Q"
		case card.CardNumber == 13:
			num = "K"
	}
	switch card.CardType {
		case CARDTYPE_SPADE: return "♠" + num
		case CARDTYPE_CLUB: return "♣" + num
		case CARDTYPE_HEART: return "♥" + num
		case CARDTYPE_DIAMOND: return "◆" + num
		case CARDTYPE_JOKER_S: return "🃟"
		case CARDTYPE_JOKER_L: return "🃏"
	}
	panic("Wrong Card")
}
