package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := r.Header.Get("token")
		if cookie == "" {
			http.Error(w, "unauthorized: no token cookie", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrInvalidKeyType
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
