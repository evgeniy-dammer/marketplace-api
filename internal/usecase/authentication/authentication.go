package authentication

import (
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

// AuthenticationGenerateToken generates authorization token.
func (s *UseCase) AuthenticationGenerateToken(userID string, username string, password string) (user.User, domain.Tokens, error) {
	var usr user.User

	var tokens domain.Tokens

	var err error

	usr, err = s.adapterStorage.AuthenticationGetUser(userID, username)

	if err != nil {
		return usr, tokens, errors.Wrap(err, "can not get user")
	}

	if username != "" {
		match, err := usecase.ComparePasswordAndHash(password, usr.Password)
		if err != nil {
			return usr, tokens, err
		}

		if !match {
			return usr, tokens, usecase.ErrInvalidPassword
		}

		userID = usr.ID
	}

	usr.Password = ""
	usr.RoleID = 0

	issuedAt := time.Now().Unix()
	expiresAt := time.Now().Add(usecase.TokenTTL).Unix()
	refreshExpiresAt := time.Now().Add(usecase.RefreshTokenTTL).Unix()

	token := usecase.CreateNewToken(userID, expiresAt, issuedAt)
	tokens.AccessToken, err = token.SignedString([]byte(usecase.SigningKey))

	if err != nil {
		return usr, tokens, errors.Wrap(err, "can not get access token")
	}

	tokens.AccessTokenExpires = expiresAt
	tokens.TokenType = "Bearer"

	// create refresh token
	refreshToken := usecase.CreateNewToken(userID, refreshExpiresAt, issuedAt)
	tokens.RefreshToken, err = refreshToken.SignedString([]byte(usecase.SigningKey))

	if err != nil {
		return usr, tokens, errors.Wrap(err, "can not get refresh token")
	}

	return usr, tokens, nil
}

// AuthenticationParseToken checks access token and returns user id.
func (s *UseCase) AuthenticationParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &domain.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, usecase.ErrInvalidSigningMethod
		}

		return []byte(usecase.SigningKey), nil
	})
	if err != nil {
		return "", errors.Wrap(err, "can not parse token")
	}

	claims, ok := token.Claims.(*domain.TokenClaims)

	if !ok {
		return "", usecase.ErrInvalidTokenClaims
	}

	return claims.UserID, nil
}

// AuthenticationCreateUser hashes the password and insert User into system.
func (s *UseCase) AuthenticationCreateUser(user user.User) (string, error) {
	pass, err := usecase.GeneratePasswordHash(user.Password, usecase.Params)
	if err != nil {
		return "", err
	}

	user.Password = pass

	userID, err := s.adapterStorage.AuthenticationCreateUser(user)

	return userID, errors.Wrap(err, "can not create user")
}

// AuthenticationGetUserRole returns users role name.
func (s *UseCase) AuthenticationGetUserRole(id string) (string, error) {
	role, err := s.adapterStorage.AuthenticationGetUserRole(id)

	return role, errors.Wrap(err, "can not get role")
}
