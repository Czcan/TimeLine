package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/Czcan/TimeLine/utils/jsonwt"
)

func JwtAuthentication(jwtClient jsonwt.JWTValidate) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth != "" {
				claim, err := jwtClient.ParseToken(auth)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				if time.Now().Unix() > claim.ExpiresAt {
					next.ServeHTTP(w, r)
					return
				}
				ctx := context.WithValue(r.Context(), "token", claim)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}
}
