package usecase

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/evgeniy-dammer/marketplace-api/internal/domain/token"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

const (
	ValuesNum = 6
)

// HashParams is a password hash params.
type HashParams struct {
	Memory      uint32
	Iterations  uint32
	SaltLength  uint32
	KeyLength   uint32
	Parallelism uint8
}

var Params = &HashParams{Memory: 4096, Iterations: 3, Parallelism: 1, SaltLength: 16, KeyLength: 32}

// GeneratePasswordHash hashes the password.
func GeneratePasswordHash(password string, Params *HashParams) (string, error) {
	salt, err := GenerateRandomBytes(Params.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, Params.Iterations, Params.Memory, Params.Parallelism, Params.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		Params.Memory,
		Params.Iterations,
		Params.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func GenerateRandomBytes(n uint32) ([]byte, error) {
	bytes := make([]byte, n)

	if _, err := rand.Read(bytes); err != nil {
		return nil, errors.Wrap(err, "can not read bytes")
	}

	return bytes, nil
}

func ComparePasswordAndHash(password, encodedHash string) (bool, error) {
	param, salt, hash, err := DecodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, param.Iterations, param.Memory, param.Parallelism, param.KeyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

func DecodeHash(encodedHash string) (*HashParams, []byte, []byte, error) {
	values := strings.Split(encodedHash, "$")

	if len(values) != ValuesNum {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int

	_, err := fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can not scan value")
	}

	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	Params = &HashParams{}

	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &Params.Memory, &Params.Iterations, &Params.Parallelism)

	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can not scan value")
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can not decode string")
	}

	Params.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can not decode string")
	}

	Params.KeyLength = uint32(len(hash))

	return Params, salt, hash, nil
}

// CreateNewToken creates new token with claims.
func CreateNewToken(userID string, expiresAt int64, issuedAt int64, hash string) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &token.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  issuedAt,
		},
		UserID: userID,
		Hash:   hash,
	})
}
