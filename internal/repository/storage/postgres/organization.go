package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/organization"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// OrganizationGetAll selects all organizations from database.
func (r *Repository) OrganizationGetAll(ctxr context.Context, userID string) ([]organization.Organization, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.OrganizationGetAll")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	var organizations []organization.Organization

	query := fmt.Sprintf("SELECT id, name, user_id, address, phone FROM %s WHERE is_deleted = false AND user_id = $1",
		organizationTable)

	err := r.database.SelectContext(ctx, &organizations, query, userID)

	return organizations, errors.Wrap(err, "organizations select query error")
}

// OrganizationGetOne select organization by id from database.
func (r *Repository) OrganizationGetOne(ctxr context.Context, userID string, organizationID string) (organization.Organization, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.OrganizationGetOne")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	var org organization.Organization

	query := fmt.Sprintf(
		"SELECT id, name, user_id, address, phone FROM %s WHERE is_deleted = false AND user_id = $1 AND id = $2",
		organizationTable)
	err := r.database.GetContext(ctx, &org, query, userID, organizationID)

	return org, errors.Wrap(err, "organization select query error")
}

// OrganizationCreate insert organization into database.
func (r *Repository) OrganizationCreate(ctxr context.Context, userID string, input organization.CreateOrganizationInput) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.OrganizationCreate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	var organizationID string

	createUserQuery := fmt.Sprintf(
		"INSERT INTO %s (name, user_id, address, phone, user_created) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		organizationTable)

	row := r.database.QueryRowContext(ctx, createUserQuery, input.Name, userID, input.Address, input.Phone, userID)

	err := row.Scan(&organizationID)

	return organizationID, errors.Wrap(err, "organization create query error")
}

// OrganizationUpdate updates organization by id in database.
func (r *Repository) OrganizationUpdate(ctxr context.Context, userID string, input organization.UpdateOrganizationInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.OrganizationUpdate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	setValues := make([]string, 0, 4)
	args := make([]interface{}, 0, 4)
	argID := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argID))
		args = append(args, *input.Name)
		argID++
	}

	if input.Address != nil {
		setValues = append(setValues, fmt.Sprintf("address=$%d", argID))
		args = append(args, *input.Address)
		argID++
	}

	if input.Phone != nil {
		setValues = append(setValues, fmt.Sprintf("phone=$%d", argID))
		args = append(args, *input.Phone)
		argID++
	}

	setValues = append(setValues, fmt.Sprintf("user_updated=$%d", argID))
	args = append(args, userID)
	argID++

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE is_deleted = false AND id = '%s' AND user_id = '%s'",
		organizationTable, setQuery, *input.ID, userID)

	_, err := r.database.ExecContext(ctx, query, args...)

	return errors.Wrap(err, "organization update query error")
}

// OrganizationDelete deletes organization by id from database.
func (r *Repository) OrganizationDelete(ctxr context.Context, userID string, organizationID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.OrganizationDelete")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	query := fmt.Sprintf(
		"UPDATE %s SET is_deleted = true, deleted_at = $1, user_deleted = $2 WHERE is_deleted = false AND id = $3",
		organizationTable,
	)

	_, err := r.database.ExecContext(ctx, query, time.Now().Format("2006-01-02 15:04:05"), userID, organizationID)

	return errors.Wrap(err, "organization delete query error")
}
