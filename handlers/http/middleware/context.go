package middleware

import (
	"Tasks/utils"
	"context"
	"fmt"
	"log"
	"net/http"
)

/*
	func SetLimit(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
			defer cancel()
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
*/
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := recover(); err != nil {
			log.Printf("Panic recovered: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		next.ServeHTTP(w, r)
	})
}
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessJWT := r.Header.Get("accessJWT")
		if accessJWT == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(utils.ConvertToJson("unauthorized", "invalid accessJWT"))
			return
		}
		refreshJWT := r.Header.Get("refreshJWT")
		_, err := utils.ValidateRefreshJWT(refreshJWT)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(utils.ConvertToJson("unauthorized", "invalid refreshJWT"))
			return
		}
		id, err := utils.ValidateAccessJWT(accessJWT)
		if err != nil {
			if err.Error() != "token expired" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(utils.ConvertToJson("unsuccessfully", "invalid accessJWT, get again"))
				return
			} else if err.Error() == "token expired" {
				_, err = utils.ValidateRefreshJWT(refreshJWT)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(utils.ConvertToJson("unauthorized", "login again"))
					return
				}
				access, err := utils.GetAccessFromRefresh(refreshJWT)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write(utils.ConvertToJson("unauthorized", "login again"))
					return
				}
				w.Header().Set("accessJWT", access)
				w.Header().Set("refreshJWT", refreshJWT)
				return
			}
		}
		fmt.Println(123)
		r = r.WithContext(context.WithValue(r.Context(), "id", int(id)))
		next.ServeHTTP(w, r)
	})
}
