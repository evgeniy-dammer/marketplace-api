package postgres

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/item"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// ItemGetAll selects all items from database.
func (r *Repository) ItemGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]item.Item, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ItemGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var items []item.Item

	egroup := &errgroup.Group{}

	qry, args, err := r.itemGetAllQuery(meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.SelectContext(ctx, &items, qry, args...)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build a query string")
	}

	for i := 0; i < len(items); i++ {
		index := i

		egroup.Go(func() error {
			builder := r.genSQL.Select(
				"id", "object_id", "type", "origin", "middle", "small", "organization_id", "is_main").
				From(imageTable).
				Where(squirrel.Eq{"is_main": true, "object_id": items[index].ID})

			qryImages, argsImages, err := builder.ToSql()
			if err != nil {
				return errors.Wrap(err, "unable to build a query string")
			}

			err = r.database.SelectContext(ctx, &items[index].Images, qryImages, argsImages...)

			return errors.Wrap(err, "images select query error")
		})
	}

	err = egroup.Wait()

	return items, errors.Wrap(err, "items select query error")
}

// itemGetAllQuery creates sql query.
func (r *Repository) itemGetAllQuery(meta query.MetaData, params queryparameter.QueryParameter) (string, []interface{}, error) { //nolint:lll
	builder := r.genSQL.Select("id", "name_tm", "name_ru", "name_tr", "name_en", "description_tm",
		"description_ru", "description_tr", "description_en", "internal_id", "price", "rating", "comments_qty",
		"category_id", "organization_id", "brand_id", "created_at").
		From(itemTable)

	if params.Search != "" {
		search := "%" + params.Search + "%"

		builder = builder.Where(squirrel.And{
			squirrel.Eq{"is_deleted": false},
			squirrel.Or{
				squirrel.Like{"name_tm": search},
				squirrel.Like{"name_ru": search},
				squirrel.Like{"name_tr": search},
				squirrel.Like{"name_en": search},
				squirrel.Like{"description_tm": search},
				squirrel.Like{"description_ru": search},
				squirrel.Like{"description_tr": search},
				squirrel.Like{"description_en": search},
				squirrel.Like{"internal_id": search},
				squirrel.Like{"price": search},
				squirrel.Like{"rating": search},
				squirrel.Like{"comments_qty": search},
				squirrel.Like{"category_id": search},
				squirrel.Like{"organization_id": search},
				squirrel.Like{"brand_id": search},
				squirrel.Like{"created_at": search},
			},
		})
	} else {
		builder = builder.Where(squirrel.Eq{"is_deleted": false})
	}

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
		builder = builder.OrderBy(params.Sorts.Parsing(mappingSortItem)...)
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

// ItemGetOne select item by id from database.
func (r *Repository) ItemGetOne(ctxr context.Context, meta query.MetaData, itemID string) (item.Item, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ItemGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var itm item.Item

	egroup := &errgroup.Group{}

	builder := r.genSQL.Select("id", "name_tm", "name_ru", "name_tr", "name_en", "description_tm",
		"description_ru", "description_tr", "description_en", "internal_id", "price", "rating", "comments_qty",
		"category_id", "organization_id", "brand_id", "created_at").
		From(itemTable).
		Where(squirrel.Eq{"is_deleted": false, "id": itemID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return itm, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.GetContext(ctx, &itm, qry, args...)
	if err != nil {
		return itm, errors.Wrap(err, "item select query error")
	}

	egroup.Go(func() error {
		builderImages := r.genSQL.Select(
			"id", "object_id", "type", "origin", "middle", "small", "organization_id", "is_main").
			From(imageTable).
			Where(squirrel.Eq{"object_id": itm.ID})

		qryImages, argsImages, err := builderImages.ToSql()
		if err != nil {
			return errors.Wrap(err, "unable to build a query string")
		}

		err = r.database.SelectContext(ctx, &itm.Images, qryImages, argsImages...)

		return errors.Wrap(err, "images select query error")
	})

	egroup.Go(func() error {
		builderSpecifications := r.genSQL.Select(
			"id", "item_id", "organization_id", "name_tm", "name_ru", "name_tr", "name_en",
			"description_tm", "description_ru", "description_tr", "description_en", "value").
			From(specificationTable).
			Where(squirrel.Eq{"item_id": itm.ID})

		qrySpecifications, argsSpecifications, err := builderSpecifications.ToSql()
		if err != nil {
			return errors.Wrap(err, "unable to build a query string")
		}

		err = r.database.SelectContext(ctx, &itm.Specification, qrySpecifications, argsSpecifications...)

		return errors.Wrap(err, "specification select query error")
	})

	egroup.Go(func() error {
		builderComments := r.genSQL.Select(
			"id", "item_id", "organization_id", "content", "status_id", "rating", "user_created", "created_at").
			From(specificationTable).
			Where(squirrel.Eq{"is_deleted": false, "item_id": itm.ID})

		qryComments, argsComments, err := builderComments.ToSql()
		if err != nil {
			return errors.Wrap(err, "unable to build a query string")
		}

		err = r.database.SelectContext(ctx, &itm.Comments, qryComments, argsComments...)

		return errors.Wrap(err, "comments select query error")
	})

	err = egroup.Wait()

	return itm, errors.Wrap(err, "images select query error")
}

// ItemCreate insert item into database.
func (r *Repository) ItemCreate(ctxr context.Context, meta query.MetaData, input item.CreateItemInput) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ItemCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var itemID string

	builder := r.genSQL.Insert(itemTable).
		Columns(
			"name_tm", "name_ru", "name_tr", "name_en", "description_tm", "description_ru", "description_tr",
			"description_en", "internal_id", "price", "category_id", "organization_id", "brand_id", "created_at").
		Values(
			input.NameTm, input.NameRu, input.NameTr, input.NameEn, input.DescriptionTm, input.DescriptionRu,
			input.DescriptionTr, input.DescriptionEn, input.InternalID, input.Price, input.CategoryID,
			input.OrganizationID, input.BrandID, meta.UserID).
		Suffix("RETURNING \"id\"")

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	row := r.database.QueryRowContext(ctx, qry, args...)

	err = row.Scan(&itemID)

	return itemID, errors.Wrap(err, "item create query error")
}

// ItemUpdate updates item by id in database.
func (r *Repository) ItemUpdate(ctxr context.Context, meta query.MetaData, input item.UpdateItemInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ItemUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(itemTable)

	if input.NameTm != nil {
		builder = builder.Set("name_tm", *input.NameTm)
	}

	if input.NameRu != nil {
		builder = builder.Set("name_ru", *input.NameRu)
	}

	if input.NameTr != nil {
		builder = builder.Set("name_tr", *input.NameTr)
	}

	if input.NameEn != nil {
		builder = builder.Set("name_en", *input.NameEn)
	}

	if input.DescriptionTm != nil {
		builder = builder.Set("description_tm", *input.DescriptionTm)
	}

	if input.DescriptionRu != nil {
		builder = builder.Set("description_ru", *input.DescriptionRu)
	}

	if input.DescriptionTr != nil {
		builder = builder.Set("description_tr", *input.DescriptionTr)
	}

	if input.DescriptionEn != nil {
		builder = builder.Set("description_en", *input.DescriptionEn)
	}

	if input.InternalID != nil {
		builder = builder.Set("internal_id", *input.InternalID)
	}

	if input.Price != nil {
		builder = builder.Set("price", *input.Price)
	}

	if input.CategoryID != nil {
		builder = builder.Set("category_id", *input.CategoryID)
	}

	if input.OrganizationID != nil {
		builder = builder.Set("organization_id", *input.OrganizationID)
	}

	if input.BrandID != nil {
		builder = builder.Set("brand_id", *input.BrandID)
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

	_, err = r.database.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "item update query error")
}

// ItemDelete deletes item by id from database.
func (r *Repository) ItemDelete(ctxr context.Context, meta query.MetaData, itemID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ItemDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(itemTable).
		Set("is_deleted", true).
		Set("user_deleted", meta.UserID).
		Set("deleted_at", time.Now().UTC()).
		Where(squirrel.Eq{"is_deleted": false, "id": itemID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.database.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "item delete query error")
}
