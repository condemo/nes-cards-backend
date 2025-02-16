package types

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Game struct {
	bun.BaseModel `bun:"table:games,alias:g"`

	ID        int64     `bun:",pk,autoincrement" json:"id"`
	P1ID      int64     `json:"p1id"`
	P2ID      int64     `json:"p2id"`
	Player1   *Player   `bun:"rel:belongs-to,join:p1id=id" json:"player1"`
	Player2   *Player   `bun:"rel:belongs-to,join:p2id=id" json:"player2"`
	P1Stats   *Stats    `bun:"rel:has-one,join:p1id=player_id,join:id=game_id" json:"p1stats"`
	P2Stats   *Stats    `bun:"rel:has-one,join:p2id=player_id,join:id=game_id" json:"p2stats"`
	Winner    string    `bun:"winner" json:"winner"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"createdAt"`
}

var _ bun.BeforeAppendModelHook = (*Game)(nil)

// BeforeAppendModel gets the current time in GMT+1 and set it in CreatedAt field
func (g *Game) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		location, err := time.LoadLocation("Europe/Madrid")
		if err != nil {
			return err
		} // FIX: Añado una hora de free por el cambio de hora, debería ser automático
		g.CreatedAt = time.Now().In(location).Add(time.Hour)
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
	}

	return g
}

type GameSetup struct {
	Player1  *Player `json:"player1"`
	Player2  *Player `json:"player2"`
	PlayerHP uint8   `json:"playerHP"`
	TowerHP  uint8   `json:"towerHP"`
}
