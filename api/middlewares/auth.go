package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"
)

func RequireAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawToken := r.Header.Get("Authorization")
		splitT := strings.Split(rawToken, "Bearer ")
		token := splitT[1]

		// TODO: JWT validation, etc
		if token == "" {
			json.NewEncoder(w).Encode(map[string]any{
				"msg":    "invalid authorization",
				"status": http.StatusBadRequest,
			})
			return
		}
	}
}
