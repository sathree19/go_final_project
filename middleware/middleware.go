package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//pass := os.Getenv("TODO_PASSWORD")

		// if len(pass) > 0 {
		// 	var jwt string // JWT-токен из куки
		// 	//получаем куку
		// 	cookie, err := r.Cookie("token")
		// 	if err == nil {
		// 		jwt = cookie.Value
		// 	}
		// 	fmt.Println(jwt)

		// 	var valid bool

		// 	secret1 := []byte(pass)

		// 	var passW struct {
		// 		Password string `json:"password"`
		// 	}

		// 	err = json.NewDecoder(r.Body).Decode(&passW)

		// 	if err != nil {
		// 		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, "Ошибка"), http.StatusBadRequest)
		// 		return
		// 	}

		// 	secret2 := []byte(passW.Password)

		// 	// создаём jwt и указываем алгоритм хеширования
		// 	jwtToken := j.New(j.SigningMethodHS256)

		// 	// получаем подписанный токен
		// 	signedToken1, err := jwtToken.SignedString(secret1)
		// 	if err != nil {
		// 		fmt.Printf("failed to sign jwt: %s\n", err)
		// 	}
		// 	signedToken2, err := jwtToken.SignedString(secret2)
		// 	if err != nil {
		// 		fmt.Printf("failed to sign jwt: %s\n", err)
		// 	}

		// 	// jwtToken1, err := j.Parse(signedToken1, func(t *j.Token) (interface{}, error) {
		// 	// 	return secret1, nil
		// 	// })
		// 	// if err != nil {
		// 	// 	fmt.Printf("Failed to parse token2: %s\n", err)
		// 	// 	return
		// 	// }

		// 	jwtToken2, err := j.Parse(jwt, func(t *j.Token) (interface{}, error) {
		// 		return secret1, nil
		// 	})
		// 	if err != nil {
		// 		fmt.Printf("Failed to parse token2: %s\n", err)
		// 		return
		// 	}
		// 	valid = signedToken1 == signedToken2 && jwtToken2.Valid

		// 	if !valid {
		// 		// возвращаем ошибку авторизации 401
		// 		http.Error(w, "Authentification required", http.StatusUnauthorized)
		// 		return
		// 	}
		// }
		//next(w, r)
	})
}

func Auth2(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
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
			fmt.Println("1", signedToken)
			fmt.Println("2", j)

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
