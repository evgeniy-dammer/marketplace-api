package redis

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/rule"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
)

// RuleGetAll gets rules from cache.
func (r *Repository) RuleGetAll(ctxr context.Context, _ query.MetaData, _ queryparameter.QueryParameter) ([]rule.Rule, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.RuleGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	rules := &rule.ListRule{}

	bytes, err := r.client.Get(ctx, rulesKey).Bytes()
	if err != nil {
		return *rules, errors.Wrap(err, "unable to get rules from cache")
	}

	if err = easyjson.Unmarshal(bytes, rules); err != nil {
		return *rules, errors.Wrap(err, "unable to unmarshal")
	}

	return *rules, nil
}

// RuleSetAll sets rules into cache.
func (r *Repository) RuleSetAll(ctxr context.Context, _ query.MetaData, _ queryparameter.QueryParameter, rules []rule.Rule) error { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.RuleSetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	ruleSlice := rule.ListRule(rules)

	bytes, err := easyjson.Marshal(ruleSlice)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, rulesKey, bytes, r.options.TTL).Err()

	return errors.Wrap(err, "setting key")
}

// RuleGetOne gets rule by id from cache.
func (r *Repository) RuleGetOne(ctxr context.Context, ruleID string) (rule.Rule, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.RuleGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var rle rule.Rule

	bytes, err := r.client.Get(ctx, ruleKey+ruleID).Bytes()
	if err != nil {
		return rle, errors.Wrap(err, "unable to get rule from cache")
	}

	if err = easyjson.Unmarshal(bytes, &rle); err != nil {
		return rle, errors.Wrap(err, "unable to unmarshal")
	}

	return rle, nil
}

// RuleCreate sets rule into cache.
func (r *Repository) RuleCreate(ctxr context.Context, usr rule.Rule) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.RuleCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, ruleKey+usr.ID, bytes, r.options.TTL).Err()

	return errors.Wrap(err, "setting key")
}

// RuleUpdate updates rule by id in cache.
func (r *Repository) RuleUpdate(ctxr context.Context, usr rule.Rule) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.RuleUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	r.client.Set(ctx, ruleKey+usr.ID, bytes, r.options.TTL)

	return nil
}

// RuleDelete deletes rule by id from cache.
func (r *Repository) RuleDelete(ctxr context.Context, ruleID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.RuleDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := r.client.Del(ctx, ruleKey+ruleID).Err()

	return errors.Wrap(err, "deleting key")
}

// RuleInvalidate invalidate rules cache.
func (r *Repository) RuleInvalidate(ctxr context.Context) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.RuleInvalidate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	iter := r.client.Scan(ctx, 0, rulesKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			return errors.Wrap(err, "deleting key")
		}
	}

	return errors.Wrap(iter.Err(), "invalidate")
}
