package middleware

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// The HeadersMiddleware provides an example of adding headers during
// the API run, including a correlation id, content-type and CORS headers.
func HeadersMiddleware() Middleware {
	return func(h http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("...Before HeadersMiddleware")
			s1 := rand.NewSource(time.Now().UnixNano())

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Correlation-Id", fmt.Sprintf("%d", rand.New(s1).Int()))
			w.Header().Add("Access-Control-Allow-Origin", "*")
			h.ServeHTTP(w, r)
			fmt.Println("...After HeadersMiddleware")
		}

		return http.HandlerFunc(fn)
	}
}
