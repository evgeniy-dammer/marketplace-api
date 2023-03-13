package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/item"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// ItemGetAll selects all items from database.
func (r *Repository) ItemGetAll(ctxr context.Context, userID string, organizationID string) ([]item.Item, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ItemGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var items []item.Item

	egroup := &errgroup.Group{}

	query := fmt.Sprintf(
		"SELECT "+
			"id, name_tm, name_ru, name_tr, name_en, description_tm, description_ru, description_tr, description_en, "+
			"internal_id, price, rating, comments_qty, category_id, organization_id, brand_id, created_at "+
			"FROM %s WHERE is_deleted = false AND organization_id = $1 ",
		itemTable,
	)

	err := r.database.SelectContext(ctx, &items, query, organizationID)

	for i := 0; i < len(items); i++ {
		index := i

		egroup.Go(func() error {
			queryImages := fmt.Sprintf(
				"SELECT id, object_id, type, origin, middle, small, organization_id, is_main FROM %s "+
					"WHERE is_main = true AND object_id = $1 ",
				imageTable,
			)

			err = r.database.SelectContext(ctx, &items[index].Images, queryImages, items[index].ID)

			return errors.Wrap(err, "images select query error")
		})

	}

	err = egroup.Wait()

	return items, errors.Wrap(err, "items select query error")
}

// ItemGetOne select item by id from database.
func (r *Repository) ItemGetOne(ctxr context.Context, userID string, organizationID string, itemID string) (item.Item, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ItemGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var itm item.Item

	egroup := &errgroup.Group{}

	query := fmt.Sprintf(
		"SELECT "+
			"id, name_tm, name_ru, name_tr, name_en, description_tm, description_ru, description_tr, description_en, "+
			"internal_id, price, rating, comments_qty, category_id, organization_id, brand_id, created_at "+
			"FROM %s WHERE is_deleted = false AND organization_id = $1 AND id = $2 ",
		itemTable,
	)

	err := r.database.GetContext(ctx, &itm, query, organizationID, itemID)
	if err != nil {
		return itm, errors.Wrap(err, "item select query error")
	}

	egroup.Go(func() error {
		queryImages := fmt.Sprintf(
			"SELECT id, object_id, type, origin, middle, small, organization_id, is_main FROM %s WHERE object_id = $1 ",
			imageTable,
		)

		err = r.database.SelectContext(ctx, &itm.Images, queryImages, itm.ID)

		return errors.Wrap(err, "images select query error")
	})

	egroup.Go(func() error {
		querySpecifications := fmt.Sprintf(
			"SELECT id, item_id, organization_id, name_tm, name_ru, name_tr, name_en, description_tm, description_ru, "+
				"description_tr, description_en, value FROM %s WHERE item_id = $1 ",
			specificationTable,
		)

		err = r.database.SelectContext(ctx, &itm.Specification, querySpecifications, itm.ID)

		return errors.Wrap(err, "specification select query error")
	})

	egroup.Go(func() error {
		queryComments := fmt.Sprintf(
			"SELECT id, item_id, organization_id, content, status_id, rating, user_created, created_at FROM %s "+
				"WHERE is_deleted = false AND item_id = $1 ",
			commentTable,
		)

		err = r.database.SelectContext(ctx, &itm.Comments, queryComments, itm.ID)

		return errors.Wrap(err, "comments select query error")
	})

	err = egroup.Wait()

	return itm, errors.Wrap(err, "images select query error")
}

// ItemCreate insert item into database.
func (r *Repository) ItemCreate(ctxr context.Context, userID string, input item.CreateItemInput) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ItemCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var itemID string

	query := fmt.Sprintf(
		"INSERT INTO %s "+
			"(name_tm, name_ru, name_tr, name_en,  description_tm, description_ru, description_tr, description_en, "+
			"internal_id, price, category_id, organization_id, brand_id, user_created) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id",
		itemTable,
	)

	row := r.database.QueryRowContext(
		ctx,
		query,
		input.NameTm,
		input.NameRu,
		input.NameTr,
		input.NameEn,
		input.DescriptionTm,
		input.DescriptionRu,
		input.DescriptionTr,
		input.DescriptionEn,
		input.InternalID,
		input.Price,
		input.CategoryID,
		input.OrganizationID,
		input.BrandID,
		userID,
	)

	err := row.Scan(&itemID)

	return itemID, errors.Wrap(err, "item create query error")
}

