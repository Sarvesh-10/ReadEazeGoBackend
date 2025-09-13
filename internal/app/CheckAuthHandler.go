package app

import (
	"encoding/json"
	"net/http"

	"github.com/Sarvesh-10/ReadEazeBackend/utility"
	"github.com/golang-jwt/jwt/v5"
)

type AuthResponse struct {
	IsLoggedIn bool `json:"isLoggedIn"`
}

//I have to add this to the routers how to do that

func CheckAuthHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("token")
	if err != nil {

		http.Error(w, "Unauthorized - Token Missing", http.StatusUnauthorized)
		return
	}

	token, err := utility.ValidateToken(cookie.Value)

	if err != nil {
		http.Error(w, "Unauthorized Invalid Token", http.StatusUnauthorized)
		return
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AuthResponse{IsLoggedIn: false})
		return
	}

	// Success
	json.NewEncoder(w).Encode(AuthResponse{IsLoggedIn: true})
}
