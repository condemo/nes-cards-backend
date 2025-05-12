package handlers

import (
	"net/http"

	"github.com/condemo/nes-cards-backend/store"
)

type AuthHandler struct {
	store store.Store
}

func NewAuthHandler(s store.Store) *AuthHandler {
	return &AuthHandler{store: s}
}

func (h *AuthHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("POST /login", makeHandler(h.login))
	r.HandleFunc("POST /signup", makeHandler(h.signup))
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) error {
	// TODO:
	return nil
}

func (h *AuthHandler) signup(w http.ResponseWriter, r *http.Request) error {
	// TODO:
	return nil
}
