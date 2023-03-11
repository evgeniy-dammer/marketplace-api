package authentication

import (
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/token"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/golang-jwt/jwt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// AuthenticationGenerateToken generates authorization token.
func (s *UseCase) AuthenticationGenerateToken(ctx context.Context, userID string, username string, password string) (user.User, token.Tokens, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.AuthenticationGenerateToken")
		defer span.Finish()

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
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.AuthenticationParseToken")
		defer span.Finish()

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
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.AuthenticationCreateUser")
		defer span.Finish()

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

// AuthenticationGetUserRole returns users role name.
func (s *UseCase) AuthenticationGetUserRole(ctx context.Context, userID string) (string, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.AuthenticationGetUserRole")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getUserRoleWithCache(ctx, s, userID)
	}

	role, err := s.adapterStorage.AuthenticationGetUserRole(ctx, userID)

	return role, errors.Wrap(err, "can not get role")
}

func getUserRoleWithCache(ctx context.Context, s *UseCase, userID string) (string, error) {
	role, err := s.adapterCache.AuthenticationGetUserRole(ctx, userID)
	if err != nil {
		logger.Logger.Error("unable to get user role from cache", zap.String("error", err.Error()))
	}

	if role != "" {
		return role, nil
	}

	role, err = s.adapterStorage.AuthenticationGetUserRole(ctx, userID)

	if err != nil {
		return role, errors.Wrap(err, "user role select failed")
	}

	if err = s.adapterCache.AuthenticationSetUserRole(ctx, userID, role); err != nil {
		logger.Logger.Error("unable to add user role into cache", zap.String("error", err.Error()))
	}

	return role, nil
}