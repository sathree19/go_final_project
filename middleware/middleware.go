package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.HandlerFunc, pass string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля

		if len(pass) > 0 {
			var j string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				j = cookie.Value
			}
			var valid bool

			secret := []byte("12345")

			jwtToken := jwt.New(jwt.SigningMethodHS256)

			signedToken, err := jwtToken.SignedString(secret)
			if err != nil {
				fmt.Printf("failed to sign jwt: %s\n", err)
			}

			valid = j == signedToken

			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
