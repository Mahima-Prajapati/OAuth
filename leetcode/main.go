package main

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var jwtSecret = []byte("your-secret-key")

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleHome)

	fmt.Println("LeetCode app running on http://localhost:8082")
	http.ListenAndServe(":8082", r)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		fmt.Println(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["email"].(string)
		fmt.Fprintf(w, "✅ Welcome to LeetCode Clone, %s!", email)
	} else {
		http.Error(w, "❌ Invalid or expired token", http.StatusUnauthorized)
	}
}
