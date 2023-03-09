package redis

import (
	"context"
	"fmt"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/pkg/errors"
)

// UserGetOne select user by id from database.
func (r *Repository) UserGetOne(userID string) (user.User, error) {
	var user user.User

	keyCache := fmt.Sprintf("%s-%s: %s", "user", "id", userID)
	err := r.cache.Get(context.Background(), keyCache, user)
	if err != nil {
		return user, errors.Wrap(err, "postCacheRepository.GetPostByID.redisClient.Get")
	}

	/*if err = json.Unmarshal([]byte(postBytes), &user); err != nil {
		return user, errors.Wrap(err, "postCacheRepository.GetPostByID.json.Unmarshal")
	}*/

	return user, nil
}

// UserCreate insert user into database.
func (r *Repository) UserCreate(userID string, user user.User) error {
	return nil
}
