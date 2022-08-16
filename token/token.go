package token

import (
	"fmt"
	"time"

	"github.com/varopxndx/chat/model"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

const defaultExpirationTime = time.Hour

// TokenService has the secret key
type TokenService struct {
	hmacSecret string
}

// Claims of the JWT
type Claims struct {
	UserName string `json:"user-name"`
	jwt.StandardClaims
}

// New creates a token service
func New(secret string) *TokenService {
	return &TokenService{hmacSecret: secret}
}

// GenerateToken creates a jwt token signed with the secret key
func (t TokenService) GenerateToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user-name": user.UserName,
		"exp":       time.Now().Add(defaultExpirationTime).Unix(),
	})
	return token.SignedString([]byte(t.hmacSecret))
}

// ValidateToken returns error if the jwt is invalid
func (t TokenService) ValidateToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Invalid signing method: %s", token.Header["alg"]))
		}
		return []byte(t.hmacSecret), nil
	})
	if _, ok := token.Claims.(*Claims); ok && token.Valid {
		return nil
	}
	return err
}
