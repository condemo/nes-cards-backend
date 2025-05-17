package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/condemo/nes-cards-backend/api/utils"
	"github.com/condemo/nes-cards-backend/store"
	"github.com/condemo/nes-cards-backend/types"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	store store.Store
}

func NewAuthHandler(s store.Store) *AuthHandler {
	return &AuthHandler{store: s}
}

func (h *AuthHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("POST /login", MakeHandler(h.login))
	r.HandleFunc("POST /signup", MakeHandler(h.signup))
	r.HandleFunc("POST /refresh", MakeHandler(h.refresh))
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) error {
	us := r.FormValue("username")
	pass := r.FormValue("password")

	user, err := h.store.GetUserByUsername(us)
	if err != nil {
		return ApiError{
			Err:    err,
			Status: http.StatusNotFound,
			Msg:    "user not found",
		}
	}

	if ok := utils.PassVerify(pass, user.Password); !ok {
		return ApiError{
			Err:    errors.New("invalid password"),
			Status: http.StatusUnauthorized,
			Msg:    "invalid password",
		}
	}

	token, err := utils.CreateJWT(user.ID)
	if err != nil {
		return err
	}

	refreshToken, err := utils.CreateRefreshJWT(user.ID)
	if err != nil {
		return err
	}

	SendJSON(w, http.StatusOK, map[string]string{
		"access_token":  token,
		"refresh_token": refreshToken,
		"token_type":    "bearer",
	})

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

func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) error {
	rawToken := r.Header.Get("Authorization")
	splitT := strings.Split(rawToken, "Bearer ")
	if len(splitT) < 2 {
		return ApiError{
			Err:    errors.New("empty authorization or bad format"),
			Msg:    "empty authorization or bad format",
			Status: http.StatusUnauthorized,
		}
	}
	token := splitT[1]

	if token == "" {
		return ApiError{
			Err:    errors.New("empty token"),
			Msg:    "Authorization header not found",
			Status: http.StatusUnauthorized,
		}
	}

	claims, err := utils.ValidateJWT(token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return ApiError{
				Err:    err,
				Msg:    "access_token is expired",
				Status: http.StatusGone,
			}
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return ApiError{
				Err:    err,
				Msg:    "invalid token format",
				Status: http.StatusUnauthorized,
			}
		}
	}

	newToken, err := utils.CreateJWT(claims.UserID)
	if err != nil {
		return err
	}

	SendJSON(w, http.StatusOK, map[string]string{
		"access_token": newToken,
		"token_type":   "bearer",
	})

	return nil
}
