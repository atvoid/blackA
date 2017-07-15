package cards

type BlackASuite struct {
	CardSuite
}

func generateCardsForBlackA() []Card {
	cards := [...]Card {
		{ CardType: CARDTYPE_CLUB, CardNumber: 1},
		{ CardType: CARDTYPE_CLUB, CardNumber: 2},
		{ CardType: CARDTYPE_CLUB, CardNumber: 5},
		{ CardType: CARDTYPE_CLUB, CardNumber: 6},
		{ CardType: CARDTYPE_CLUB, CardNumber: 7},
		{ CardType: CARDTYPE_CLUB, CardNumber: 8},
		{ CardType: CARDTYPE_CLUB, CardNumber: 9},
		{ CardType: CARDTYPE_CLUB, CardNumber: 10},
		{ CardType: CARDTYPE_CLUB, CardNumber: 11},
		{ CardType: CARDTYPE_CLUB, CardNumber: 12},
		{ CardType: CARDTYPE_CLUB, CardNumber: 13},
		{ CardType: CARDTYPE_DIAMOND, CardNumber: 1},
		{ CardType: CARDTYPE_DIAMOND, CardNumber: 2},
		{ CardType: CARDTYPE_DIAMOND, CardNumber: 5},
		{ CardType: CARDTYPE_DIAMOND, CardNumber: 6},
		{ CardType: CARDTYPE_DIAMOND, CardNumber: 7},
		{ CardType: CARDTYPE_DIAMOND, CardNumber: 8},
		{ CardType: CARDTYPE_DIAMOND, CardNumber: 9},
		{ CardType: CARDTYPE_DIAMOND, CardNumber: 10},
		{ CardType: CARDTYPE_DIAMOND, CardNumber: 11},
		{ CardType: CARDTYPE_DIAMOND, CardNumber: 12},
		{ CardType: CARDTYPE_DIAMOND, CardNumber: 13},
		{ CardType: CARDTYPE_HEART, CardNumber: 1},
		{ CardType: CARDTYPE_HEART, CardNumber: 2},
		{ CardType: CARDTYPE_HEART, CardNumber: 5},
		{ CardType: CARDTYPE_HEART, CardNumber: 6},
		{ CardType: CARDTYPE_HEART, CardNumber: 7},
		{ CardType: CARDTYPE_HEART, CardNumber: 8},
		{ CardType: CARDTYPE_HEART, CardNumber: 9},
		{ CardType: CARDTYPE_HEART, CardNumber: 10},
		{ CardType: CARDTYPE_HEART, CardNumber: 11},
		{ CardType: CARDTYPE_HEART, CardNumber: 12},
		{ CardType: CARDTYPE_HEART, CardNumber: 13},
		{ CardType: CARDTYPE_SPADE, CardNumber: 1},
		{ CardType: CARDTYPE_SPADE, CardNumber: 2},
		{ CardType: CARDTYPE_SPADE, CardNumber: 5},
		{ CardType: CARDTYPE_SPADE, CardNumber: 6},
		{ CardType: CARDTYPE_SPADE, CardNumber: 7},
		{ CardType: CARDTYPE_SPADE, CardNumber: 8},
		{ CardType: CARDTYPE_SPADE, CardNumber: 9},
		{ CardType: CARDTYPE_SPADE, CardNumber: 10},
		{ CardType: CARDTYPE_SPADE, CardNumber: 11},
		{ CardType: CARDTYPE_SPADE, CardNumber: 12},
		{ CardType: CARDTYPE_SPADE, CardNumber: 13},
		{ CardType: CARDTYPE_JOKER_S},
		{ CardType: CARDTYPE_JOKER_L},
	}
	return cards[:];
}

func CreateCardSuiteForBlackA() BlackASuite {
	return BlackASuite{ CreateCardSuite(1, generateCardsForBlackA) }
}

func (suite *BlackASuite)ToString() string {
	return "BlackASuite"
}

func (suite *BlackASuite)Compare(a *Card, b *Card) int {
	switch {
		case a.CardType >= CARDTYPE_JOKER_S && b.CardType >= CARDTYPE_JOKER_S:
			return (int)(a.CardType - b.CardType)
	}
	return 0
}