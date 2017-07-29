package room

type UserInfo struct {
	UserId     int
	PlayerInfo string
	Connected  bool
	Ready      bool
}

func NewUserInfo(uid int) UserInfo {
	u := UserInfo{}
	u.UserId = uid
	u.PlayerInfo = ""
	u.Connected = true
	u.Ready = false
	return u
}
