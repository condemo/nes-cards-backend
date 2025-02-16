package middlewares

import (
	"fmt"
	"net/http"
)

type Middleware func(next http.Handler) http.HandlerFunc

func MiddlewareStack(m ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next.ServeHTTP
	}
}

func Recover(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// TODO: Implementar para usar con vue
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func AddCors(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// FIX: Actualmente hay una brecha de seguridad al permitir a todo el mundo
		// acceder a esta ruta
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
