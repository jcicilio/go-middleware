package main

import (
	"net/http"
	"time"
	"fmt"
	"go-middleware/handlers"
	"go-middleware/middleware"
)


// LoggingMiddleware provides an example of before and after logging in an
// API call stack, in this case providing duration of the call
func LoggingMiddleware() middleware.Adapter {
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
	withHeaders := middleware.Compose(handlers.Basehandler(),middleware.HeadersMiddleware())
	withHeadersAndLogging := middleware.Compose(handlers.Basehandler(),middleware.HeadersMiddleware(), LoggingMiddleware())

	http.Handle("/",handlers.Basehandler())
	http.Handle("/withHeaders",withHeaders)
	http.Handle("/withHeadersAndLogging",withHeadersAndLogging)

	http.ListenAndServe(":8080", nil)
}
