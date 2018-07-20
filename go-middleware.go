package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Borrowed extensively from
// https://gowebexamples.com/advanced-middleware/

type Middleware func(http.HandlerFunc) http.HandlerFunc

func HeadersMiddleware() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			s1 := rand.NewSource(time.Now().UnixNano())

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Correlation-Id", fmt.Sprintf("%d", rand.New(s1).Int()))
			w.Header().Add("Access-Control-Allow-Origin", "*")

			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Compose(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

type timeval struct {
	TimeValue int64 `json:"value"`
}

func Basehandler(w http.ResponseWriter, r *http.Request) {
	v := timeval{time.Now().UnixNano()}
	enc := json.NewEncoder(w)
	enc.Encode(v)
}

func main() {
	withHeaders := Compose(Basehandler, HeadersMiddleware())

	http.HandleFunc("/", Basehandler)
	http.HandleFunc("/withHeaders", withHeaders)
	http.ListenAndServe(":8080", nil)
}
