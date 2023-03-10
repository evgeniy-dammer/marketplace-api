package table

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/table"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// TableGetAll returns all tables from the system.
func (s *UseCase) TableGetAll(ctx context.Context, userID string, organizationID string) ([]table.Table, error) {
	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.TableGetAll")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	tables, err := s.adapterStorage.TableGetAll(ctx, userID, organizationID)

	return tables, errors.Wrap(err, "tables select error")
}

// TableGetOne returns table by id from the system.
func (s *UseCase) TableGetOne(ctx context.Context, userID string, organizationID string, tableID string) (table.Table, error) {
	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.TableGetOne")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	tbl, err := s.adapterStorage.TableGetOne(ctx, userID, organizationID, tableID)

	return tbl, errors.Wrap(err, "table select error")
}

// TableCreate inserts table into system.
func (s *UseCase) TableCreate(ctx context.Context, userID string, input table.CreateTableInput) (string, error) {
	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.TableCreate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	tableID, err := s.adapterStorage.TableCreate(ctx, userID, input)

	return tableID, errors.Wrap(err, "table create error")
}

// TableUpdate updates table by id in the system.
func (s *UseCase) TableUpdate(ctx context.Context, userID string, input table.UpdateTableInput) error {
	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.TableUpdate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.TableUpdate(ctx, userID, input)

	return errors.Wrap(err, "table update error")
}

// TableDelete deletes table by id from the system.
func (s *UseCase) TableDelete(ctx context.Context, userID string, organizationID string, tableID string) error {
	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.TableDelete")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.TableDelete(ctx, userID, organizationID, tableID)

	return errors.Wrap(err, "table delete error")
}
