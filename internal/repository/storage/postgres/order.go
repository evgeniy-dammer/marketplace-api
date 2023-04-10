package postgres

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/order"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// OrderGetAll selects all orders from database.
func (r *Repository) OrderGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]order.Order, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.OrderGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var orders []order.Order

	qry, args, err := r.orderGetAllQuery(meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build a query string")
	}

	err = r.databaseSlave.SelectContext(ctx, &orders, qry, args...)

	return orders, errors.Wrap(err, "orders select query error")
}

// orderGetAllQuery creates sql query.
func (r *Repository) orderGetAllQuery(meta query.MetaData, params queryparameter.QueryParameter) (string, []interface{}, error) { //nolint:lll
	builder := r.genSQL.Select(
		"id", "user_id", "organization_id", "table_id", "status_id", "totalsum", "created_at").
		From(orderTable).
		Where(squirrel.Eq{"is_deleted": false})

	switch {
	case !params.StartDate.IsZero() && params.EndDate.IsZero():
		builder = builder.Where(squirrel.And{
			squirrel.GtOrEq{"created_at": params.StartDate.Format("2006-01-02 15:04:05")},
			squirrel.LtOrEq{"created_at": time.Now().Format("2006-01-02 15:04:05")},
		})
	case params.StartDate.IsZero() && !params.EndDate.IsZero():
		builder = builder.Where(squirrel.And{
			squirrel.GtOrEq{"created_at": time.Now().Format("2006-01-02 15:04:05")},
			squirrel.LtOrEq{"created_at": params.EndDate.Format("2006-01-02 15:04:05")},
		})
	case !params.StartDate.IsZero() && !params.EndDate.IsZero():
		builder = builder.Where(squirrel.And{
			squirrel.GtOrEq{"created_at": params.StartDate.Format("2006-01-02 15:04:05")},
			squirrel.LtOrEq{"created_at": params.EndDate.Format("2006-01-02 15:04:05")},
		})
	}

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	if len(params.Sorts) > 0 {
		builder = builder.OrderBy(params.Sorts.Parsing(mappingSortOrder)...)
	} else {
		builder = builder.OrderBy("created_at DESC")
	}

	if params.Pagination.Limit > 0 {
		builder = builder.Limit(params.Pagination.Limit)
	}

	if params.Pagination.Offset > 0 {
		builder = builder.Offset(params.Pagination.Offset)
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", nil, errors.Wrap(err, "unable to build a query string")
	}

	return qry, args, nil
}

// OrderGetOne select order by id from database.
func (r *Repository) OrderGetOne(ctxr context.Context, meta query.MetaData, orderID string) (order.Order, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.OrderGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var ordr order.Order

	builder := r.genSQL.Select(
		"id", "user_id", "organization_id", "table_id", "status_id", "totalsum", "created_at").
		From(orderTable).
		Where(squirrel.Eq{"is_deleted": false, "id": orderID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return ordr, errors.Wrap(err, "unable to build a query string")
	}

	err = r.databaseSlave.GetContext(ctx, &ordr, qry, args...)

	return ordr, errors.Wrap(err, "order select query error")
}

// OrderCreate insert order into database.
func (r *Repository) OrderCreate(ctxr context.Context, meta query.MetaData, input order.CreateOrderInput) (string, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.OrderCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var orderID string

	trx, err := r.databaseMaster.Begin()
	if err != nil {
		return "", errors.Wrap(err, "transaction begin error")
	}

	builder := r.genSQL.Insert(orderTable).
		Columns("user_id", "organization_id", "table_id", "status_id", "totalsum", "user_created").
		Values(input.UserID, input.OrganizationID, input.TableID, input.StatusID, input.TotalSum, meta.UserID).
		Suffix("RETURNING \"id\"")

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	row := trx.QueryRowContext(ctx, qry, args...)

	if err = row.Scan(&orderID); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "orders rollback error")
		}

		return "", errors.Wrap(err, "order id scan error")
	}

	for _, item := range input.Items {
		builderOrderItem := r.genSQL.Insert(orderItemTable).
			Columns("order_id", "item_id", "quantity", "unitprise", "totalprice").
			Values(orderID, item.ItemID, item.Quantity, item.UnitPrice, item.TotalPrice).
			Suffix("RETURNING \"id\"")

		createOrderItemQuery, argsOrderItem, err := builderOrderItem.ToSql()
		if err != nil {
			return "", errors.Wrap(err, "unable to build a query string")
		}

		_, err = trx.ExecContext(ctx, createOrderItemQuery, argsOrderItem...)

		if err != nil {
			if err = trx.Rollback(); err != nil {
				return "", errors.Wrap(err, "orders_items table rollback error")
			}

			return "", errors.Wrap(err, "item insert query error")
		}
	}

	return orderID, trx.Commit()
}

