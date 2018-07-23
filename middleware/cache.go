package middleware

import (
	"time"
	"sync"
	"net/http"
	"fmt"
)

type CacheItem struct {
	Key string
	Value   interface{}
	timeout int64
}

var cache_map map[string]CacheItem
var cache_channel = make(chan CacheItem )
var mutex = &sync.Mutex{}

func init() {
	cache_map = make(map[string]CacheItem)
	// start the cache gargage collector
	go cacheGc()
	go processCache()
}

func cacheGc() {
	for {
		// Run every n seconds
		time.Sleep(5)

		for k, v := range cache_map {
			if time.Now().UnixNano() > v.timeout {
				mutex.Lock()
				delete(cache_map, k)
				fmt.Println("deleting cache item " + k)
				mutex.Unlock()
			}
		}
	}
}

func processCache(){
	for {
		// Read the channel
		// Update the cache
		v := <-cache_channel

		// Update the cache, with timeout n seconds ahead
		v.timeout = time.Now().UnixNano() + 5000000000
		mutex.Lock()
		cache_map[v.Key] = v
		mutex.Unlock()
	}
}

func GetCacheItem(k string)( bool, CacheItem){
	if v, ok := cache_map[k]; ok {
		return true, v
	}

	return false, CacheItem{}
}

func CacheMiddleware() Adapter {
	return func(h http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("...Before CacheMiddleware")

			// Check if item in cache, if it is then return the cached item

			h.ServeHTTP(w, r)


			exists, _ := GetCacheItem(r.RequestURI); if !exists {
				cache_channel <- CacheItem{Key: r.RequestURI, Value: "{1}"}
			}

			fmt.Println("...After CacheMiddleware")
		}

		return http.HandlerFunc(fn)
	}
}
