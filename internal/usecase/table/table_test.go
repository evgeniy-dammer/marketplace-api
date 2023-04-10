package table

import (
	"os"
	"testing"

	"github.com/evgeniy-dammer/marketplace-api/internal/domain/role"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/table"
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
	storageRepo       = new(mockStorage.Table)
	cacheRepo         = new(mockCache.Table)
	isTracingOn       = true
	isCacheOn         = false
	ucDialog          *UseCase
	ucDialogWithCache *UseCase
	createTables      []table.CreateTableInput
	createTableWrong  table.CreateTableInput
	updateTables      []table.UpdateTableInput
	updateTableWrong  table.UpdateTableInput
	tables            []table.Table
	roles             []role.Role
	tble              table.Table
)

func TestMain(m *testing.M) {
	tableID := "49c9b955-8511-4b53-81ef-82e3d0259fed"
	wrongTableID := "49c9b955-8511-4b53-81ef-82e32d059fed"
	name := "Name1"
	organizationID := "1c6886a9-419e-4880-90cb-57da57d03155"

	createTable := table.CreateTableInput{
		Name:           "Name",
		OrganizationID: "1c6886a9-419e-4880-90cb-57da57d03155",
	}

	createTableWrong = table.CreateTableInput{
		Name:           "Name1",
		OrganizationID: "794f9bbd-fb40-44f2-b30b-5b82d3375ddb",
	}

	updateTable := table.UpdateTableInput{
		ID:             &tableID,
		Name:           &name,
		OrganizationID: &organizationID,
	}

	updateTableWrong = table.UpdateTableInput{
		ID:             &wrongTableID,
		Name:           &name,
		OrganizationID: &organizationID,
	}

	tble = table.Table{
		ID:             "49c9b955-8511-4b53-81ef-82e3d0259fed",
		Name:           "Name",
		OrganizationID: "1c6886a9-419e-4880-90cb-57da57d03155",
	}

	createTables = append(createTables, createTable)
	updateTables = append(updateTables, updateTable)
	tables = append(tables, tble)

	_ = logger.InitLogger()

	os.Exit(m.Run())
}

