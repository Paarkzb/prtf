package game

type State struct {
	Timestamp    int64     `json:"timestamp"`
	Player       *Player   `json:"player"`
	OtherPlayers []*Player `json:"otherPlayers"`
	Bullets      []*Bullet `json:"bullets"`
}

type Setting struct {
	GameWidth  int32 `json:"game_width"`
	GameHeight int32 `json:"game_height"`
}

type EndGame struct {
	Data string `json:"data"`
}
