package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/order"
	"github.com/pkg/errors"
)

// OrderGetAll selects all orders from database.
func (r *Repository) OrderGetAll(userID string, organizationID string) ([]order.Order, error) {
	var orders []order.Order

	query := fmt.Sprintf(
		"SELECT id, user_id, organization_id, table_id, status_id, totalsum, created_at FROM %s "+
			"WHERE is_deleted = false AND organization_id = $1 ",
		orderTable,
	)

	err := r.db.Select(&orders, query, organizationID)

	return orders, errors.Wrap(err, "orders select query error")
}

// OrderGetOne select order by id from database.
func (r *Repository) OrderGetOne(userID string, organizationID string, orderID string) (order.Order, error) {
	var order order.Order

	query := fmt.Sprintf(
		"SELECT id, user_id, organization_id, table_id, status_id, totalsum, created_at FROM %s "+
			"WHERE is_deleted = false AND organization_id = $1 AND id = $2 ",
		orderTable,
	)
	err := r.db.Get(&order, query, organizationID, orderID)

	return order, errors.Wrap(err, "order select query error")
}

// OrderCreate insert order into database.
func (r *Repository) OrderCreate(userID string, order order.Order) (string, error) {
	var orderID string

	trx, err := r.db.Begin()
	if err != nil {
		return "", errors.Wrap(err, "transaction begin error")
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (user_id, organization_id, table_id, status_id, totalsum, user_created) "+
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		orderTable)
	row := trx.QueryRow(query, order.UserID, order.OrganizationID, order.TableID, order.StatusID, order.TotalSum, userID)

	if err = row.Scan(&orderID); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "orders rollback error")
		}

		return "", errors.Wrap(err, "order id scan error")
	}

	for _, item := range order.Items {
		createOrderItemQuery := fmt.Sprintf(
			"INSERT INTO %s (order_id, item_id, quantity, unitprise, totalprice) VALUES ($1, $2, $3, $4, $5)",
			orderItemTable,
		)

		_, err = trx.Exec(createOrderItemQuery, orderID, item.ItemID, item.Quantity, item.UnitPrice, item.TotalPrice)

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
func (r *Repository) OrderUpdate(userID string, input order.UpdateOrderInput) error {
	setValues := make([]string, 0, 6)
	args := make([]interface{}, 0, 6)
	argID := 1

	trx, err := r.db.Begin()
	if err != nil {
		return errors.Wrap(err, "transaction begin error")
	}

	if input.TableID != nil {
		setValues = append(setValues, fmt.Sprintf("table_id=$%d", argID))
		args = append(args, *input.TableID)
		argID++
	}

	if input.Status != nil {
		setValues = append(setValues, fmt.Sprintf("status_id=$%d", argID))
		args = append(args, *input.Status)
		argID++
	}

	if input.OrganizationID != nil {
		setValues = append(setValues, fmt.Sprintf("organization_id=$%d", argID))
		args = append(args, *input.OrganizationID)
		argID++
	}

	if input.TotalSum != nil {
		setValues = append(setValues, fmt.Sprintf("totalsum=$%d", argID))
		args = append(args, *input.TotalSum)
		argID++
	}

	setValues = append(setValues, fmt.Sprintf("user_updated=$%d", argID))
	args = append(args, userID)
	argID++

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE is_deleted = false AND organization_id = '%s' AND id = '%s'",
		orderTable, setQuery, *input.OrganizationID, *input.ID)

	if _, err = trx.Exec(query, args...); err != nil {
		if err = trx.Rollback(); err != nil {
			return errors.Wrap(err, "orders rollback error")
		}

		return errors.Wrap(err, "order update error")
	}

	deleteOrderItemsQuery := fmt.Sprintf("DELETE FROM  %s WHERE order_id = $1", orderItemTable)

	if _, err = trx.Exec(deleteOrderItemsQuery, *input.ID); err != nil {
		if err = trx.Rollback(); err != nil {
			return errors.Wrap(err, "order items rollback error")
		}

		return errors.Wrap(err, "order items delete error")
	}

	for _, item := range *input.Items {
		createOrderItemQuery := fmt.Sprintf(
			"INSERT INTO %s (order_id, item_id, quantity, unitprise, totalprice) VALUES ($1, $2, $3, $4, $5)",
			orderItemTable,
		)

		_, err = trx.Exec(createOrderItemQuery, *input.ID, item.ItemID, item.Quantity, item.UnitPrice, item.TotalPrice)

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
func (r *Repository) OrderDelete(userID string, organizationID string, orderID string) error {
	query := fmt.Sprintf(
		"UPDATE %s SET is_deleted = true, deleted_at = $1, user_deleted = $2 "+
			"WHERE is_deleted = false AND id = $3 AND organization_id = $4",
		orderTable,
	)

	_, err := r.db.Exec(query, time.Now().Format("2006-01-02 15:04:05"), userID, orderID, organizationID)

	return errors.Wrap(err, "order delete query error")
}
