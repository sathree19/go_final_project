package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func PostSign(w http.ResponseWriter, r *http.Request) {

	var pass struct {
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&pass)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, "Ошибка"), http.StatusBadRequest)
		return
	}

	if pass.Password == os.Getenv("TODO_PASSWORD") {

		jwtToken := jwt.New(jwt.SigningMethodHS256)

		secret := []byte(pass.Password)

		signedToken, err := jwtToken.SignedString(secret)
		if err != nil {
			fmt.Printf("failed to sign jwt: %s\n", err)
		}

		fmt.Fprintf(w, `{"token": "%s"}`, signedToken)

	} else {

		fmt.Fprintf(w, `{"error": "%s"}`, "Ошибка")
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

}
