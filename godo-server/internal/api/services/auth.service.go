package services

import (
	"errors"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthService interface {
	GenerateJWT(email string, username string) (string, error)
	ValidateToken(signedToken string) (err error)
}

type authService struct {
	log   ilog.StdLogger
	query repository.ApiUserQuery
}

func NewAuthService(apiUserQuery repository.ApiUserQuery, logger ilog.StdLogger) AuthService {
	return &authService{
		log:   logger,
		query: apiUserQuery,
	}
}

// TODO: Place key in config
var jwtKey = []byte("remove this key")

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func (s *authService) GenerateJWT(email string, username string) (token string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := JWTClaim{
		Email:    email,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err = t.SignedString(jwtKey)
	return
}

func (s *authService) ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("Could not parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("Token has expired")
	}

	return
}
