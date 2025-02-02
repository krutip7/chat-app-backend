package auth

import (
	"errors"
	"fmt"
	"net/http"
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
	CookieName           string
}

type TokenPair struct {
	AuthToken    string
	RefreshToken string
}

type JWTClaims struct {
	Name     string
	Username string
	jwt.RegisteredClaims
}

func (auth *Auth) GenerateJWTToken(user *models.User) (tokenPair TokenPair, err error) {

	authTokenClaims := JWTClaims{
		Name:     fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
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

func (auth *Auth) GetRefreshTokenCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     auth.CookieName,
		Value:    refreshToken,
		MaxAge:   int(auth.RefreshTokenValidity.Seconds()),
		Expires:  time.Now().UTC().Add(auth.RefreshTokenValidity),
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
	}
}

func (auth *Auth) ClearRefreshTokenCookie() *http.Cookie {
	return &http.Cookie{
		Name:     auth.CookieName,
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
	}
}

func (auth *Auth) VerifyJWT(tokenString string) (JWTClaims, error) {
	var claims JWTClaims

	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return auth.SigningKey, nil
	})

	if err != nil {
		return claims, errors.New("invalid auth token")
	}

	return claims, nil
}
