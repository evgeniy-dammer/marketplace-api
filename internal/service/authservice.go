package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

const (
	signingKey      = "jndc83uhf3fdoenfoe0iededf9uecijnuche9cc"
	tokenTTL        = 30 * time.Minute
	refreshTokenTTL = 12 * time.Hour
	valuesNum       = 6
)

var (
	errInvalidHash          = errors.New("the encoded hash is not in the correct format")
	errIncompatibleVersion  = errors.New("incompatible version of argon2")
	errInvalidPassword      = errors.New("invalid password")
	errInvalidSigningMethod = errors.New("invalid signing method")
	errInvalidTokenClaims   = errors.New("token claims are not of type *tokenClaims")
	params                  = &model.HashParams{Memory: 4096, Iterations: 3, Parallelism: 1, SaltLength: 16, KeyLength: 32}
)

// AuthService is an authentication service.
type AuthService struct {
	repo repository.Authorization
}

// NewAuthService constructor for AuthService.
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// GenerateToken generates authorization token.
func (s *AuthService) GenerateToken(userID string, username string, password string) (model.User, model.Tokens, error) {
	var user model.User

	var tokens model.Tokens

	var err error

	user, err = s.repo.GetUser(userID, username)

	if err != nil {
		return user, tokens, errors.Wrap(err, "can not get user")
	}

	if username != "" {
		match, err := comparePasswordAndHash(password, user.Password)
		if err != nil {
			return user, tokens, err
		}

		if !match {
			return user, tokens, errInvalidPassword
		}

		userID = user.ID
	}

	user.Password = ""
	user.RoleID = 0

	issuedAt := time.Now().Unix()
	expiresAt := time.Now().Add(tokenTTL).Unix()
	refreshExpiresAt := time.Now().Add(refreshTokenTTL).Unix()

	token := createNewToken(userID, expiresAt, issuedAt)
	tokens.AccessToken, err = token.SignedString([]byte(signingKey))

	if err != nil {
		return user, tokens, errors.Wrap(err, "can not get access token")
	}

	tokens.AccessTokenExpires = expiresAt
	tokens.TokenType = "Bearer"

	// create refresh token
	refreshToken := createNewToken(userID, refreshExpiresAt, issuedAt)
	tokens.RefreshToken, err = refreshToken.SignedString([]byte(signingKey))

	if err != nil {
		return user, tokens, errors.Wrap(err, "can not get refresh token")
	}

	return user, tokens, nil
}

// createNewToken creates new token with claims.
func createNewToken(userID string, expiresAt int64, issuedAt int64) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &model.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  issuedAt,
		},
		UserID: userID,
	})
}

// ParseToken checks access token and returns user id.
func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &model.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidSigningMethod
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", errors.Wrap(err, "can not parse token")
	}

	claims, ok := token.Claims.(*model.TokenClaims)

	if !ok {
		return "", errInvalidTokenClaims
	}

	return claims.UserID, nil
}

// CreateUser hashes the password and insert User into system.
func (s *AuthService) CreateUser(user model.User) (string, error) {
	pass, err := generatePasswordHash(user.Password, params)
	if err != nil {
		return "", err
	}

	user.Password = pass

	userID, err := s.repo.CreateUser(user)

	return userID, errors.Wrap(err, "can not create user")
}

// GetUserRole returns users role name.
func (s *AuthService) GetUserRole(id string) (string, error) {
	role, err := s.repo.GetUserRole(id)

	return role, errors.Wrap(err, "can not get role")
}

// generatePasswordHash hashes the password.
func generatePasswordHash(password string, params *model.HashParams) (string, error) {
	salt, err := generateRandomBytes(params.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	bytes := make([]byte, n)

	if _, err := rand.Read(bytes); err != nil {
		return nil, errors.Wrap(err, "can not read bytes")
	}

	return bytes, nil
}

func comparePasswordAndHash(password, encodedHash string) (bool, error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	param, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, param.Iterations, param.Memory, param.Parallelism, param.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

func decodeHash(encodedHash string) (*model.HashParams, []byte, []byte, error) {
	values := strings.Split(encodedHash, "$")

	if len(values) != valuesNum {
		return nil, nil, nil, errInvalidHash
	}

	var version int

	_, err := fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can not scan value")
	}

	if version != argon2.Version {
		return nil, nil, nil, errIncompatibleVersion
	}

	params = &model.HashParams{}

	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)

	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can not scan value")
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can not decode string")
	}

	params.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can not decode string")
	}

	params.KeyLength = uint32(len(hash))

	return params, salt, hash, nil
}
