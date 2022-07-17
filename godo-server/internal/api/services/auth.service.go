package services

import (
	"errors"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthService interface {
	GenerateJWT(email string, username string, accountId string) (string, error)
	ValidateTokenClaims(signedToken string) (err error)
	GetClaims(signedToken string) (*JWTClaim, error)
	BearerTokenToToken(token string) (string, error)
}

type authService struct {
	jwtSecret []byte
	log       ilog.StdLogger
	query     repository.ApiUserQuery
}

func NewAuthService(apiUserQuery repository.ApiUserQuery, jwtSecret []byte, logger ilog.StdLogger) AuthService {
	return &authService{
		query:     apiUserQuery,
		jwtSecret: jwtSecret,
		log:       logger,
	}
}

type JWTClaim struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	AccountId string `json:"account_id"`
	jwt.StandardClaims
}

func (s *authService) GenerateJWT(email string, username string, accountId string) (token string, err error) {
	expirationTime := time.Now().Add(6 * time.Hour)

	claims := JWTClaim{
		Email:     email,
		Username:  username,
		AccountId: accountId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err = t.SignedString(s.jwtSecret)
	return
}

func (s *authService) ValidateTokenClaims(signedToken string) (err error) {
	claims, err := s.GetClaims(signedToken)
	if err != nil {
		return err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		s.log.Info("The token has expired")
		err = errors.New("Token has expired")
	}

	return nil
}

func (s *authService) GetClaims(signedToken string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return s.jwtSecret, nil
		},
	)

	if err != nil {
		s.log.Info("There has been an issue processing the claims from the signed token")
		s.log.Error(err)
		return nil, errors.New("There has been an issue processing the claims from the signed token")
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		s.log.Info("Could not parse claims")
		return nil, errors.New("Could not parse claims")
	}

	return claims, nil
}

func (s *authService) BearerTokenToToken(token string) (string, error) {
	err := s.validateTokenString(token)
	if err != nil {
		return "", err
	}

	tokenParts := strings.SplitAfter(token, " ")
	if len(tokenParts) != 2 {
		s.log.Info("The bearer token does not follow the correct format")
		return "", errors.New("The bearer token does not follow the correct format")
	}

	return tokenParts[1], nil
}

func (s *authService) validateTokenString(token string) error {
	if token == "" {
		s.log.Info("Token is a blank string")
		return errors.New("The token is an empty string")
	}

	// Validate that the token string looks alright
	if !strings.HasPrefix(token, "Bearer") {
		s.log.Info("Bad token: the token does not start with `Bearer`")
		return errors.New("Bearer not included at the beginning of the token")
	}

	tokenParts := strings.SplitAfter(token, " ")
	if len(tokenParts) != 2 {
		s.log.Info("The bearer token does not follow the correct format")
		return errors.New("The bearer token does not follow the correct format")
	}

	return nil
}
