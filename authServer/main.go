package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var jwtSecret = []byte("your-secret-key")

var oauthConf = &oauth2.Config{
	ClientID:     "YOUR_GOOGLE_CLIENT_ID",
	ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
	RedirectURL:  "http://localhost:8080/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleLogin)
	r.HandleFunc("/callback", handleCallback)

	fmt.Println("Auth server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConf.AuthCodeURL("random-state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "OAuth exchange failed", http.StatusInternalServerError)
		return
	}

	client := oauthConf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var user struct {
		Email string `json:"email"`
	}
	json.NewDecoder(resp.Body).Decode(&user)

	// Create JWT
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := jwtToken.SignedString(jwtSecret)

	// Redirect to dashboard with token
	redirectURL := fmt.Sprintf("http://localhost:8081/dashboard?token=%s", tokenString)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
