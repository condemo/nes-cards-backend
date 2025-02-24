package service

import "github.com/condemo/nes-cards-backend/types"

type GameService struct {
	game *types.Game
}

func NewGameService() *GameService {
	return &GameService{}
}

func (gs *GameService) SetGame(g *types.Game) {
	gs.game = g
}
