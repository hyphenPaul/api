package main

import (
	"context"
	"net/http"
	"time"
)

func jsonMw(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		r = r.WithContext(ctx)
		w.Header().Add("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func logMw(actx AppContext, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(t time.Time) {
			l := loggerPayload{
				Duration: time.Since(t).String(),
				URL:      r.URL.String(),
				Method:   r.Method,
			}

			actx.logger.info(l)
		}(time.Now())

		h.ServeHTTP(w, r)
	})
}
