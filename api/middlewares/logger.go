package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

type wrapperResponse struct {
	http.ResponseWriter
	status int
}

func (wr *wrapperResponse) WriteHeader(status int) {
	wr.ResponseWriter.WriteHeader(status)
	wr.status = status
}

var (
	fail    = "\033[38;5;160m"
	info    = "\033[38;5;39m"
	success = "\033[38;5;42m"
)

func Logger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wr := &wrapperResponse{w, http.StatusOK}
		next.ServeHTTP(wr, r)

		var c string
		if wr.status >= 400 {
			c = fail
		} else if wr.status >= 300 {
			c = info
		} else {
			c = success
		}

		fmt.Printf("%s[%s] %s [%s] %d - %s\033[0m\n", c, start.Format(time.DateTime), r.URL.Path, r.Method, wr.status, time.Since(start))
	}
}
