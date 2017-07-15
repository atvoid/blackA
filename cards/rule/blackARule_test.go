package rule

import (
	"blackA/cards"
	"blackA/cards/pattern"
	"testing"
	"fmt"
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

func TestBlackAGame_DealCards(t *testing.T) {
	game := ICardRule(&BlackAGame{})
	game.Init()
	result := game.DealCards()
	if (len(result) != 4) {
		t.Errorf("BlackAGame.DealCards() result length should be %v, but is %v", 4, len(result))
	}

	for i, v := range result {
		fmt.Printf("Player %v:", i)
		t.Logf("Player %v:", i)
		for _, c := range v {
			fmt.Print(c.ToString())
			t.Log(c.ToString())
		}
		fmt.Printf("\t\t%v\n", len(v))
		t.Log("\n")
	}
}

func TestBlackAGame_Compare(t *testing.T) {
	game := &BlackAGame{}
	matcher := pattern.GenerateBlackAPatternMatcher()
	type args struct {
		a pattern.CardPattern
		b pattern.CardPattern
	}
	tests := []struct {
		name  string
		this  *BlackAGame
		args  args
		want  int
		want1 bool
	}{
		{
			name: "11 33",
			this: game,
			args: args{
				a: matcher(append(makeCards(cards.CARDTYPE_DIAMOND, "1"), makeCards(cards.CARDTYPE_SPADE, "1")...)),
				b: matcher(append(makeCards(cards.CARDTYPE_CLUB, "3"), makeCards(cards.CARDTYPE_SPADE, "3")...))},
			want: -2,
			want1: true,
		},
		{
			name: "black11 33",
			this: game,
			args: args{
				a: matcher(append(makeCards(cards.CARDTYPE_CLUB, "1"), makeCards(cards.CARDTYPE_SPADE, "1")...)),
				b: matcher(append(makeCards(cards.CARDTYPE_CLUB, "3"), makeCards(cards.CARDTYPE_SPADE, "3")...))},
			want: 1,
			want1: true,
		},
		{
			name: "123 123",
			this: game,
			args: args{
				a: matcher(append(makeCards(cards.CARDTYPE_DIAMOND, "1"), makeCards(cards.CARDTYPE_SPADE, "23")...)),
				b: matcher(append(makeCards(cards.CARDTYPE_CLUB, "13"), makeCards(cards.CARDTYPE_SPADE, "2")...))},
			want: 0,
			want1: true,
		},
		{
			name: "123 333",
			this: game,
			args: args{
				a: matcher(append(makeCards(cards.CARDTYPE_DIAMOND, "1"), makeCards(cards.CARDTYPE_SPADE, "23")...)),
				b: matcher(append(makeCards(cards.CARDTYPE_CLUB, "33"), makeCards(cards.CARDTYPE_SPADE, "3")...))},
			want: -1,
			want1: true,
		},
		{
			name: "123 1234",
			this: game,
			args: args{
				a: matcher(append(makeCards(cards.CARDTYPE_DIAMOND, "1"), makeCards(cards.CARDTYPE_SPADE, "23")...)),
				b: matcher(append(makeCards(cards.CARDTYPE_CLUB, "24"), makeCards(cards.CARDTYPE_SPADE, "13")...))},
			want: 0,
			want1: false,
		},
		{
			name: "1234 5678",
			this: game,
			args: args{
				a: matcher(append(makeCards(cards.CARDTYPE_DIAMOND, "243"), makeCards(cards.CARDTYPE_SPADE, "1")...)),
				b: matcher(append(makeCards(cards.CARDTYPE_CLUB, "58"), makeCards(cards.CARDTYPE_SPADE, "67")...))},
			want: -1,
			want1: true,
		},
		{
			name: "QKA XJQ",
			this: game,
			args: args{
				a: matcher(append(makeCards(cards.CARDTYPE_DIAMOND, "1"), makeCards(cards.CARDTYPE_SPADE, "KQ")...)),
				b: matcher(append(makeCards(cards.CARDTYPE_CLUB, "X"), makeCards(cards.CARDTYPE_SPADE, "JQ")...))},
			want: 1,
			want1: true,
		},
		{
			name: "JJQQKK QQKKAA",
			this: game,
			args: args{
				a: matcher(append(makeCards(cards.CARDTYPE_DIAMOND, "QQKK"), makeCards(cards.CARDTYPE_SPADE, "JJ")...)),
				b: matcher(append(makeCards(cards.CARDTYPE_CLUB, "QK1"), makeCards(cards.CARDTYPE_SPADE, "QK1")...))},
			want: -1,
			want1: true,
		},
		{
			name: "JJQQKK 778899",
			this: game,
			args: args{
				a: matcher(append(makeCards(cards.CARDTYPE_DIAMOND, "QQKK"), makeCards(cards.CARDTYPE_SPADE, "JJ")...)),
				b: matcher(append(makeCards(cards.CARDTYPE_CLUB, "789"), makeCards(cards.CARDTYPE_SPADE, "879")...))},
			want: 4,
			want1: true,
		},
		{
			name: "444 5555",
			this: game,
			args: args{
				a: matcher(append(makeCards(cards.CARDTYPE_DIAMOND, "44"), makeCards(cards.CARDTYPE_SPADE, "4")...)),
				b: matcher(append(makeCards(cards.CARDTYPE_CLUB, "5"), makeCards(cards.CARDTYPE_SPADE, "555")...))},
			want: -1,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.this.Compare(&tt.args.a, &tt.args.b)
			if got != tt.want {
				t.Errorf("BlackAGame.Compare() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BlackAGame.Compare() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
