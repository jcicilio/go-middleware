package middleware

import (
	"time"
	"sync"
	"net/http"
	"fmt"
)

type CacheItem struct {
	Key string
	Value   []byte
	timeout int64
}

var cache_map map[string]CacheItem
var cache_channel = make(chan CacheItem )
var mutex = &sync.Mutex{}
const SECONDS_TO_CACHE = 5

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
				fmt.Printf("...deleting cache key: %s\n",k)
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
		v.timeout = time.Now().UnixNano() + (SECONDS_TO_CACHE * 100000000)
		mutex.Lock()
		fmt.Printf("...adding cache key: %s\n", v.Key)
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

type cacheRecorder struct {
	http.ResponseWriter
	r []byte
}

func (rec *cacheRecorder) Write(b []byte)(int, error) {
	rec.r = b
	return rec.ResponseWriter.Write(b)
}


func (rec *cacheRecorder) WriteHeader(code int) {
	rec.ResponseWriter.WriteHeader(code)
}


func (rec *cacheRecorder) Header() http.Header{
	return rec.ResponseWriter.Header()
}

func CacheMiddleware() Adapter {
	return func(h http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("...Before CacheMiddleware")

			// Check if item is in cache, if it is, then stream to output and return
			exists, v := GetCacheItem(r.RequestURI); if exists {
				//enc := json.NewEncoder(w)
				//enc.Encode(v.Value)
				w.Write(v.Value)
				return
			}

			c := cacheRecorder{w, nil}

			h.ServeHTTP(&c, r)


			exists, _ = GetCacheItem(r.RequestURI); if !exists {
				//stringToCache := string(c.r[:])
				cache_channel <- CacheItem{Key: r.RequestURI, Value: c.r[:]}
			}
			
			fmt.Println("...After CacheMiddleware")
		}

		return http.HandlerFunc(fn)
	}
}