// ItemUpdate updates item by id in database.
func (r *Repository) ItemUpdate(ctxr context.Context, userID string, input item.UpdateItemInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ItemUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	setValues := make([]string, 0, 14)
	args := make([]interface{}, 0, 14)
	argID := 1

	if input.NameTm != nil {
		setValues = append(setValues, fmt.Sprintf("name_tm=$%d", argID))
		args = append(args, *input.NameTm)
		argID++
	}

	if input.NameRu != nil {
		setValues = append(setValues, fmt.Sprintf("name_ru=$%d", argID))
		args = append(args, *input.NameRu)
		argID++
	}

	if input.NameTr != nil {
		setValues = append(setValues, fmt.Sprintf("name_tr=$%d", argID))
		args = append(args, *input.NameTr)
		argID++
	}

	if input.NameEn != nil {
		setValues = append(setValues, fmt.Sprintf("name_en=$%d", argID))
		args = append(args, *input.NameEn)
		argID++
	}

	if input.DescriptionTm != nil {
		setValues = append(setValues, fmt.Sprintf("description_tm=$%d", argID))
		args = append(args, *input.DescriptionTm)
		argID++
	}

	if input.DescriptionRu != nil {
		setValues = append(setValues, fmt.Sprintf("description_ru=$%d", argID))
		args = append(args, *input.DescriptionRu)
		argID++
	}

	if input.DescriptionTr != nil {
		setValues = append(setValues, fmt.Sprintf("description_tr=$%d", argID))
		args = append(args, *input.DescriptionTr)
		argID++
	}

	if input.DescriptionEn != nil {
		setValues = append(setValues, fmt.Sprintf("description_en=$%d", argID))
		args = append(args, *input.DescriptionEn)
		argID++
	}

	if input.InternalID != nil {
		setValues = append(setValues, fmt.Sprintf("internal_id=$%d", argID))
		args = append(args, *input.InternalID)
		argID++
	}

	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argID))
		args = append(args, *input.Price)
		argID++
	}

	if input.CategoryID != nil {
		setValues = append(setValues, fmt.Sprintf("category_id=$%d", argID))
		args = append(args, *input.CategoryID)
		argID++
	}

	if input.OrganizationID != nil {
		setValues = append(setValues, fmt.Sprintf("organization_id=$%d", argID))
		args = append(args, *input.OrganizationID)
		argID++
	}

	if input.BrandID != nil {
		setValues = append(setValues, fmt.Sprintf("brand_id=$%d", argID))
		args = append(args, *input.BrandID)
		argID++
	}

	setValues = append(setValues, fmt.Sprintf("user_updated=$%d", argID))
	args = append(args, userID)
	argID++

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE is_deleted = false AND organization_id = '%s' AND id = '%s'",
		itemTable, setQuery, *input.OrganizationID, *input.ID)

	_, err := r.database.ExecContext(ctx, query, args...)

	return errors.Wrap(err, "item update query error")
}

// ItemDelete deletes item by id from database.
func (r *Repository) ItemDelete(ctxr context.Context, userID string, organizationID string, itemID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ItemDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	query := fmt.Sprintf(
		"UPDATE %s SET is_deleted = true, deleted_at = $1, user_deleted = $2 "+
			"WHERE is_deleted = false AND id = $3 AND organization_id = $4",
		itemTable,
	)

	_, err := r.database.ExecContext(ctx, query, time.Now().Format("2006-01-02 15:04:05"), userID, itemID, organizationID)

	return errors.Wrap(err, "item delete query error")
}
