package cards

type CardList []Card;

func (this CardList) Len() int {
	return len(this)
}

func (this CardList) Less(i, j int) bool {
	if (this[i].CardNumber != this[j].CardNumber) {
		return this[i].CardNumber < this[j].CardNumber
	} else {
		return this[i].CardType < this[i].CardType
	}
}

func (this CardList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}