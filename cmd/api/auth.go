package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/krutip7/chat-app-server/internals/models"
)

type Auth struct {
	Issuer               string
	Audience             string
	AuthTokenValidity    time.Duration
	RefreshTokenValidity time.Duration
	SigningKey           []byte
}

type TokenPair struct {
	AuthToken    string
	RefreshToken string
}

type JWTClaims struct {
	name     string
	username string
	jwt.RegisteredClaims
}

func (auth *Auth) GenerateJWTToken(user *models.User) (tokenPair TokenPair, err error) {

	authTokenClaims := JWTClaims{
		name:     fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "",
			Issuer:    auth.Issuer,
			Subject:   user.Id,
			Audience:  jwt.ClaimStrings{auth.Audience},
			IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().UTC().Add(auth.AuthTokenValidity)},
		},
	}

	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authTokenClaims)
	signedAuthToken, err := authToken.SignedString(auth.SigningKey)
	if err != nil {
		return
	}

	refreshTokenClaims := jwt.RegisteredClaims{
		ID:        "",
		Issuer:    auth.Issuer,
		Subject:   user.Id,
		Audience:  jwt.ClaimStrings{auth.Audience},
		IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
		ExpiresAt: &jwt.NumericDate{Time: time.Now().UTC().Add(auth.RefreshTokenValidity)},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedRefreshToken, err := refreshToken.SignedString(auth.SigningKey)

	tokenPair = TokenPair{
		AuthToken:    signedAuthToken,
		RefreshToken: signedRefreshToken,
	}

	return
}
