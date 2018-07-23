package cache

import "time"

type cacheItem struct{
	value interface{}
	timeout int64

}

// The cache
var cache_map map[string]cacheItem

func init(){
	// start the cache gargage collector
	go cacheGc()
}

func cacheGc(){
	// Run every n seconds
	time.Sleep(5)
}