// OrderUpdate updates order by id in database.
func (r *Repository) OrderUpdate(ctxr context.Context, meta query.MetaData, input order.UpdateOrderInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.OrderUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	trx, err := r.databaseMaster.Begin()
	if err != nil {
		return errors.Wrap(err, "transaction begin error")
	}

	builder := r.genSQL.Update(orderTable)

	if input.TableID != nil {
		builder = builder.Set("table_id", *input.TableID)
	}

	if input.Status != nil {
		builder = builder.Set("status_id", *input.Status)
	}

	if input.OrganizationID != nil {
		builder = builder.Set("organization_id", *input.OrganizationID)
	}

	if input.TotalSum != nil {
		builder = builder.Set("totalsum", *input.TotalSum)
	}

	builder = builder.Set("user_updated", meta.UserID).
		Set("updated_at", time.Now().UTC()).
		Where(squirrel.Eq{"is_deleted": false, "id": *input.ID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	if _, err = trx.ExecContext(ctx, qry, args...); err != nil {
		if err = trx.Rollback(); err != nil {
			return errors.Wrap(err, "orders rollback error")
		}

		return errors.Wrap(err, "order update error")
	}

	builderDeleteOrderItems := r.genSQL.Delete(orderItemTable).
		Where(squirrel.Eq{"order_id": *input.ID})

	deleteOrderItemsQuery, argsDeleteOrderItemsQuery, err := builderDeleteOrderItems.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	if _, err = trx.ExecContext(ctx, deleteOrderItemsQuery, argsDeleteOrderItemsQuery...); err != nil {
		if err = trx.Rollback(); err != nil {
			return errors.Wrap(err, "order items rollback error")
		}

		return errors.Wrap(err, "order items delete error")
	}

	for _, item := range *input.Items {
		builderCreateOrderItemQuery := r.genSQL.Insert(orderItemTable).
			Columns("order_id", "item_id", "quantity", "unitprise", "totalprice").
			Values(*input.ID, item.ItemID, item.Quantity, item.UnitPrice, item.TotalPrice)

		createOrderItemQuery, argsCreateOrderItemQuery, err := builderCreateOrderItemQuery.ToSql()
		if err != nil {
			return errors.Wrap(err, "unable to build a query string")
		}

		_, err = trx.ExecContext(ctx, createOrderItemQuery, argsCreateOrderItemQuery...)

		if err != nil {
			if err = trx.Rollback(); err != nil {
				return errors.Wrap(err, "orders_items table rollback error")
			}

			return errors.Wrap(err, "item insert query error")
		}
	}

	return trx.Commit()
}

// OrderDelete deletes order by id from database.
func (r *Repository) OrderDelete(ctxr context.Context, meta query.MetaData, orderID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.OrderDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(orderTable).
		Set("is_deleted", true).
		Set("user_deleted", meta.UserID).
		Set("deleted_at", time.Now().UTC()).
		Where(squirrel.Eq{"is_deleted": false, "id": orderID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.databaseMaster.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "order delete query error")
}
