package handlers

import (
	"net/http"

	"github.com/condemo/nes-cards-backend/api/utils"
	"github.com/condemo/nes-cards-backend/store"
	"github.com/condemo/nes-cards-backend/types"
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
	us := r.FormValue("username")
	pass := r.FormValue("password")

	user, err := h.store.GetUserByUsername(us)
	if err != nil {
		return ApiError{
			Err:    err,
			Status: http.StatusNotFound,
			Msg:    "user not found"}
	}

	if ok := utils.PassVerify(pass, user.Password); !ok {
		return ApiError{
			Err:    err,
			Status: http.StatusUnauthorized,
			Msg:    "invalid password",
		}
	}

	// TODO: JWT stuff

	SendJSON(w, http.StatusOK, user)
	return nil
}

func (h *AuthHandler) signup(w http.ResponseWriter, r *http.Request) error {
	us := r.FormValue("username")
	pass, err := utils.PassEncrypt(r.FormValue("password"))
	if err != nil {
		return ApiError{
			Err:    err,
			Status: http.StatusBadRequest,
			Msg:    "invalid password format",
		}
	}

	user := &types.User{
		Username: us,
		Password: pass,
	}

	if err := user.Validate(); err != nil {
		return ApiError{
			Err:    err,
			Status: http.StatusBadRequest,
			Msg:    "credentials validation error",
		}
	}

	if err := h.store.CreateUser(user); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}
