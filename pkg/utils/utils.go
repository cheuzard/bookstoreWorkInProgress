package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func Chain(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func TimeoutMiddleware(h http.Handler) http.Handler {
	fmt.Printf("timeout middleware is being run \n")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		go func() {
			h.ServeHTTP(w, r)
			cancel()
		}()
		ctx.Done()
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			w.WriteHeader(http.StatusGatewayTimeout)
		}
		return
	})
}
