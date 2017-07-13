package cards

import (
	"math/rand"
	"time"
)

type CardSuite struct {
	CardList	[]Card
}

func CreateCardSuite(suiteCount int, cardsGenerator func () []Card) CardSuite {
	cards := cardsGenerator()
	for i := 0; i < suiteCount; i++ {
		cards = append(cards, cardsGenerator()...)
	}
	return CardSuite{ CardList: cards }
}

func (suite *CardSuite)Shuffle() {
	rand.Seed(time.Now().UnixNano())  
	for i := range suite.CardList {
		j := rand.Intn(i+1)
		suite.CardList[i], suite.CardList[j] = suite.CardList[j], suite.CardList[i]
	}
}

func (suite *CardSuite)ToString() string {
	return "CardSuite"
}