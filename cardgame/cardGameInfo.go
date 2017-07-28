package cardgame

import (
	"encoding/json"
)

type CardGameInfo struct {
	IsEnd bool
	Turn  int
	Wind  bool
}

func (this *CardGameInfo) ToMessage() string {
	b, _ := json.Marshal(*this)
	return string(b)
}
