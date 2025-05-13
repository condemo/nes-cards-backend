package handlers

import (
	"net/http"
	"strconv"

	"github.com/condemo/nes-cards-backend/service"
	"github.com/condemo/nes-cards-backend/store"
)

type CurrentGameHandler struct {
	gs    *service.GameService
	store store.Store
}

func NewCurrentGameHandlder(gs *service.GameService, s store.Store) *CurrentGameHandler {
	return &CurrentGameHandler{
		gs:    gs,
		store: s,
	}
}

func (h *CurrentGameHandler) RegisterRoutes(r *http.ServeMux) {
	// NOTE: Igual no hace falta esta ruta pero no viene mal, reevaluar de vez en cuando
	r.HandleFunc("POST /set/{id}", MakeHandler(h.setGame))
}

func (h *CurrentGameHandler) setGame(w http.ResponseWriter, r *http.Request) error {
	gID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return err
	}
	// TODO: checks...

	game, err := h.store.GetGameById(gID)
	if err != nil {
		return err
	}

	h.gs.SetGame(game)
	return nil
}
