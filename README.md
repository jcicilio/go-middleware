# go middleware

The functionality of the API is to return a random number in json format.
There are two routes to the API, one is cached the other is not cached.
Caching is done in memory for simplicity.

## the structure of middleware

It starts with func(http.Handler) http.Handler
and then function composition

type Adapter func(http.Handler) http.Handler

given a http.Handler h

then wrapping of the http.Handler
h  <== f1(f2(f3(f4(h))))

```
func SomeHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// perform some actions here

		// serve http
		h.ServeHTTP(w, r)
	}

    // pass back handler
	return http.HandlerFunc(fn)
}
```

## middleware examples

### console logging

responsibility:  log to console using [Apache Common Log Format `CLF`](http://httpd.apache.org/docs/2.2/logs.html#common)
responsibility:  log start time
responsibility:  log end time and duration

### caching

responsibility: cache by query URI

### cors, and common headers

responsibility: add header correlation id
responsibility: add header cors
responsibility: add header Content-Type


## references

[Writing middleware in #golang and how Go makes it so much fun.](https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81)

[Simple HTTP middleware with Go](https://hackernoon.com/simple-http-middleware-with-go-79a4ad62889b)

[Gorilla Handlers](http://www.gorillatoolkit.org/pkg/handlers#LoggingHandler)


https://hackernoon.com/simple-http-middleware-with-go-79a4ad62889b
https://www.nicolasmerouze.com/middlewares-golang-best-practices-examples/

best so far
https://gowebexamples.com/advanced-middleware/