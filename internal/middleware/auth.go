package middleware

import (
	"context"
	"net/http"

	"github.com/Sarvesh-10/ReadEazeBackend/utility"
	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

func JWTMiddleWare(next http.Handler) http.Handler {
	logger := utility.NewLogger()
	logger.Info("GOT REAUEST")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			logger.Error("token is missing")
			http.Error(w, "Unauthorized - Token Missing", http.StatusUnauthorized)
			return
		}
		logger.Info("Priinting token", cookie.Value)
		token, err := utility.ValidateToken(cookie.Value)

		if err != nil {
			http.Error(w, "Unauthorized Invalid Token", http.StatusUnauthorized)
			return
		}
		cliams, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Unauthorized - Invalid Token Claims", http.StatusUnauthorized)
			return
		}
		userId, ok := cliams["user_id"].(float64)
		if !ok {
			http.Error(w, "Unauthorized - Invalid UserID", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, int(userId))

		// Do stuff here
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
