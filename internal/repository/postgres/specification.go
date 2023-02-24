package postgres

import (
	"fmt"
	"strings"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// SpecificationPostgresql repository.
type SpecificationPostgresql struct {
	db *sqlx.DB
}

// NewSpecificationPostgresql is a constructor for SpecificationPostgresql.
func NewSpecificationPostgresql(db *sqlx.DB) *SpecificationPostgresql {
	return &SpecificationPostgresql{db: db}
}

// GetAll selects all specifications from database.
func (r *SpecificationPostgresql) GetAll(userID string, organizationID string) ([]model.Specification, error) {
	var specifications []model.Specification

	query := fmt.Sprintf(
		"SELECT id, item_id, organization_id, name_tm, name_ru, name_tr, name_en, description_tm, description_ru, "+
			"description_tr, description_en, value FROM %s WHERE organization_id = $1 ",
		specificationTable,
	)

	err := r.db.Select(&specifications, query, organizationID)

	return specifications, errors.Wrap(err, "specifications select query error")
}

// GetOne select specification by id from database.
func (r *SpecificationPostgresql) GetOne(userID string, organizationID string, specificationID string) (model.Specification, error) {
	var specification model.Specification

	query := fmt.Sprintf(
		"SELECT id, item_id, organization_id, name_tm, name_ru, name_tr, name_en, description_tm, description_ru, "+
			"description_tr, description_en, value FROM %s WHERE organization_id = $1 AND id = $2 ",
		specificationTable,
	)

	err := r.db.Get(&specification, query, organizationID, specificationID)

	return specification, errors.Wrap(err, "specification select query error")
}

// Create insert specification into database.
func (r *SpecificationPostgresql) Create(userID string, specification model.Specification) (string, error) {
	var specificationID string

	query := fmt.Sprintf(
		"INSERT INTO %s (item_id, organization_id, name_tm, name_ru, name_tr, name_en,  description_tm, description_ru, "+
			"description_tr, description_en, value) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		specificationTable,
	)

	row := r.db.QueryRow(
		query,
		specification.ItemID,
		specification.OrganizationID,
		specification.NameTm,
		specification.NameRu,
		specification.NameTr,
		specification.NameEn,
		specification.DescriptionTm,
		specification.DescriptionRu,
		specification.DescriptionTr,
		specification.DescriptionEn,
		specification.Value,
	)

	err := row.Scan(&specificationID)

	return specificationID, errors.Wrap(err, "specification create query error")
}

// Update updates specification by id in database.
func (r *SpecificationPostgresql) Update(userID string, input model.UpdateSpecificationInput) error {
	setValues := make([]string, 0, 12)
	args := make([]interface{}, 0, 12)
	argID := 1

	if input.ItemID != nil {
		setValues = append(setValues, fmt.Sprintf("item_id=$%d", argID))
		args = append(args, *input.ItemID)
		argID++
	}

	if input.OrganizationID != nil {
		setValues = append(setValues, fmt.Sprintf("organization_id=$%d", argID))
		args = append(args, *input.OrganizationID)
		argID++
	}

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

	if input.Value != nil {
		setValues = append(setValues, fmt.Sprintf("value=$%d", argID))
		args = append(args, *input.Value)
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE organization_id = '%s' AND id = '%s'",
		specificationTable, setQuery, *input.OrganizationID, *input.ID)

	_, err := r.db.Exec(query, args...)

	return errors.Wrap(err, "specification update query error")
}

// Delete deletes specification by id from database.
func (r *SpecificationPostgresql) Delete(userID string, organizationID string, specificationID string) error {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id = $1 AND organization_id = $2",
		specificationTable,
	)

	_, err := r.db.Exec(query, specificationID, organizationID)

	return errors.Wrap(err, "specification delete query error")
}
