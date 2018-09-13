package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func Basehandler() http.Handler {
	return http.HandlerFunc(BasehandlerFunc)
}

func BasehandlerFunc(w http.ResponseWriter, r *http.Request) {
	type timeval struct {
		TimeValue int64 `json:"value"`
	}


	// Add 10 seconds delay just for demo purposes
	// time.Sleep(10 * time.Millisecond)

	fmt.Println("...Before Basehandler")
	v := timeval{time.Now().UnixNano()}
	enc := json.NewEncoder(w)
	enc.Encode(v)
	fmt.Println("...After Basehandler")
}
