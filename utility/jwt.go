package utility

import (
	"errors"
	"time"

	"github.com/Sarvesh-10/ReadEazeBackend/config"
	"github.com/golang-jwt/jwt/v5"
)

// Secret key for signing tokens
var jwtSecret = []byte(config.AppConfig.JWTSecret)

// GenerateToken creates a new JWT token
func GenerateToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken parses and validates a JWT token
func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
}