func initTestUseCaseTable(t *testing.T) {
	storageRepo.On("New",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("bool"),
		mock.AnythingOfType("bool")).
		Return(func(store storage.Table, cache cache.Table, isTracingOn bool, isCacheOn bool) *UseCase {
			return &UseCase{
				adapterStorage: store,
				adapterCache:   cache,
				isTracingOn:    isTracingOn,
				isCacheOn:      isCacheOn,
			}
		})

	storageRepo.On("TableGetAll",
		mock.Anything,
		mock.Anything,
		mock.Anything).
		Return(func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) []table.Table {
			if params.Search == errorString {
				return nil
			}

			return tables
		}, func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) error {
			if params.Search == errorString {
				return usecase.ErrTablesNotFound
			}

			return nil
		})

	storageRepo.On("getAllWithCache",
		mock.Anything,
		mock.Anything,
		mock.Anything).
		Return(func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) []table.Table {
			if params.Search == errorString {
				return nil
			}

			return tables
		}, func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) error {
			if params.Search == errorString {
				return usecase.ErrTablesNotFound
			}

			return nil
		})

	cacheRepo.On("TableGetAll",
		mock.Anything,
		mock.Anything,
		mock.Anything).
		Return(func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) []table.Table {
			if params.Search == errorString {
				return nil
			}

			return tables
		}, func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) error {
			if params.Search == errorString {
				return usecase.ErrTablesNotFound
			}

			return nil
		})

	cacheRepo.On("TableSetAll",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything).
		Return(func(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, tables []table.Table) error { //nolint:lll
			if params.Search == errorString {
				return usecase.ErrTablesNotFound
			}

			return nil
		})

	storageRepo.On("TableGetOne",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("string")).
		Return(func(ctx context.Context, meta query.MetaData, tableID string) table.Table {
			for k, v := range tables {
				if v.ID == tableID {
					return tables[k]
				}
			}

			return table.Table{}
		}, func(ctx context.Context, meta query.MetaData, tableID string) error {
			for _, v := range tables {
				if v.ID == tableID {
					return nil
				}
			}

			return usecase.ErrTableNotFound
		})

	storageRepo.On("getOneWithCache",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("string")).
		Return(func(ctx context.Context, meta query.MetaData, tableID string) table.Table {
			for k, v := range tables {
				if v.ID == tableID {
					return tables[k]
				}
			}

			return table.Table{}
		}, func(ctx context.Context, meta query.MetaData, tableID string) error {
			for _, v := range tables {
				if v.ID == tableID {
					return nil
				}
			}

			return usecase.ErrTableNotFound
		})

	cacheRepo.On("TableGetOne",
		mock.Anything,
		mock.AnythingOfType("string")).
		Return(func(ctx context.Context, tableID string) table.Table {
			for k, v := range tables {
				if v.ID == tableID {
					return tables[k]
				}
			}

			return table.Table{}
		}, func(ctx context.Context, tableID string) error {
			for _, v := range tables {
				if v.ID == tableID {
					return nil
				}
			}

			return usecase.ErrTableNotFound
		})

	storageRepo.On("TableCreate",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("table.CreateTableInput")).
		Return(func(ctx context.Context, meta query.MetaData, input table.CreateTableInput) string {
			for k, v := range tables {
				if input.OrganizationID == v.OrganizationID {
					return tables[k].ID
				}
			}

			return usecase.ErrTableNotFound.Error()
		}, func(ctx context.Context, meta query.MetaData, input table.CreateTableInput) error {
			for _, v := range tables {
				if input.OrganizationID != v.OrganizationID {
					return usecase.ErrTableNotFound
				}
			}

			return nil
		})

	cacheRepo.On("TableCreate",
		mock.Anything,
		mock.AnythingOfType("table.Table")).
		Return(func(ctx context.Context, input table.Table) error {
			for _, v := range tables {
				if input.OrganizationID != v.OrganizationID {
					return usecase.ErrTableNotFound
				}
			}

			return nil
		})

	cacheRepo.On("TableInvalidate",
		mock.Anything).
		Return(func(ctx context.Context) error {
			if len(tables) > 0 {
				return nil
			}

			return usecase.ErrTablesNotFound
		})

	storageRepo.On("TableUpdate",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("table.UpdateTableInput")).
		Return(func(ctx context.Context, meta query.MetaData, input table.UpdateTableInput) error {
			for _, v := range tables {
				if *input.ID == v.ID {
					v.Name = *input.Name
					v.OrganizationID = *input.OrganizationID

					return nil
				}
			}

			return usecase.ErrTableNotFound
		})

	cacheRepo.On("TableUpdate",
		mock.Anything,
		mock.AnythingOfType("table.Table")).
		Return(func(ctx context.Context, input table.Table) error {
			for _, v := range tables {
				if input.ID == v.ID {
					v.Name = input.Name
					v.OrganizationID = input.OrganizationID

					return nil
				}
			}

			return usecase.ErrTableNotFound
		})

	storageRepo.On("TableDelete",
		mock.Anything,
		mock.Anything,
		mock.AnythingOfType("string")).
		Return(func(ctx context.Context, meta query.MetaData, dTableID string) error {
			for k, v := range tables {
				if dTableID == v.ID {
					tables = append(tables[:k], tables[k+1:]...)

					return nil
				}
			}

			return usecase.ErrTableNotFound
		})

	cacheRepo.On("TableDelete",
		mock.Anything,
		mock.AnythingOfType("string")).
		Return(func(ctx context.Context, dTableID string) error {
			for k, v := range tables {
				if dTableID == v.ID {
					tables = append(tables[:k], tables[k+1:]...)

					return nil
				}
			}

			return usecase.ErrTableNotFound
		})

	storageRepo.On("TableGetAllRoles",
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

	cacheRepo.On("TableGetAllRoles",
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

func TestTable(t *testing.T) {
	initTestUseCaseTable(t)

	ucDialog = New(storageRepo, cacheRepo, isTracingOn, isCacheOn)
	ucDialogWithCache = New(storageRepo, cacheRepo, isTracingOn, true)
	assertion := assert.New(t)

	meta := query.MetaData{}
	params := queryparameter.QueryParameter{}

	t.Run("New", func(t *testing.T) {
		result := New(storageRepo, cacheRepo, isTracingOn, isCacheOn)
		assertion.Equal(result, ucDialog)
	})

	t.Run("TableGetAll", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.TableGetAll(ctx, meta, params)
		assertion.NoError(err)
		assertion.Equal(result, tables)
	})

	t.Run("TableGetAllWithError", func(t *testing.T) {
		ctx := context.Empty()

		var nilSlice []table.Table
		params.Search = errorString

		result, err := ucDialog.TableGetAll(ctx, meta, params)
		assertion.Error(err)
		assertion.Equal(result, nilSlice)

		params.Search = ""
	})

	t.Run("TableGetAllWithCache", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialogWithCache.TableGetAll(ctx, meta, params)
		assertion.NoError(err)
		assertion.Equal(result, tables)
	})

	t.Run("getAllWithCache", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.getAllWithCache(ctx, meta, params)
		assertion.NoError(err)
		assertion.Equal(result, tables)
	})

	t.Run("getAllWithCacheWithError", func(t *testing.T) {
		ctx := context.Empty()

		var nilSlice []table.Table
		params.Search = errorString

		result, err := ucDialog.getAllWithCache(ctx, meta, params)
		assertion.Error(err)
		assertion.Equal(result, nilSlice)

		params.Search = ""
	})

	t.Run("TableGetOne", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.TableGetOne(ctx, meta, tables[0].ID)
		assertion.NoError(err)
		assertion.Equal(result, tables[0])
	})

	t.Run("TableGetOneWithError", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.TableGetOne(ctx, meta, "")
		assertion.Error(err)
		assertion.Equal(result, table.Table{})
	})

	t.Run("TableGetOneWithCache", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialogWithCache.TableGetOne(ctx, meta, tables[0].ID)
		assertion.NoError(err)
		assertion.Equal(result, tables[0])
	})

	t.Run("getOneWithCache", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.getOneWithCache(ctx, meta, tables[0].ID)
		assertion.NoError(err)
		assertion.Equal(result, tables[0])
	})

	t.Run("getOneWithCacheWithError", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.getOneWithCache(ctx, meta, "")
		assertion.Error(err)
		assertion.Equal(result, table.Table{})
	})

	t.Run("TableCreate", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.TableCreate(ctx, meta, createTables[0])
		assertion.NoError(err)
		assertion.Equal(result, tables[0].ID)
	})

	t.Run("TableCreateWithError", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialog.TableCreate(ctx, meta, createTableWrong)
		assertion.Error(err)
		assertion.Equal(result, "")
	})

	t.Run("TableCreateWithCache", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialogWithCache.TableCreate(ctx, meta, createTables[0])
		assertion.NoError(err)
		assertion.Equal(result, tables[0].ID)
	})

	t.Run("TableCreateWithCacheWithError", func(t *testing.T) {
		ctx := context.Empty()

		result, err := ucDialogWithCache.TableCreate(ctx, meta, createTableWrong)
		assertion.Error(err)
		assertion.Equal(result, "")
	})

	t.Run("TableUpdate", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialog.TableUpdate(ctx, meta, updateTables[0])
		assertion.NoError(err)
		assertion.Equal(err, nil)
	})

	t.Run("TableUpdateWithError", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialog.TableUpdate(ctx, meta, updateTableWrong)
		assertion.Error(err)
	})

	t.Run("TableUpdateWithCache", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialogWithCache.TableUpdate(ctx, meta, updateTables[0])
		assertion.NoError(err)
		assertion.Equal(err, nil)
	})

	t.Run("TableUpdateWithCacheWithError", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialogWithCache.TableUpdate(ctx, meta, updateTableWrong)
		assertion.Error(err)
	})

	t.Run("TableSetAll", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialogWithCache.adapterCache.TableSetAll(ctx, meta, params, tables)
		assertion.NoError(err)
		assertion.Equal(err, nil)
	})

	t.Run("TableSetAllWithError", func(t *testing.T) {
		ctx := context.Empty()

		params.Search = errorString

		err := ucDialogWithCache.adapterCache.TableSetAll(ctx, meta, params, tables)
		assertion.Error(err)

		params.Search = ""
	})

	t.Run("TableDelete", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialog.TableDelete(ctx, meta, tables[0].ID)
		assertion.NoError(err)
		assertion.Equal(err, nil)
	})

	t.Run("TableDeleteWithError", func(t *testing.T) {
		ctx := context.Empty()

		err := ucDialog.TableDelete(ctx, meta, "")
		assertion.Error(err)
	})

	t.Run("TableDeleteWithCache", func(t *testing.T) {
		ctx := context.Empty()
		tables = append(tables, tble, tble, tble)

		err := ucDialogWithCache.TableDelete(ctx, meta, tble.ID)
		assertion.NoError(err)
		assertion.Equal(err, nil)
	})

	t.Run("TableDeleteWithCacheWithError", func(t *testing.T) {
		ctx := context.Empty()

		tables = append(tables, tble)

		err := ucDialogWithCache.TableDelete(ctx, meta, *updateTableWrong.ID)
		assertion.Error(err)
	})

	t.Run("TableInvalidate", func(t *testing.T) {
		ctx := context.Empty()

		tables = append(tables, tble)

		err := ucDialogWithCache.adapterCache.TableInvalidate(ctx)
		assertion.NoError(err)
		assertion.Equal(err, nil)
	})

	t.Run("TableInvalidateWithError", func(t *testing.T) {
		ctx := context.Empty()

		tables = []table.Table{}

		err := ucDialogWithCache.adapterCache.TableInvalidate(ctx)
		assertion.Error(err)
	})
}
