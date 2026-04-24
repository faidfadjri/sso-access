package jwt

import (
	"akastra-access/internal/app/config"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID   uint64  `json:"user_id"`
	Email    string  `json:"email"`
	Username string  `json:"username"`
	FullName string  `json:"fullname"`
	Phone    *string `json:"phone"`
	Photo    *string `json:"photo"`
	ServiceName *string    `json:"service_name"`
	RoleName *string    `json:"role_name"`
	IDPRole  *string `json:"idp_role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// TokenUser is satisfied by both *domain.Users and *domain.UserWithService.
type TokenUser interface {
	GetUserID() uint64
	GetEmail() string
	GetUsername() string
	GetFullName() string
	GetPhone() *string
	GetPhoto() *string
	GetServiceName() *string
	GetRoleName() *string
	GetIDPRole() *string
}

func GenerateAccessToken(user TokenUser) (string, error) {
	conf := config.Load()
	claims := &AccessClaims{
		UserID:   user.GetUserID(),
		Email:    user.GetEmail(),
		Username: user.GetUsername(),
		FullName: user.GetFullName(),
		Phone:    user.GetPhone(),
		Photo:    user.GetPhoto(),
		ServiceName: user.GetServiceName(),
		RoleName: user.GetRoleName(),
		IDPRole:  user.GetIDPRole(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.JWTSecret))
}

func ValidateAccessToken(tokenString string) (*AccessClaims, error) {
	conf := config.Load()
	tokenString = ExtractBearer(tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(conf.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GenerateRefreshToken(user TokenUser) (string, error) {
	conf := config.Load()
	claims := &RefreshClaims{
		UserID: user.GetUserID(),
		Email:  user.GetEmail(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.RefreshTokenSecret))
}

func ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	conf := config.Load()
	tokenString = ExtractBearer(tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(conf.RefreshTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// helper untuk strip "Bearer "
func ExtractBearer(token string) string {
	parts := strings.SplitN(token, " ", 2)
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return token
}
