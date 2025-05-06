package types

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Game struct {
	bun.BaseModel `bun:"table:games,alias:g"`

	ID         int64     `bun:",pk,autoincrement" json:"id"`
	P1ID       int64     `json:"p1id"`
	P2ID       int64     `json:"p2id"`
	Player1    *Player   `bun:"rel:belongs-to,join:p1id=id" json:"player1"`
	Player2    *Player   `bun:"rel:belongs-to,join:p2id=id" json:"player2"`
	P1Stats    *Stats    `bun:"rel:has-one,join:p1id=player_id,join:id=game_id" json:"p1stats"`
	P2Stats    *Stats    `bun:"rel:has-one,join:p2id=player_id,join:id=game_id" json:"p2stats"`
	Winner     string    `bun:"winner" json:"winner"`
	Round      uint16    `bun:"round,notnull" json:"round"`
	TurnMode   uint8     `bun:",nullzero" json:"turnMode"`
	PlayerTurn uint8     `bun:",nullzero" json:"playerTurn"`
	CreatedAt  time.Time `bun:",nullzero,notnull" json:"createdAt"`
	UpdateAt   time.Time `bun:"updateAt" json:"updateAt"`
}

var _ bun.BeforeAppendModelHook = (*Game)(nil)

// BeforeAppendModel gets the current time in Madrid and set it in CreatedAt field
func (g *Game) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	location, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		return err
	}

	switch query.(type) {
	case *bun.InsertQuery:
		g.CreatedAt = time.Now().In(location)
	case *bun.UpdateQuery:
		g.UpdateAt = time.Now().In(location)
	}

	return nil
}

func NewGame(gs *GameSetup) *Game {
	g := &Game{
		P1ID:    gs.Player1.ID,
		P2ID:    gs.Player2.ID,
		Winner:  "none",
		Player1: gs.Player1,
		Player2: gs.Player2,
		Round:   1,
	}

	return g
}

type GameSetup struct {
	Player1  *Player `json:"player1"`
	Player2  *Player `json:"player2"`
	PlayerHP uint8   `json:"playerHP"`
	TowerHP  uint8   `json:"towerHP"`
}
