package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/condemo/nes-cards-backend/config"
	"github.com/condemo/nes-cards-backend/store"
)

type PlayerHandler struct {
	store store.Store
}

func NewPlayerHandler(s store.Store) *PlayerHandler {
	return &PlayerHandler{
		store: s,
	}
}

func (h *PlayerHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /", makeHandler(h.getPlayerList))
	r.HandleFunc("POST /", makeHandler(h.createPlayer))
}

func (h *PlayerHandler) getPlayerList(w http.ResponseWriter, r *http.Request) error {
	var limit int
	l := r.FormValue("limit")
	if l != "" {
		limitInt, err := strconv.ParseInt(l, 10, 64)
		if err != nil {
			return err
		}
		limit = int(limitInt)
	} else {
		limit = config.ServerConfig.PlayerLimit
	}

	pl, err := h.store.GetPlayerList(int(limit))
	if err != nil {
		return err
	}
	fmt.Println(limit)

	return SendJSON(w, http.StatusOK, pl)
}

func (h *PlayerHandler) createPlayer(w http.ResponseWriter, r *http.Request) error {
	// TODO:
	return nil
}
