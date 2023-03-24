package token

import "github.com/golang-jwt/jwt"

// Tokens is a signin and refresh token response.
type Tokens struct {
	// Token Type
	TokenType string `json:"tokenType"`
	// Access Token
	AccessToken string `json:"accessToken"`
	// Refresh Token
	RefreshToken string `json:"refreshToken"`
	// Refresh Token Hash
	RefreshTokenHash string `json:"refreshTokenHash,omitempty"`
	// Access Token Expires datetime
	AccessTokenExpires int64 `json:"accessTokenExpires"`
}

// RefreshToken is a refresh token input.
type RefreshToken struct {
	// Refresh Token
	Authorization string `json:"token"`
}

// Claims is a token claims.
type Claims struct {
	jwt.StandardClaims
	UserID string `json:"userId"`
	Hash   string `json:"hash"`
}
