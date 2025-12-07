package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Ambil secret dari environment, fallback ke default
func jwtSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "MY_SECRET_KEY" // fallback
	}
	return []byte(secret)
}

// GenerateJWT -> membuat token baru
func GenerateJWT(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // token berlaku 24 jam
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey())
}

// ParseToken -> memvalidasi token dan mengembalikan claims
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey(), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
