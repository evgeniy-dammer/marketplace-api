package authentication

import (
	"time"

	"github.com/evgeniy-dammer/marketplace-api/internal/domain/token"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/user"
	"github.com/evgeniy-dammer/marketplace-api/internal/usecase"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// AuthenticationGenerateToken generates authorization token.
func (s *UseCase) AuthenticationGenerateToken(ctx context.Context, userID string, username string, password string) (user.User, token.Tokens, error) { //nolint:lll
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
			return usr, tokens, errors.Wrap(err, "comparing")
		}

		if !match {
			return usr, tokens, usecase.ErrInvalidPassword
		}

		userID = usr.ID
	}

	usr.Password = ""
	usr.RoleID = 0

	// TODO: replace with time.Minute
	accessTokenTTL := time.Duration(viper.GetInt("authentication.access_token_ttl")) * time.Hour
	refreshTokenTTL := time.Duration(viper.GetInt("authentication.refresh_token_ttl")) * time.Hour

	issuedAt := time.Now().Unix()
	expiresAt := time.Now().Add(accessTokenTTL).Unix()
	refreshExpiresAt := time.Now().Add(refreshTokenTTL).Unix()

	tkn := usecase.CreateNewToken(userID, expiresAt, issuedAt, "")
	tokens.AccessToken, err = tkn.SignedString([]byte(viper.GetString("JWT_KEY")))

	if err != nil {
		return usr, tokens, errors.Wrap(err, "can not get access token")
	}

	tokens.AccessTokenExpires = expiresAt
	tokens.TokenType = "Bearer"

	hash, err := uuid.NewUUID()
	if err != nil {
		return usr, tokens, errors.Wrap(err, "can not create new hash")
	}

	refreshToken := usecase.CreateNewToken(userID, refreshExpiresAt, issuedAt, hash.String())
	tokens.RefreshTokenHash = hash.String()
	tokens.RefreshToken, err = refreshToken.SignedString([]byte(viper.GetString("JWT_KEY")))

	if err != nil {
		return usr, tokens, errors.Wrap(err, "can not get refresh token")
	}

	return usr, tokens, nil
}

// AuthenticationParseToken checks access token and returns user id.
func (s *UseCase) AuthenticationParseToken(ctx context.Context, accessToken string) (string, string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.AuthenticationParseToken")
		defer span.End()

		_ = context.New(ctxt)
	}

	tkn, err := jwt.ParseWithClaims(accessToken, &token.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, usecase.ErrInvalidSigningMethod
		}

		return []byte(viper.GetString("JWT_KEY")), nil
	})
	if err != nil {
		return "", "", errors.Wrap(err, "can not parse token")
	}

	claims, ok := tkn.Claims.(*token.Claims)

	if !ok {
		return "", "", usecase.ErrInvalidTokenClaims
	}

	return claims.UserID, claims.Hash, nil
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
		return "", errors.Wrap(err, "can not generate password hash")
	}

	input.Password = pass

	userID, err := s.adapterStorage.AuthenticationCreateUser(ctx, input)

	return userID, errors.Wrap(err, "can not create user")
}

// AuthenticationCreateTokenHash creates token hash in database.
func (s *UseCase) AuthenticationCreateTokenHash(ctx context.Context, userID string, hash string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.AuthenticationCreateTokenHash")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.AuthenticationCreateTokenHash(ctx, userID, hash)

	return errors.Wrap(err, "token create failed")
}

// AuthenticationGetTokenHash returns token hash id from database.
func (s *UseCase) AuthenticationGetTokenHash(ctx context.Context, userID string, hash string) (string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.AuthenticationGetTokenHash")
		defer span.End()

		ctx = context.New(ctxt)
	}

	tokenID, err := s.adapterStorage.AuthenticationGetTokenHash(ctx, userID, hash)

	return tokenID, errors.Wrap(err, "token id select error")
}

// AuthenticationUpdateTokenHash updates token hash in database.
func (s *UseCase) AuthenticationUpdateTokenHash(ctx context.Context, tokenID string, hash string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.AuthenticationUpdateTokenHash")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.AuthenticationUpdateTokenHash(ctx, tokenID, hash)

	return errors.Wrap(err, "token update failed")
}
