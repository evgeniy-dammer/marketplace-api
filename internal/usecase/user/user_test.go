package user

import (
	"os"
	"testing"

	"github.com/evgeniy-dammer/marketplace-api/internal/domain/role"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/user"
	mockStorage "github.com/evgeniy-dammer/marketplace-api/internal/repository/storage/mockpostgres"
	mockCache "github.com/evgeniy-dammer/marketplace-api/internal/repository/storage/mockredis"
	"github.com/evgeniy-dammer/marketplace-api/internal/usecase"
	"github.com/evgeniy-dammer/marketplace-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/marketplace-api/internal/usecase/adapters/storage"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/logger"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const errorString = "error"

var (
	storageRepo       = new(mockStorage.User)
	cacheRepo         = new(mockCache.User)
	isTracingOn       = true
	isCacheOn         = false
	ucDialog          *UseCase
	ucDialogWithCache *UseCase
	createUsers       []user.CreateUserInput
	createUserWrong   user.CreateUserInput
	updateUsers       []user.UpdateUserInput
	updateUserWrong   user.UpdateUserInput
	users             []user.User
	roles             []role.Role
	usr               user.User
)

func TestMain(m *testing.M) {
	userID := "49c9b955-8511-4b53-81ef-82e3d0259fed"
	idWrong := "49c9b955-8511-4b53-81ef-82e32d059fed"
	firstName := "FirstName1"
	lastName := "LastName1"
	password := "password1"

	createUser := user.CreateUserInput{
		Phone:     "123456789",
		Password:  "password",
		FirstName: "FirstName",
		LastName:  "LastName",
		RoleID:    1,
	}

	createUserWrong = user.CreateUserInput{
		Phone:     "89324987",
		Password:  "password",
		FirstName: "FirstName",
		LastName:  "LastName",
		RoleID:    1,
	}

	updateUser := user.UpdateUserInput{
		ID:        &userID,
		FirstName: &firstName,
		LastName:  &lastName,
		Password:  &password,
	}

	updateUserWrong = user.UpdateUserInput{
		ID:        &idWrong,
		FirstName: &firstName,
		LastName:  &lastName,
		Password:  &password,
	}

	usr = user.User{
		ID:        "49c9b955-8511-4b53-81ef-82e3d0259fed",
		Phone:     "123456789",
		Password:  "password",
		FirstName: "FirstName",
		LastName:  "LastName",
		RoleID:    1,
		RoleName:  "admin",
		Status:    "active",
	}

	rle := role.Role{
		ID:   userID,
		Name: firstName,
	}

	createUsers = append(createUsers, createUser)
	updateUsers = append(updateUsers, updateUser)
	users = append(users, usr)
	roles = append(roles, rle)

	_ = logger.InitLogger()

	os.Exit(m.Run())
}

func initTestUseCaseUser(t *testing.T) {
	storageRepo.On("New",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("bool"),
		mock.AnythingOfType("bool")).
		Return(func(store storage.User, cache cache.User, isTracingOn bool, isCacheOn bool) *UseCase {
			return &UseCase{
				adapterStorage: store,
				adapterCache:   cache,
				isTracingOn:    isTracingOn,
				isCacheOn:      isCacheOn,
			}
		})

	storageRepo.On("UserGetAll",
		mock.Anything,
		mock.Anything,
		mock.Anything).
		Return(func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) []user.User {
			if params.Search == errorString {
				return nil
			}

			return users
		}, func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) error {
			if params.Search == errorString {
				return usecase.ErrUsersNotFound
			}

			return nil
		})

	storageRepo.On("getAllWithCache",
		mock.Anything,
		mock.Anything,
		mock.Anything).
		Return(func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) []user.User {
			if params.Search == errorString {
				return nil
			}

			return users
		}, func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) error {
			if params.Search == errorString {
				return usecase.ErrUsersNotFound
			}

			return nil
		})

	cacheRepo.On("UserGetAll",
		mock.Anything,
		mock.Anything,
		mock.Anything).
		Return(func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) []user.User {
			if params.Search == errorString {
				return nil
			}

			return users
		}, func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) error {
			if params.Search == errorString {
				return usecase.ErrUsersNotFound
			}

			return nil
		})

	cacheRepo.On("UserSetAll",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything).
		Return(func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, users []user.User) error {
			if params.Search == errorString {
				return usecase.ErrUsersNotFound
			}

			return nil
		})

	storageRepo.On("UserGetOne",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("string")).
		Return(func(ctx context.Context, meta query.MetaData, userID string) user.User {
			for k, v := range users {
				if v.ID == userID {
					return users[k]
				}
			}

			return user.User{}
		}, func(ctx context.Context, meta query.MetaData, userID string) error {
			for _, v := range users {
				if v.ID == userID {
					return nil
				}
			}

			return usecase.ErrUserNotFound
		})

	storageRepo.On("getOneWithCache",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("string")).
		Return(func(ctx context.Context, meta query.MetaData, userID string) user.User {
			for k, v := range users {
				if v.ID == userID {
					return users[k]
				}
			}

			return user.User{}
		}, func(ctx context.Context, meta query.MetaData, userID string) error {
			for _, v := range users {
				if v.ID == userID {
					return nil
				}
			}

			return usecase.ErrUserNotFound
		})

	cacheRepo.On("UserGetOne",
		mock.Anything,
		mock.AnythingOfType("string")).
		Return(func(ctx context.Context, userID string) user.User {
			for k, v := range users {
				if v.ID == userID {
					return users[k]
				}
			}

			return user.User{}
		}, func(ctx context.Context, userID string) error {
			for _, v := range users {
				if v.ID == userID {
					return nil
				}
			}

			return usecase.ErrUserNotFound
		})

	storageRepo.On("UserCreate",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("user.CreateUserInput")).
		Return(func(ctx context.Context, meta query.MetaData, input user.CreateUserInput) string {
			for k, v := range users {
				if input.Phone == v.Phone {
					return users[k].ID
				}
			}

			return usecase.ErrUserNotFound.Error()
		}, func(ctx context.Context, meta query.MetaData, input user.CreateUserInput) error {
			for _, v := range users {
				if input.Phone != v.Phone {
					return usecase.ErrUserNotFound
				}
			}

			return nil
		})

	cacheRepo.On("UserCreate",
		mock.Anything,
		mock.AnythingOfType("user.User")).
		Return(func(ctx context.Context, input user.User) error {
			for _, v := range users {
				if input.Phone != v.Phone {
					return usecase.ErrUserNotFound
				}
			}

			return nil
		})

	cacheRepo.On("UserInvalidate",
		mock.Anything).
		Return(func(ctx context.Context) error {
			if len(users) > 0 {
				return nil
			}

			return usecase.ErrUsersNotFound
		})

	storageRepo.On("UserUpdate",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("user.UpdateUserInput")).
		Return(func(ctx context.Context, meta query.MetaData, input user.UpdateUserInput) error {
			for _, v := range users {
				if *input.ID == v.ID {
					v.FirstName = *input.FirstName
					v.LastName = *input.LastName
					v.Password = *input.Password

					return nil
				}
			}

			return usecase.ErrUserNotFound
		})

	cacheRepo.On("UserUpdate",
		mock.Anything,
		mock.AnythingOfType("user.User")).
		Return(func(ctx context.Context, input user.User) error {
			for _, v := range users {
				if input.ID == v.ID {
					v.FirstName = input.FirstName
					v.LastName = input.LastName
					v.Password = input.Password

					return nil
				}
			}

			return usecase.ErrUserNotFound
		})

	storageRepo.On("UserDelete",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("string")).
		Return(func(ctx context.Context, meta query.MetaData, dUserID string) error {
			for k, v := range users {
				if dUserID == v.ID {
					users = append(users[:k], users[k+1:]...)

					return nil
				}
			}

			return usecase.ErrUserNotFound
		})

	cacheRepo.On("UserDelete",
		mock.Anything,
		mock.AnythingOfType("string")).
		Return(func(ctx context.Context, dUserID string) error {
			for k, v := range users {
				if dUserID == v.ID {
					users = append(users[:k], users[k+1:]...)

					return nil
				}
			}

			return usecase.ErrUserNotFound
		})

	storageRepo.On("UserGetAllRoles",
		mock.Anything,
		mock.Anything,
		mock.Anything).
		Return(func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) []role.Role {
			if len(roles) > 0 {
				return roles
			}

			var nilSlice []role.Role

			return nilSlice
		}, func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) error {
			if len(roles) > 0 {
				return nil
			}

			return usecase.ErrRolesNotFound
		})

	storageRepo.On("getAllRolesWithCache",
		mock.Anything,
		mock.Anything,
		mock.Anything).
		Return(func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) []role.Role {
			if len(roles) > 0 {
				return roles
			}

			var nilSlice []role.Role

			return nilSlice
		}, func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) error {
			if len(roles) > 0 {
				return nil
			}

			return usecase.ErrRolesNotFound
		})

	cacheRepo.On("UserGetAllRoles",
		mock.Anything,
		mock.Anything,
		mock.Anything).
		Return(func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) []role.Role {
			if len(roles) > 0 {
				return roles
			}

			var nilSlice []role.Role

			return nilSlice
		}, func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) error {
			if len(roles) > 0 {
				return nil
			}

			return usecase.ErrRolesNotFound
		})
}

