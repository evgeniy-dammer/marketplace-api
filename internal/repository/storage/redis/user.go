package redis

import (
	"fmt"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// UserGetOne select user by id from database.
func (r *Repository) UserGetOne(ctxr context.Context, userID string) (user.User, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.CategoryGetAll")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	var usr user.User

	keyCache := fmt.Sprintf("%s-%s: %s", "user", "id", userID)
	err := r.cache.Get(ctx, keyCache, usr)
	if err != nil {
		return usr, errors.Wrap(err, "postCacheRepository.GetPostByID.redisClient.Get")
	}

	/*if err = json.Unmarshal([]byte(postBytes), &user); err != nil {
		return user, errors.Wrap(err, "postCacheRepository.GetPostByID.json.Unmarshal")
	}*/

	return usr, nil
}

// UserCreate insert user into database.
func (r *Repository) UserCreate(ctxr context.Context, userID string, input user.CreateUserInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.CategoryGetAll")
		defer span.Finish()

		_ = context.New(ctxt)
	}

	return nil
}
