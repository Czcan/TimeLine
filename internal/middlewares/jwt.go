package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Czcan/TimeLine/pkg/jwt"
)

func JwtAuthentication(jwtClient jwt.JWTValidate) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if strings.TrimSpace(auth) != "" {
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
