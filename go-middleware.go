package main

import (
	"fmt"
	"go-middleware/handlers"
	"go-middleware/middleware"
	"net/http"
)

func main() {
	// build out different handlers
	withHeaders := middleware.Compose(handlers.Basehandler(), middleware.HeadersMiddleware())
	withHeadersAndLogging := middleware.Compose(handlers.Basehandler(), middleware.HeadersMiddleware(), middleware.LoggingMiddleware())
	withCaching := middleware.Compose(handlers.Basehandler(), middleware.CacheMiddleware())
	withHeadersAndCaching := middleware.Compose(handlers.Basehandler(), middleware.HeadersMiddleware(), middleware.CacheMiddleware())

	http.Handle("/", handlers.Basehandler())
	http.Handle("/withHeaders", withHeaders)
	http.Handle("/withHeadersAndLogging", withHeadersAndLogging)
	http.Handle("/withCaching", withCaching)
	http.Handle("/withHeadersAndCaching", withHeadersAndCaching)

	fmt.Println("go-middleware started...")
	http.ListenAndServe(":8080", nil)
}
