package authentication

import (
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/token"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

// AuthenticationGenerateToken generates authorization token.
func (s *UseCase) AuthenticationGenerateToken(ctx context.Context, userID string, username string, password string) (user.User, token.Tokens, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.AuthenticationGenerateToken")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var usr user.User

	var tokens token.Tokens

	var err error

	usr, err = s.adapterStorage.AuthenticationGetUser(ctx, userID, username)

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

	tkn := usecase.CreateNewToken(userID, expiresAt, issuedAt)
	tokens.AccessToken, err = tkn.SignedString([]byte(usecase.SigningKey))

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
func (s *UseCase) AuthenticationParseToken(ctx context.Context, accessToken string) (string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.AuthenticationParseToken")
		defer span.End()

		_ = context.New(ctxt)
	}

	tkn, err := jwt.ParseWithClaims(accessToken, &token.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, usecase.ErrInvalidSigningMethod
		}

		return []byte(usecase.SigningKey), nil
	})
	if err != nil {
		return "", errors.Wrap(err, "can not parse token")
	}

	claims, ok := tkn.Claims.(*token.Claims)

	if !ok {
		return "", usecase.ErrInvalidTokenClaims
	}

	return claims.UserID, nil
}

// AuthenticationCreateUser hashes the password and insert User into system.
func (s *UseCase) AuthenticationCreateUser(ctx context.Context, input user.CreateUserInput) (string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.AuthenticationCreateUser")
		defer span.End()

		ctx = context.New(ctxt)
	}

	pass, err := usecase.GeneratePasswordHash(input.Password, usecase.Params)
	if err != nil {
		return "", err
	}

	input.Password = pass

	userID, err := s.adapterStorage.AuthenticationCreateUser(ctx, input)

	return userID, errors.Wrap(err, "can not create user")
}
