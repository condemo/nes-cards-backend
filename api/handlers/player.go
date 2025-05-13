package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/condemo/nes-cards-backend/config"
	"github.com/condemo/nes-cards-backend/store"
	"github.com/condemo/nes-cards-backend/types"
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
	r.HandleFunc("GET /", MakeHandler(h.getPlayerList))
	r.HandleFunc("POST /", MakeHandler(h.createPlayer))
	r.HandleFunc("PUT /", MakeHandler(h.updatePlayer))
	r.HandleFunc("DELETE /{id}", MakeHandler(h.deletePlayer))
}

func (h *PlayerHandler) getPlayerList(w http.ResponseWriter, r *http.Request) error {
	var limit int
	l := r.URL.Query().Get("limit")
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

	return SendJSON(w, http.StatusOK, pl)
}

func (h *PlayerHandler) createPlayer(w http.ResponseWriter, r *http.Request) error {
	p := new(types.Player)
	json.NewDecoder(r.Body).Decode(p)

	if p.Name == "" {
		return NewApiError(
			errors.New("player name is empty"),
			"Player name is empty",
			http.StatusBadRequest)
	}

	if ok := h.store.CheckPlayer(p.Name); ok {
		return NewApiError(
			errors.New("player already exists"),
			fmt.Sprintf("%s already exists", p.Name),
			http.StatusConflict,
		)
	}

	if err := p.Validate(); err != nil {
		return NewApiError(err, "invalid player name", http.StatusBadRequest)
	}

	if err := h.store.CreatePlayer(p); err != nil {
		return err
	}

	return SendJSON(w, http.StatusCreated, p)
}

func (h *PlayerHandler) updatePlayer(w http.ResponseWriter, r *http.Request) error {
	p := new(types.Player)

	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		return err
	}
	if err := p.Validate(); err != nil {
		return NewApiError(err, "invalid player name", http.StatusBadRequest)
	}

	if err := h.store.UpdatePlayer(p); err != nil {
		return err
	}

	return SendJSON(w, http.StatusOK, p)
}

func (h *PlayerHandler) deletePlayer(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		return err
	}

	if _, err = h.store.GetPlayerById(id); err != nil {
		return NewApiError(err, "player not found", http.StatusNotFound)
	}

	if err := h.store.DeletePlayer(id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
