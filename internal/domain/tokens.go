package domain

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

// TokenClaims is a token claims.
type TokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"userId"`
}

// HashParams is a password hash params.
type HashParams struct {
	Memory      uint32
	Iterations  uint32
	SaltLength  uint32
	KeyLength   uint32
	Parallelism uint8
}
