package main

import (
	"net/http"
	"math/rand"
	"time"
	"fmt"
	"encoding/json"
)

type Adapter func(http.Handler) http.Handler

type timeval struct {
	TimeValue int64 `json:"value"`
}

func Compose(h http.Handler, adapters ...Adapter) http.Handler {
	for i:=len(adapters)-1; i>=0; i--{
		h = adapters[i](h)
	}

	return h
}



func Basehandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("base 2 start")
		v := timeval{time.Now().UnixNano()}
		enc := json.NewEncoder(w)
		enc.Encode(v)
		fmt.Println("base 2 end")
	}

	return http.HandlerFunc(fn)
}

func HeadersMiddleware() Adapter {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("header 3 start")
			s1 := rand.NewSource(time.Now().UnixNano())

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Correlation-Id", fmt.Sprintf("%d", rand.New(s1).Int()))
			w.Header().Add("Access-Control-Allow-Origin", "*")
			h.ServeHTTP(w, r)
			fmt.Println("header 3 end")
		}

		return http.HandlerFunc(fn)
	}
}

func Logger() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("logger start")
			start := time.Now().UnixNano()
			h.ServeHTTP(w, r)

			//time.Sleep(2 * time.Second)
			end := time.Now().UnixNano()
			fmt.Printf("start: %d  end: %d duration:%d\n", start, end, end-start)
			fmt.Println("logger end")
		})
	}
}



func main() {
	composed := Compose(Basehandler(), Logger(), HeadersMiddleware())
	http.Handle("/",composed)

	http.ListenAndServe(":8080", nil)
}
