package main

import (
	"net/http"
	"log"
	"encoding/json"
	"time"
)

type randval struct{
	V int64 `json:"value""`
}

func basehandler(w http.ResponseWriter, r *http.Request) {
	v := randval{time.Now().UnixNano()}

	enc := json.NewEncoder(w)
	enc.Encode(v)
}


func main() {
	http.HandleFunc("/", basehandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
