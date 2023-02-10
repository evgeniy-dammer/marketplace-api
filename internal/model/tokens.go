package model

import "github.com/dgrijalva/jwt-go"

type Tokens struct {
	TokenType          string `json:"tokenType"`
	AccessToken        string `json:"accessToken"`
	AccessTokenExpires int64  `json:"accessTokenExpires"`
	RefreshToken       string `json:"refreshToken"`
}

type RefreshToken struct {
	Authorization string `json:"token"`
}

// TokenClaims
type TokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

// HashParams
type HashParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}
