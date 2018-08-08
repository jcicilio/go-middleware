package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type CacheItem struct {
	Key     string
	Value   []byte
	timeout int64
}

var cacheMap map[string]CacheItem
var cacheChannel = make(chan CacheItem)
var mutex = &sync.Mutex{}

const SecondsToCache = 10

// Responsible for initializing the cache map
// Responseible for starting the garbage collector (expired cache items)
// Responsible for starting the cache processor
func init() {
	cacheMap = make(map[string]CacheItem)
	// start the cache gargage collector
	go cacheGc()
	go processCache()
}

// Checks every second for cache entries that
// have been deleted and then expires them.
func cacheGc() {
	for {
		// Run GC every n seconds
		time.Sleep(1)

		for k, v := range cacheMap {
			if time.Now().UnixNano() > v.timeout {
				mutex.Lock()
				delete(cacheMap, k)
				fmt.Printf("...deleting cache key: %s\n", k)
				mutex.Unlock()
			}
		}
	}
}

// Watches the cache channel and updates the cache as
// new entries arrive at the channel
func processCache() {
	for {
		// Read the channel
		// Update the cache
		v := <-cacheChannel

		// Update the cache, with timeout n seconds ahead
		v.timeout = time.Now().UnixNano() + (SecondsToCache * 1000000000)
		mutex.Lock()
		fmt.Printf("...adding cache key: %s\n", v.Key)
		cacheMap[v.Key] = v
		mutex.Unlock()
	}
}

// Returns a cache item
func GetCacheItem(k string) (bool, CacheItem) {
	if v, ok := cacheMap[k]; ok {
		return true, v
	}

	return false, CacheItem{}
}

// A type used to implement the
// response writer interface in order to be able
// to capture bytes written to the client for
// caching.
type cacheRecorder struct {
	http.ResponseWriter
	r []byte
}

// response writer interface
func (rec *cacheRecorder) Write(b []byte) (int, error) {
	rec.r = b
	return rec.ResponseWriter.Write(b)
}

// response writer interface
func (rec *cacheRecorder) WriteHeader(code int) {
	rec.ResponseWriter.WriteHeader(code)
}

// response writer interface
func (rec *cacheRecorder) Header() http.Header {
	return rec.ResponseWriter.Header()
}

// returns an item from cache if found, otherwise
// captures the response stream an caches it
func CacheMiddleware() Adapter {
	return func(h http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("...Before CacheMiddleware")

			// Check if item is in cache, if it is, then stream to output and return
			exists, v := GetCacheItem(r.RequestURI)
			if exists {
				w.Write(v.Value)
				return
			}

			c := cacheRecorder{w, nil}

			h.ServeHTTP(&c, r)

			exists, _ = GetCacheItem(r.RequestURI)
			if !exists {
				cacheChannel <- CacheItem{Key: r.RequestURI, Value: c.r[:]}
			}

			fmt.Println("...After CacheMiddleware")
		}

		return http.HandlerFunc(fn)
	}
}
