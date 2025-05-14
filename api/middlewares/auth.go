package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/condemo/nes-cards-backend/api/handlers"
	"github.com/condemo/nes-cards-backend/api/utils"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(next http.Handler) http.HandlerFunc {
	return handlers.MakeHandler(func(w http.ResponseWriter, r *http.Request) error {
		rawToken := r.Header.Get("Authorization")
		splitT := strings.Split(rawToken, "Bearer ")
		if len(splitT) < 2 {
			return handlers.ApiError{
				Err:    errors.New("empty authorization or no bearer prefix"),
				Msg:    "empty authorization or bad format",
				Status: http.StatusUnauthorized,
			}
		}
		token := splitT[1]

		if token == "" {
			return handlers.ApiError{
				Err:    errors.New("empty token"),
				Msg:    "authorization header not found",
				Status: http.StatusUnauthorized,
			}
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return handlers.ApiError{
					Err:    err,
					Msg:    "access_token, is expired",
					Status: http.StatusGone,
				}
			}
			if errors.Is(err, jwt.ErrTokenMalformed) {
				return handlers.ApiError{
					Err:    err,
					Msg:    "invalid token format",
					Status: http.StatusUnauthorized,
				}
			}
			return err
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.Clone(ctx))

		return nil
	})
}
