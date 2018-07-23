package main

import (
	"go-middleware/handlers"
	"go-middleware/middleware"
	"net/http"
)

func main() {
	withHeaders := middleware.Compose(handlers.Basehandler(), middleware.HeadersMiddleware())
	withHeadersAndLogging := middleware.Compose(handlers.Basehandler(), middleware.HeadersMiddleware(), middleware.LoggingMiddleware())
	withCaching:= middleware.Compose(handlers.Basehandler(), middleware.CacheMiddleware())

	http.Handle("/", handlers.Basehandler())
	http.Handle("/withHeaders", withHeaders)
	http.Handle("/withHeadersAndLogging", withHeadersAndLogging)
	http.Handle("/withCaching", withCaching)

	http.ListenAndServe(":8080", nil)
}
