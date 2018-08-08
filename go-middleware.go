package main

import (
	"fmt"
	"go-middleware/handlers"
	"go-middleware/middleware"
	"net/http"
	"log"
)

func main() {
	// build out different handlers
	withHeaders := middleware.Compose(handlers.Basehandler(), middleware.HeadersMiddleware())
	withHeadersAndLogging := middleware.Compose(handlers.Basehandler(), middleware.HeadersMiddleware(), middleware.LoggingMiddleware())
	withCaching := middleware.Compose(handlers.Basehandler(), middleware.CacheMiddleware())
	withHeadersAndCaching := middleware.Compose(handlers.Basehandler(), middleware.HeadersMiddleware(), middleware.CacheMiddleware())
	withAllDemoMiddleware := middleware.Compose(handlers.Basehandler(), middleware.LoggingMiddleware(),middleware.HeadersMiddleware(), middleware.CacheMiddleware())

	http.Handle("/", handlers.Basehandler())
	http.Handle("/withHeaders", withHeaders)
	http.Handle("/withHeadersAndLogging", withHeadersAndLogging)
	http.Handle("/withCaching", withCaching)
	http.Handle("/withHeadersAndCaching", withHeadersAndCaching)
	http.Handle("/withAllDemoMiddleware", withAllDemoMiddleware)

	fmt.Println("go-middleware started...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
