package middleware

import "net/http"

type Adapter func(http.Handler) http.Handler

func Compose(h http.Handler, adapters ...Adapter) http.Handler {
	for i:=len(adapters)-1; i>=0; i--{
		h = adapters[i](h)
	}

	return h
}