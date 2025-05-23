package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/condemo/nes-cards-backend/config"
	"github.com/condemo/nes-cards-backend/service"
	"github.com/condemo/nes-cards-backend/store"
	"github.com/condemo/nes-cards-backend/types"
)

type GameHandler struct {
	store store.Store
	gs    *service.GameService
}

func NewGameHandler(s store.Store, gs *service.GameService) *GameHandler {
	return &GameHandler{
		store: s,
		gs:    gs,
	}
}

func (h *GameHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /", MakeHandler(h.getGameRecords))
	r.HandleFunc("GET /{id}", MakeHandler(h.getGame))
	r.HandleFunc("GET /last", MakeHandler(h.getLastGame))
	r.HandleFunc("POST /", MakeHandler(h.createGame))
	r.HandleFunc("PUT /", MakeHandler(h.updateGame))
	r.HandleFunc("DELETE /{id}", MakeHandler(h.deleteGame))
	r.HandleFunc("PUT /stats", MakeHandler(h.updateStats))
}

func (h *GameHandler) getGame(w http.ResponseWriter, r *http.Request) error {
	var updateCurrent bool
	uc := r.URL.Query().Get("updateCurrent")
	if uc != "" {
		ok, err := strconv.ParseBool(uc)
		if err != nil {
			return NewApiError(err, "updateCurrent must be a boolean value", http.StatusBadRequest)
		}
		updateCurrent = ok
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return err
	}

	game, err := h.store.GetGameById(id)
	if err != nil {
		return err
	}

	if updateCurrent {
		h.gs.SetGame(game)
	}

	return SendJSON(w, http.StatusOK, game)
}

func (h *GameHandler) getLastGame(w http.ResponseWriter, r *http.Request) error {
	var updateCurrent bool

	cu := r.FormValue("updateCurrent")
	if cu != "" {
		ok, err := strconv.ParseBool(cu)
		if err != nil {
			return NewApiError(err, "updateCurrent must be a boolean value", http.StatusBadRequest)
		}
		updateCurrent = ok
	}

	game, err := h.store.GetLastGame()
	if err != nil {
		return NewApiError(err, "there is no game yet", http.StatusNotFound)
	}

	if updateCurrent && game != nil {
		h.gs.SetGame(game)
	}

	return SendJSON(w, http.StatusOK, game)
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
		return err
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

	h.gs.SetGame(game)

	return SendJSON(w, http.StatusCreated, game)
}

func (h *GameHandler) getGameRecords(w http.ResponseWriter, r *http.Request) error {
	var limit int
	l := r.URL.Query().Get("limit")
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

func (h *GameHandler) updateGame(w http.ResponseWriter, r *http.Request) error {
	g := new(types.Game)

	if err := json.NewDecoder(r.Body).Decode(g); err != nil {
		return err
	}

	if err := h.store.UpdateGame(g); err != nil {
		return err
	}

	return SendJSON(w, http.StatusOK, g)
}

func (h *GameHandler) deleteGame(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return err
	}

	if _, err := h.store.GetGameById(id); err != nil {
		return NewApiError(err, "game not found", http.StatusNotFound)
	}

	if err := h.store.DeleteGame(id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *GameHandler) updateStats(w http.ResponseWriter, r *http.Request) error {
	st := new(types.Stats)

	if err := json.NewDecoder(r.Body).Decode(st); err != nil {
		return err
	}

	if err := h.store.UpdateStats(st); err != nil {
		return err
	}

	return SendJSON(w, http.StatusOK, st)
}
