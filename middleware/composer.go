package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Compose(h http.Handler, adapters ...Middleware) http.Handler {
	for i := len(adapters) - 1; i >= 0; i-- {
		h = adapters[i](h)
	}

	return h
}
