package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
	"golang.org/x/crypto/argon2"
	"strings"
	"time"
)

const (
	signingKey      = "nfljsdflbvjsdlfvblsjfnv"
	tokenTTL        = 30 * time.Minute
	refreshTokenTTL = 12 * time.Hour
)

var (
	errInvalidHash         = errors.New("the encoded hash is not in the correct format")
	errIncompatibleVersion = errors.New("incompatible version of argon2")
	errInvalidPassword     = errors.New("invalid password")
	params                 = &model.HashParams{Memory: 4096, Iterations: 3, Parallelism: 1, SaltLength: 16, KeyLength: 32}
)

// AuthService is an authentication service
type AuthService struct {
	repo repository.Authorization
}

// NewAuthService constructor for AuthService
func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// GenerateToken generates authorization token
func (s *AuthService) GenerateToken(id string, username string, password string) (model.User, model.Tokens, error) {
	var user model.User
	var tokens model.Tokens
	var err error

	if username != "" && password != "" {
		// get user from database
		user, err = s.repo.GetUser("", username)

		if err != nil {
			return user, tokens, err
		}

		match, err := comparePasswordAndHash(password, user.Password)
		user.Password = ""
		user.RoleId = ""

		if err != nil {
			return user, tokens, err
		}

		if !match {
			return user, tokens, errInvalidPassword
		}

		id = user.Id
	} else {
		user, err = s.repo.GetUser(id, "")

		if err != nil {
			return user, tokens, err
		}
		user.Password = ""
		user.RoleId = ""
	}

	issuedAt := time.Now().Unix()
	expiresAt := time.Now().Add(tokenTTL).Unix()
	refreshExpiresAt := time.Now().Add(refreshTokenTTL).Unix()

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: refreshExpiresAt,
			IssuedAt:  issuedAt,
		},
		id,
	})

	tokens.AccessToken, err = token.SignedString([]byte(signingKey))

	if err != nil {
		return user, tokens, err
	}

	tokens.AccessTokenExpires = expiresAt
	tokens.TokenType = "Bearer"

	// create refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  issuedAt,
		},
		id,
	})

	tokens.RefreshToken, err = refreshToken.SignedString([]byte(signingKey))

	if err != nil {
		return user, tokens, err
	}

	return user, tokens, nil
}

// ParseToken checks access token and returns user id
func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &model.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*model.TokenClaims)

	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

// CreateUser hashes the password and insert User into system
func (s *AuthService) CreateUser(user model.User, statusId string) (string, error) {
	pass, err := generatePasswordHash(user.Password, params)

	if err != nil {
		return "", err
	}

	user.Password = pass

	return s.repo.CreateUser(user, statusId)
}

// GetUserRole returns users role name
func (s *AuthService) GetUserRole(id string) (string, error) {
	return s.repo.GetUserRole(id)
}

// generatePasswordHash hashes the password
func generatePasswordHash(password string, params *model.HashParams) (string, error) {
	salt, err := generateRandomBytes(params.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, params.Memory, params.Iterations, params.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func comparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *model.HashParams, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")

	if len(vals) != 6 {
		return nil, nil, nil, errInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errIncompatibleVersion
	}

	p = &model.HashParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}
