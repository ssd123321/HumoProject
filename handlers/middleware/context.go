package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

func SetLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
		defer cancel()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := recover(); err != nil {
			log.Printf("Panic recovered: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		next.ServeHTTP(w, r)
	})
}
