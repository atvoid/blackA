package main

import (
	"blackA/cards"
	"fmt"
)

func main() {
	suite := cards.CreateCardSuiteForBlackA()
	suite.Shuffle()
	for _, v := range suite.CardList {
		fmt.Print(v.ToString())
	}
}
