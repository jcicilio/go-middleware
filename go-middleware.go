package main

import (
	"net/http"
	"math/rand"
	"time"
	"fmt"
	"go-middleware/handlers"
)

type Adapter func(http.Handler) http.Handler

func Compose(h http.Handler, adapters ...Adapter) http.Handler {
	for i:=len(adapters)-1; i>=0; i--{
		h = adapters[i](h)
	}

	return h
}



// The HeadersMiddleware provides an example of adding headers during
// the API run, including a correlation id, content-type and CORS headers.
func HeadersMiddleware() Adapter {
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

// LoggingMiddleware provides an example of before and after logging in an
// API call stack, in this case providing duration of the call
func LoggingMiddleware() Adapter {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("...Before LoggingMiddleware")
			start := time.Now().UnixNano()
			h.ServeHTTP(w, r)
			end := time.Now().UnixNano()
			fmt.Printf("...LoggingMiddleware, start-time: %d  end-time: %d duration(ns):%d\n", start, end, end-start)
			fmt.Println("...After LoggingMiddleware")
		}
		return http.HandlerFunc(fn)
	}
}



func main() {
	withHeaders := Compose(handlers.Basehandler(),HeadersMiddleware())
	withHeadersAndLogging := Compose(handlers.Basehandler(),HeadersMiddleware(), LoggingMiddleware())

	http.Handle("/",handlers.Basehandler())
	http.Handle("/withHeaders",withHeaders)
	http.Handle("/withHeadersAndLogging",withHeadersAndLogging)

	http.ListenAndServe(":8080", nil)
}
