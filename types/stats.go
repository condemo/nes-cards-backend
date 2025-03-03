package types

import "github.com/uptrace/bun"

type Stats struct {
	bun.BaseModel `bun:"table:stats,alias:t"`

	ID         int64 `bun:",pk,autoincrement" json:"id"`
	GameID     int64 `bun:",notnull" json:"gameID"`
	PlayerID   int64 `bun:",notnull" json:"playerID"`
	HP         uint8 `bun:",nullzero" validator:"gte=0,lte=255" json:"hp"`
	Strength   int16 `bun:",nullzero" json:"strength"`
	Intangible uint8 `bun:",nullzero" json:"intangible"`
	Confusion  uint8 `bun:",nullzero" json:"confusion"`
	T1HP       uint8 `bun:",nullzero" json:"t1hp"`
	T2HP       uint8 `bun:",nullzero" json:"t2hp"`
}

// NewStats recibe GameID y PlayerID además de una cantidad de HP e instancia dos Torres
func NewStats(gid, pid int64, php, thp uint8) *Stats {
	s := &Stats{
		GameID:   gid,
		PlayerID: pid,
		HP:       php,
		T1HP:     thp,
		T2HP:     thp,
	}

	return s
}

func (t *Stats) Validate() error {
	err := validate.Struct(t)
	return err
}
