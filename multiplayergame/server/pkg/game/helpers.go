package game

import (
	"log"
	"math"
	"time"
)

type MessageToClient struct {
	MessageType string      `json:"type"`
	Payload     interface{} `json:"payload"`
}

func (g *Game) AddBullet(bullet *Bullet) {
	g.Bullets.Store(bullet.ID, bullet)
}

func (g *Game) DeleteBullet(bullet *Bullet) {
	g.Bullets.Delete(bullet.ID)
}

func Clamp(val float64, lo float64, hi float64) float64 {
	return math.Max(lo, math.Min(val, hi))
}

func (g *Game) UpdateBullets(dt float64) {
	g.Bullets.Range(func(_, ob interface{}) bool {
		b := ob.(*Bullet)

		b.update(dt)

		if b.Position.X > float64(WorldWidth) || b.Position.X < 0 || b.Position.Y > float64(WorldHeight) || b.Position.Y < 0 {
			g.DeleteBullet(b)
		}

		return true
	})
}

func (g *Game) WriteMessage(msgType string, player *Player) {

	var msg MessageToClient
	msg.MessageType = msgType

	switch msgType {
	case "start":
		var payload *Setting

		payload = g.getSetting()

		msg.Payload = payload
	case "state":
		var payload *State

		payload = g.getState(player)

		msg.Payload = payload
	case "end":
		var payload *EndGame

		payload = g.getEndGame()

		msg.Payload = payload
	}

	// err := player.Conn.WriteJSON(msg)
	// if err != nil {
	// 	log.Println(err)
	// }

}

func (g *Game) getSetting() *Setting {
	var setting Setting

	setting.GameWidth = WorldWidth
	setting.GameHeight = WorldHeight

	return &setting
}

func (g *Game) getState(player *Player) *State {
	var state State

	state.Timestamp = time.Now().UnixMilli()
	state.Player = player

	state.OtherPlayers = make([]*Player, 0)
	g.Players.Range(func(_, op interface{}) bool {
		p := op.(*Player)

		state.OtherPlayers = append(state.OtherPlayers, p)

		return true
	})

	state.Bullets = make([]*Bullet, 0)
	g.Bullets.Range(func(_, ob interface{}) bool {
		b := ob.(*Bullet)

		state.Bullets = append(state.Bullets, b)

		return true
	})

	playerLength := 0
	g.Players.Range(func(_, _ interface{}) bool {
		playerLength++

		return true
	})
	log.Println("players ", playerLength)

	return &state
}

func (g *Game) getEndGame() *EndGame {
	var endGame EndGame

	endGame.Data = "Отчушпанен"

	return &endGame
}
