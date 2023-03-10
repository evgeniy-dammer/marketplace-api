package token

import "github.com/golang-jwt/jwt"

// Tokens is a signin and refresh token response.
type Tokens struct {
	TokenType          string `json:"tokenType"`
	AccessToken        string `json:"accessToken"`
	RefreshToken       string `json:"refreshToken"`
	AccessTokenExpires int64  `json:"accessTokenExpires"`
}

// RefreshToken is a refresh token input.
type RefreshToken struct {
	Authorization string `json:"token"`
}

// Claims is a token claims.
type Claims struct {
	jwt.StandardClaims
	UserID string `json:"userId"`
}
