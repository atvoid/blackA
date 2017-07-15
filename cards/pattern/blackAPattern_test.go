package pattern

import (
	"blackA/cards"
	"testing"
)

func makeCards(cardType int8, list string) []cards.Card {
	ans := make([]cards.Card, len(list))
	for i, ll := 0, len(list); i < ll; i++ {
		var num byte = 0
		switch {
			case list[i] >= '0' && list[i] <= '9':
				num = list[i] - byte('0')
			case list[i] == 'X':
				num = 10
			case list[i] == 'J':
				num = 11
			case list[i] == 'Q':
				num = 12
			case list[i] == 'K':
				num = 13
		}
		ans[i] = cards.Card{ CardType: cardType, CardNumber: int(num) }
	}
	return ans
}

func Test_blackAPattern_Double(t *testing.T) {
	type args struct {
		list []cards.Card
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{ name: "11", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "1"), makeCards(cards.CARDTYPE_SPADE, "1")...)}, want: PATTERN_INVALID },
		{ name: "11", args: args{list: makeCards(cards.CARDTYPE_CLUB, "11")}, want: BLACKAPATTERN_DOUBLE },
		{ name: "12", args: args{list: makeCards(cards.CARDTYPE_CLUB, "12")}, want: PATTERN_INVALID },
		{ name: "111", args: args{list: makeCards(cards.CARDTYPE_CLUB, "12")}, want: PATTERN_INVALID },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := blackAPattern_Double(tt.args.list); got != tt.want {
				t.Errorf("blackAPattern_Double() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_blackAPattern_DoubleStraight(t *testing.T) {
	type args struct {
		list []cards.Card
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{ name: "JJQQKK11", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "JK1"), makeCards(cards.CARDTYPE_SPADE, "KJQ1Q")...)}, want: BLACKAPATTERN_DOUBLESTRAIGHT },
		{ name: "11223344", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "134"), makeCards(cards.CARDTYPE_SPADE, "21243")...)}, want: BLACKAPATTERN_DOUBLESTRAIGHT },
		{ name: "6677788899", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "677898"), makeCards(cards.CARDTYPE_SPADE, "9876")...)}, want: PATTERN_INVALID },
		{ name: "12345", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "23"), makeCards(cards.CARDTYPE_SPADE, "145")...)}, want: PATTERN_INVALID },
		{ name: "XJQK1", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "XKJ"), makeCards(cards.CARDTYPE_SPADE, "1Q")...)}, want: PATTERN_INVALID },
		{ name: "6789X", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "X6"), makeCards(cards.CARDTYPE_SPADE, "798")...)}, want: PATTERN_INVALID },
		{ name: "36789X", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "3X6"), makeCards(cards.CARDTYPE_SPADE, "798")...)}, want: PATTERN_INVALID },
		{ name: "66778899", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "6789"), makeCards(cards.CARDTYPE_SPADE, "9876")...)}, want: BLACKAPATTERN_DOUBLESTRAIGHT },
		{ name: "11", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "1"), makeCards(cards.CARDTYPE_DIAMOND, "1")...)}, want: PATTERN_INVALID },
		{ name: "11", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "1"), makeCards(cards.CARDTYPE_CLUB, "1")...)}, want: PATTERN_INVALID },
		{ name: "12", args: args{list: makeCards(cards.CARDTYPE_CLUB, "12")}, want: PATTERN_INVALID },
		{ name: "111", args: args{list: makeCards(cards.CARDTYPE_CLUB, "12")}, want: PATTERN_INVALID },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := blackAPattern_DoubleStraight(tt.args.list); got != tt.want {
				t.Errorf("blackAPattern_DoubleStraight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_blackAPattern_DoubleKing(t *testing.T) {
	type args struct {
		list []cards.Card
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{ name: "00", args: args{list: append(makeCards(cards.CARDTYPE_JOKER_L, "0"), makeCards(cards.CARDTYPE_JOKER_S, "0")...)}, want: BLACKAPATTERN_DOUBLEKING },
		{ name: "01", args: args{list: append(makeCards(cards.CARDTYPE_JOKER_L, "0"), makeCards(cards.CARDTYPE_CLUB, "1")...)}, want: PATTERN_INVALID },
		{ name: "12", args: args{list: makeCards(cards.CARDTYPE_CLUB, "12")}, want: PATTERN_INVALID },
		{ name: "111", args: args{list: makeCards(cards.CARDTYPE_CLUB, "12")}, want: PATTERN_INVALID },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := blackAPattern_DoubleKing(tt.args.list); got != tt.want {
				t.Errorf("blackAPattern_DoubleKing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_blackAPattern_DoubleBlackA(t *testing.T) {
	type args struct {
		list []cards.Card
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{ name: "11", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "1"), makeCards(cards.CARDTYPE_SPADE, "1")...)}, want: BLACKAPATTERN_DOUBLEBLACKA },
		{ name: "11", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "1"), makeCards(cards.CARDTYPE_DIAMOND, "1")...)}, want: PATTERN_INVALID },
		{ name: "11", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "1"), makeCards(cards.CARDTYPE_CLUB, "1")...)}, want: PATTERN_INVALID },
		{ name: "12", args: args{list: makeCards(cards.CARDTYPE_CLUB, "12")}, want: PATTERN_INVALID },
		{ name: "111", args: args{list: makeCards(cards.CARDTYPE_CLUB, "12")}, want: PATTERN_INVALID },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := blackAPattern_DoubleBlackA(tt.args.list); got != tt.want {
				t.Errorf("blackAPattern_DoubleBlackA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_blackAPattern_Straight(t *testing.T) {
	type args struct {
		list []cards.Card
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{ name: "12345", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "23"), makeCards(cards.CARDTYPE_SPADE, "145")...)}, want: BLACKAPATTERN_STRAIGHT },
		{ name: "XJQK1", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "XKJ"), makeCards(cards.CARDTYPE_SPADE, "1Q")...)}, want: BLACKAPATTERN_STRAIGHT },
		{ name: "6789X", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "X6"), makeCards(cards.CARDTYPE_SPADE, "798")...)}, want: BLACKAPATTERN_STRAIGHT },
		{ name: "36789X", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "3X6"), makeCards(cards.CARDTYPE_SPADE, "798")...)}, want: PATTERN_INVALID },
		{ name: "66778899", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "6789"), makeCards(cards.CARDTYPE_SPADE, "9876")...)}, want: PATTERN_INVALID },
		{ name: "11", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "1"), makeCards(cards.CARDTYPE_DIAMOND, "1")...)}, want: PATTERN_INVALID },
		{ name: "11", args: args{list: append(makeCards(cards.CARDTYPE_CLUB, "1"), makeCards(cards.CARDTYPE_CLUB, "1")...)}, want: PATTERN_INVALID },
		{ name: "12", args: args{list: makeCards(cards.CARDTYPE_CLUB, "12")}, want: PATTERN_INVALID },
		{ name: "111", args: args{list: makeCards(cards.CARDTYPE_CLUB, "12")}, want: PATTERN_INVALID },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := blackAPattern_Straight(tt.args.list); got != tt.want {
				t.Errorf("blackAPattern_Straight() = %v, want %v", got, tt.want)
			}
		})
	}
}
