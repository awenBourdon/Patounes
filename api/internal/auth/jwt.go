package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

const TokenDuration = 7 * 24 * time.Hour

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("Error jwtSecret is undefined.")
	}

	jwtSecret = []byte(secret)
}

func GenerateToken(userID uuid.UUID, email string) (string, error) {
	claims := Claims{
		UserID: userID.String(),
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "patounes-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString((jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("méthode de signature inattendue")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token invalide")
	}

	return claims, nil
}

func GetUserIDFromToken(tokenString string) (string, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
