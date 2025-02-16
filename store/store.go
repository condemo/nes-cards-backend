package store

import (
	"context"

	"github.com/condemo/nes-cards-backend/types"
	"github.com/uptrace/bun"
)

type Store interface {
	// TODO: Añadir Updates y Deletes
	CreatePlayer(*types.Player) error
	CheckPlayer(string) bool
	GetPlayerById(*types.Player) error
	GetPlayerByName(*types.Player) error
	GetPlayerList(limit int) ([]types.Player, error)
	CreatePlayerStats([]*types.Stats) error
	CreateGame(*types.Game) error
	GetGameById(id int64) (*types.Game, error)
	GetLastGame() (*types.Game, error)
	GetGameList(limit int) ([]*types.Game, error)
}

type Storage struct {
	db *bun.DB
}

func NewStorage(db *bun.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreatePlayer(p *types.Player) error {
	_, err := s.db.NewInsert().Model(p).
		Returning("*").Exec(context.Background())
	return err
}

func (s *Storage) CheckPlayer(name string) bool {
	if err := s.db.NewSelect().
		Model(&types.Player{}).
		Where("name = ?", name).
		Scan(context.Background()); err != nil {
		return false
	}

	return true
}

func (s *Storage) GetPlayerById(p *types.Player) error {
	err := s.db.NewSelect().Model(p).
		Where("id = ?", p.ID).Scan(context.Background())

	return err
}

func (s *Storage) GetPlayerByName(p *types.Player) error {
	err := s.db.NewSelect().Model(p).
		Where("name = ?", p.Name).Scan(context.Background())

	return err
}

func (s *Storage) GetPlayerList(limit int) ([]types.Player, error) {
	var pl []types.Player
	err := s.db.NewSelect().Model(&pl).
		Order("id ASC").Limit(20).Limit(limit).
		Scan(context.Background())

	return pl, err
}

func (s *Storage) CreatePlayerStats(ps []*types.Stats) error {
	_, err := s.db.NewInsert().Model(&ps).
		Returning("*").Exec(context.Background())
	return err
}

func (s *Storage) CreateGame(g *types.Game) error {
	// TODO: Ineficiente, dos querys en lugar de una
	_, err := s.db.NewInsert().Model(g).
		Returning("*").Exec(context.Background())
	if err != nil {
		return err
	}

	return err
}

func (s *Storage) GetLastGame() (*types.Game, error) {
	g := new(types.Game)
	err := s.db.NewSelect().Model(g).
		Relation("Player1").Where("p1id = player1.id").
		Relation("Player2").Where("p2id = player2.id").
		Relation("P1Stats").Where("p1id = p1_stats.player_id").
		Relation("P2Stats").Where("p2id = p2_stats.player_id").
		Order("g.created_at DESC").Limit(1).
		Scan(context.Background())

	return g, err
}

func (s *Storage) GetGameList(limit int) ([]*types.Game, error) {
	var pl []*types.Game

	err := s.db.NewSelect().Model(&pl).
		Relation("Player1").Where("p1id=player1.id").
		Relation("Player2").Where("p2id=Player2.id").
		Relation("P1Stats").Where("p1id = p1_stats.player_id").
		Relation("P2Stats").Where("p2id = p2_stats.player_id").
		Order("g.created_at DESC").Limit(limit).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return pl, nil
}

func (s *Storage) GetGameById(id int64) (*types.Game, error) {
	g := new(types.Game)
	err := s.db.NewSelect().Model(g).
		Relation("Player1").Where("p1id = player1.id").
		Relation("Player2").Where("p2id = player2.id").
		Relation("P1Stats").Where("p1id = p1_stats.player_id").
		Relation("P2Stats").Where("p2id = p2_stats.player_id").
		Where("g.id = ?", id).Scan(context.Background())

	return g, err
}