func TestUser(t *testing.T) {
	initTestUseCaseUser(t)

	ucDialog = New(storageRepo, cacheRepo, isTracingOn, isCacheOn)
	ucDialogWithCache = New(storageRepo, cacheRepo, isTracingOn, true)
	assertion := assert.New(t)

	meta := query.MetaData{}
	params := queryparameter.QueryParameter{}

	t.Run("New", func(t *testing.T) {
		result := New(storageRepo, cacheRepo, isTracingOn, isCacheOn)
		assertion.Equal(result, ucDialog)
	})

	t.Run("UserGetAll", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.UserGetAll(ctx, meta, params)
		assertion.NoError(err)
		assertion.Equal(result, users)
	})

	t.Run("UserGetAllWithError", func(t *testing.T) {
		ctx := context.Empty()

		var nilSlice []user.User
		params.Search = errorString

		result, err := ucDialog.UserGetAll(ctx, meta, params)
		assertion.Error(err)
		assertion.Equal(result, nilSlice)

		params.Search = ""
	})

	t.Run("UserGetAllWithCache", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialogWithCache.UserGetAll(ctx, meta, params)
		assertion.NoError(err)
		assertion.Equal(result, users)
	})

	t.Run("getAllWithCache", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.getAllWithCache(ctx, meta, params)
		assertion.NoError(err)
		assertion.Equal(result, users)
	})

	t.Run("getAllWithCacheWithError", func(t *testing.T) {
		ctx := context.Empty()

		var nilSlice []user.User
		params.Search = errorString

		result, err := ucDialog.getAllWithCache(ctx, meta, params)
		assertion.Error(err)
		assertion.Equal(result, nilSlice)

		params.Search = ""
	})

	t.Run("UserGetOne", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.UserGetOne(ctx, meta, users[0].ID)
		assertion.NoError(err)
		assertion.Equal(result, users[0])
	})

	t.Run("UserGetOneWithError", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.UserGetOne(ctx, meta, "")
		assertion.Error(err)
		assertion.Equal(result, user.User{})
	})

	t.Run("UserGetOneWithCache", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialogWithCache.UserGetOne(ctx, meta, users[0].ID)
		assertion.NoError(err)
		assertion.Equal(result, users[0])
	})

	t.Run("getOneWithCache", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.getOneWithCache(ctx, meta, users[0].ID)
		assertion.NoError(err)
		assertion.Equal(result, users[0])
	})

	t.Run("getOneWithCacheWithError", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.getOneWithCache(ctx, meta, "")
		assertion.Error(err)
		assertion.Equal(result, user.User{})
	})

	t.Run("UserCreate", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.UserCreate(ctx, meta, createUsers[0])
		assertion.NoError(err)
		assertion.Equal(result, users[0].ID)
	})

	t.Run("UserCreateWithError", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.UserCreate(ctx, meta, createUserWrong)
		assertion.Error(err)
		assertion.Equal(result, "")
	})

	t.Run("UserCreateWithCache", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialogWithCache.UserCreate(ctx, meta, createUsers[0])
		assertion.NoError(err)
		assertion.Equal(result, users[0].ID)
	})

	t.Run("UserCreateWithCacheWithError", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialogWithCache.UserCreate(ctx, meta, createUserWrong)
		assertion.Error(err)
		assertion.Equal(result, "")
	})

	t.Run("UserUpdate", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialog.UserUpdate(ctx, meta, updateUsers[0])
		assertion.NoError(err)
		assertion.Equal(err, nil)
	})

	t.Run("UserUpdateWithError", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialog.UserUpdate(ctx, meta, updateUserWrong)
		assertion.Error(err)
	})

	t.Run("UserUpdateWithCache", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialogWithCache.UserUpdate(ctx, meta, updateUsers[0])
		assertion.NoError(err)
		assertion.Equal(err, nil)
	})

	t.Run("UserUpdateWithCacheWithError", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialogWithCache.UserUpdate(ctx, meta, updateUserWrong)
		assertion.Error(err)
	})

	t.Run("UserSetAll", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialogWithCache.adapterCache.UserSetAll(ctx, meta, params, users)
		assertion.NoError(err)
		assertion.Equal(err, nil)
	})

	t.Run("UserSetAllWithError", func(t *testing.T) {
		ctx := context.Empty()

		params.Search = errorString

		err := ucDialogWithCache.adapterCache.UserSetAll(ctx, meta, params, users)
		assertion.Error(err)

		params.Search = ""
	})

	t.Run("UserDelete", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialog.UserDelete(ctx, meta, users[0].ID)
		assertion.NoError(err)
		assertion.Equal(err, nil)
	})

	t.Run("UserDeleteWithError", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialog.UserDelete(ctx, meta, "")
		assertion.Error(err)
	})

	t.Run("UserDeleteWithCache", func(t *testing.T) {
		ctx := context.Empty()
		users = append(users, usr, usr, usr)

		err := ucDialogWithCache.UserDelete(ctx, meta, usr.ID)
		assertion.NoError(err)
		assertion.Equal(err, nil)
	})

	t.Run("UserDeleteWithCacheWithError", func(t *testing.T) {
		ctx := context.Empty()

		users = append(users, usr)

		err := ucDialogWithCache.UserDelete(ctx, meta, *updateUserWrong.ID)
		assertion.Error(err)
	})

	t.Run("UserGetAllRoles", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.UserGetAllRoles(ctx, meta, params)
		assertion.NoError(err)
		assertion.Equal(result, roles)
	})

	t.Run("UserGetAllRolesWithCache", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialogWithCache.UserGetAllRoles(ctx, meta, params)
		assertion.NoError(err)
		assertion.Equal(result, roles)
	})

	t.Run("getAllRolesWithCache", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.getAllRolesWithCache(ctx, meta, params)
		assertion.NoError(err)
		assertion.Equal(result, roles)
	})

	t.Run("getAllRolesWithCacheWithError", func(t *testing.T) {
		ctx := context.Empty()

		roles = []role.Role{}

		var nilSlice []role.Role

		result, err := ucDialog.getAllRolesWithCache(ctx, meta, params)
		assertion.Error(err)
		assertion.Equal(result, nilSlice)
	})

	t.Run("UserInvalidate", func(t *testing.T) {
		ctx := context.Empty()

		users = append(users, usr)

		err := ucDialogWithCache.adapterCache.UserInvalidate(ctx)
		assertion.NoError(err)
		assertion.Equal(err, nil)
	})

	t.Run("UserInvalidateWithError", func(t *testing.T) {
		ctx := context.Empty()

		users = []user.User{}

		err := ucDialogWithCache.adapterCache.UserInvalidate(ctx)
		assertion.Error(err)
	})
}
