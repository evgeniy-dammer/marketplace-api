package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/rule"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// RuleGetAll selects all rules from database.
func (r *Repository) RuleGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]rule.Rule, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.RuleGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var rules []rule.Rule

	qry, args, err := r.ruleGetAllQuery(meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.SelectContext(ctx, &rules, qry, args...)

	return rules, errors.Wrap(err, "rules select query error")
}

// ruleGetAllQuery creates sql query.
func (r *Repository) ruleGetAllQuery(_ query.MetaData, params queryparameter.QueryParameter) (string, []interface{}, error) { //nolint:lll
	builder := r.genSQL.Select("id", "ptype", "v0", "v1", "v2", "v3", "v4", "v5").From(ruleTable)

	if params.Search != "" {
		search := "%" + params.Search + "%"

		builder = builder.Where(squirrel.Or{
			squirrel.Like{"ptype": search},
			squirrel.Like{"v0": search},
			squirrel.Like{"v1": search},
			squirrel.Like{"v2": search},
			squirrel.Like{"v3": search},
			squirrel.Like{"v4": search},
			squirrel.Like{"v5": search},
		})
	}

	if len(params.Sorts) > 0 {
		builder = builder.OrderBy(params.Sorts.Parsing(mappingSortRule)...)
	} else {
		builder = builder.OrderBy("id ASC")
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

// RuleGetOne select rule by id from database.
func (r *Repository) RuleGetOne(ctxr context.Context, _ query.MetaData, ruleID string) (rule.Rule, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.RuleGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var rle rule.Rule

	builder := r.genSQL.Select("id", "ptype", "v0", "v1", "v2", "v3", "v4", "v5").From(ruleTable).
		Where(squirrel.Eq{"id": ruleID})

	qry, args, err := builder.ToSql()
	if err != nil {
		return rle, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.GetContext(ctx, &rle, qry, args...)

	return rle, errors.Wrap(err, "rule select query error")
}

// RuleCreate insert rule into database.
func (r *Repository) RuleCreate(ctxr context.Context, _ query.MetaData, input rule.CreateRuleInput) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.RuleCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var ruleID string

	builder := r.genSQL.Insert(ruleTable).
		Columns("ptype", "v0", "v1", "v2", "v3", "v4", "v5").
		Values(input.Ptype, input.V0, input.V1, input.V2, input.V3, input.V4, input.V5).
		Suffix("RETURNING \"id\"")

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	row := r.database.QueryRowContext(ctx, qry, args...)
	err = row.Scan(&ruleID)

	return ruleID, errors.Wrap(err, "rule create query error")
}

// RuleUpdate updates rule by id in database.
func (r *Repository) RuleUpdate(ctxr context.Context, _ query.MetaData, input rule.UpdateRuleInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.RuleUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(tableTable)

	if input.Ptype != nil {
		builder = builder.Set("ptype", *input.Ptype)
	}

	if input.V0 != nil {
		builder = builder.Set("v0", *input.V0)
	}

	if input.V1 != nil {
		builder = builder.Set("v1", *input.V1)
	}

	if input.V2 != nil {
		builder = builder.Set("v2", *input.V2)
	}

	if input.V3 != nil {
		builder = builder.Set("v3", *input.V3)
	}

	if input.V4 != nil {
		builder = builder.Set("v4", *input.V4)
	}

	if input.V5 != nil {
		builder = builder.Set("v5", *input.V5)
	}

	builder = builder.Where(squirrel.Eq{"id": *input.ID})

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.database.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "rule update query error")
}

// RuleDelete deletes rule by id from database.
func (r *Repository) RuleDelete(ctxr context.Context, _ query.MetaData, ruleID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.RuleDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Delete(ruleTable).Where(squirrel.Eq{"id": ruleID})

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.database.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "rule delete query error")
}
