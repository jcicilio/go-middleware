package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// LoggingMiddleware provides an example of before and after logging in an
// API call stack, in this case providing duration of the call
func LoggingMiddleware() Adapter {
	return func(h http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("...Before LoggingMiddleware")
			start := time.Now().UnixNano()
			h.ServeHTTP(w, r)
			end := time.Now().UnixNano()
			fmt.Printf("start-time(ns): %d  end-time(ns): %d duration(ns):%d  \"%s %s\"\n", start, end,(end-start),r.Method, r.RequestURI)
			fmt.Println("...After LoggingMiddleware")
		}

		return http.HandlerFunc(fn)
	}
}



