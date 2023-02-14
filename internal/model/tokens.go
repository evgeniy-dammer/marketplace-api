package model

import "github.com/golang-jwt/jwt"

// Tokens is a signin and refresh token response.
type Tokens struct {
	TokenType          string `json:"tokenType"`
	AccessToken        string `json:"accessToken"`
	AccessTokenExpires int64  `json:"accessTokenExpires"`
	RefreshToken       string `json:"refreshToken"`
}

// RefreshToken is a refresh token input.
type RefreshToken struct {
	Authorization string `json:"token"`
}

// TokenClaims is a token claims.
type TokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"userId"`
}

// HashParams is a password hash params.
type HashParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}
