package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/condemo/nes-cards-backend/config"
	"github.com/condemo/nes-cards-backend/store"
	"github.com/condemo/nes-cards-backend/types"
)

type GameHandler struct {
	store store.Store
}

func NewGameHandler(s store.Store) *GameHandler {
	return &GameHandler{
		store: s,
	}
}

func (h *GameHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /", makeHandler(h.getGameRecords))
	r.HandleFunc("GET /{id}", makeHandler(h.getGame))
	r.HandleFunc("POST /", makeHandler(h.createGame))
}

func (h *GameHandler) getGame(w http.ResponseWriter, r *http.Request) error {
	// TODO: Mejorar implementación
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return err
	}

	g, err := h.store.GetGameById(id)
	if err != nil {
		return err
	}

	return SendJSON(w, http.StatusOK, g)
}

func (h *GameHandler) createGame(w http.ResponseWriter, r *http.Request) error {
	g := new(types.GameSetup)
	err := json.NewDecoder(r.Body).Decode(g)
	if err != nil {
		return err
	}

	// Stats
	if g.TowerHP == 0 {
		g.TowerHP = 60
	}
	if g.PlayerHP == 0 {
		g.PlayerHP = 80
	}

	// Get Player from DB
	if err := h.store.GetPlayerByName(g.Player1); err != nil {
		log.Fatal(err)
		return nil
	}
	if err := h.store.GetPlayerByName(g.Player2); err != nil {
		log.Fatal(err)
		return nil
	}

	// CreateGame
	game := types.NewGame(g)
	if err := h.store.CreateGame(game); err != nil {
		fmt.Println(err)
		return nil
	}

	p1Stats := types.NewStats(game.ID, game.Player1.ID, g.PlayerHP, g.TowerHP)
	p2Stats := types.NewStats(game.ID, game.Player2.ID, g.PlayerHP, g.TowerHP)

	tl := []*types.Stats{p1Stats, p2Stats}
	for _, t := range tl {
		if err := t.Validate(); err != nil {
			return NewApiError(
				err,
				"towers hp must be more than 0 and less than 255",
				http.StatusBadRequest)
		}
	}

	if err := h.store.CreatePlayerStats(tl); err != nil {
		return err
	}

	game, err = h.store.GetLastGame()
	if err != nil {
		return err
	}

	return SendJSON(w, http.StatusCreated, game)
}

func (h *GameHandler) getGameRecords(w http.ResponseWriter, r *http.Request) error {
	var limit int
	l := r.FormValue("limit")
	if l != "" {
		limitInt64, err := strconv.ParseInt(l, 10, 64)
		if err != nil {
			return err
		}
		limit = int(limitInt64)
	} else {
		limit = config.ServerConfig.GameRecordLimit
	}

	gl, err := h.store.GetGameList(limit)
	if err != nil {
		return err
	}

	return SendJSON(w, http.StatusOK, gl)
}
