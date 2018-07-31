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

func init() {
	cacheMap = make(map[string]CacheItem)
	// start the cache gargage collector
	go cacheGc()
	go processCache()
}

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

func GetCacheItem(k string) (bool, CacheItem) {
	if v, ok := cacheMap[k]; ok {
		return true, v
	}

	return false, CacheItem{}
}

type cacheRecorder struct {
	http.ResponseWriter
	r []byte
}

func (rec *cacheRecorder) Write(b []byte) (int, error) {
	rec.r = b
	return rec.ResponseWriter.Write(b)
}

func (rec *cacheRecorder) WriteHeader(code int) {
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *cacheRecorder) Header() http.Header {
	return rec.ResponseWriter.Header()
}

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
