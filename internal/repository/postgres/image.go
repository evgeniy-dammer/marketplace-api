package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ImagePostgresql repository.
type ImagePostgresql struct {
	db *sqlx.DB
}

// NewImagePostgresql is a constructor for ImagePostgresql.
func NewImagePostgresql(db *sqlx.DB) *ImagePostgresql {
	return &ImagePostgresql{db: db}
}

// GetAll selects all images from database.
func (r *ImagePostgresql) GetAll(userID string, organizationID string) ([]model.Image, error) {
	var images []model.Image

	query := fmt.Sprintf(
		"SELECT id, object_id, type, origin, middle, small, organization_id, is_main FROM %s "+
			"WHERE is_deleted = false AND organization_id = $1 ",
		imageTable,
	)

	err := r.db.Select(&images, query, organizationID)

	return images, errors.Wrap(err, "images select query error")
}

// GetOne select image by id from database.
func (r *ImagePostgresql) GetOne(userID string, organizationID string, imageID string) (model.Image, error) {
	var image model.Image

	query := fmt.Sprintf(
		"SELECT id, object_id, type, origin, middle, small, organization_id, is_main FROM %s "+
			"WHERE is_deleted = false AND organization_id = $1 AND id = $2 ",
		imageTable,
	)
	err := r.db.Get(&image, query, organizationID, imageID)

	return image, errors.Wrap(err, "image select query error")
}

// Create insert image into database.
func (r *ImagePostgresql) Create(userID string, image model.Image) (string, error) {
	var imageID string

	query := fmt.Sprintf(
		"INSERT INTO %s (object_id, type, origin, middle, small, organization_id, is_main, user_created) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		imageTable,
	)

	row := r.db.QueryRow(
		query,
		image.ObjectID,
		image.ObjectType,
		image.Origin,
		image.Middle,
		image.Small,
		image.OrganizationID,
		image.IsMain,
		userID,
	)

	err := row.Scan(&imageID)

	return imageID, errors.Wrap(err, "image create query error")
}

// Update updates image by id in database.
func (r *ImagePostgresql) Update(userID string, input model.UpdateImageInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.ObjectID != nil {
		setValues = append(setValues, fmt.Sprintf("object_id=$%d", argID))
		args = append(args, *input.ObjectID)
		argID++
	}

	if input.ObjectType != nil {
		setValues = append(setValues, fmt.Sprintf("type=$%d", argID))
		args = append(args, *input.ObjectType)
		argID++
	}

	if input.Origin != nil {
		setValues = append(setValues, fmt.Sprintf("origin=$%d", argID))
		args = append(args, *input.Origin)
		argID++
	}

	if input.Middle != nil {
		setValues = append(setValues, fmt.Sprintf("middle=$%d", argID))
		args = append(args, *input.Middle)
		argID++
	}

	if input.Small != nil {
		setValues = append(setValues, fmt.Sprintf("small=$%d", argID))
		args = append(args, *input.Small)
		argID++
	}

	if input.OrganizationID != nil {
		setValues = append(setValues, fmt.Sprintf("organization_id=$%d", argID))
		args = append(args, *input.OrganizationID)
		argID++
	}

	if input.IsMain != nil {
		setValues = append(setValues, fmt.Sprintf("is_main=$%d", argID))
		args = append(args, *input.IsMain)
		argID++
	}

	setValues = append(setValues, fmt.Sprintf("user_updated=$%d", argID))
	args = append(args, userID)
	argID++

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE is_deleted = false AND organization_id = '%s' AND id = '%s'",
		imageTable, setQuery, *input.OrganizationID, *input.ID)

	_, err := r.db.Exec(query, args...)

	return errors.Wrap(err, "image update query error")
}

// Delete deletes image by id from database.
func (r *ImagePostgresql) Delete(userID string, organizationID string, imageID string) error {
	query := fmt.Sprintf(
		"UPDATE %s SET is_deleted = true, deleted_at = $1, user_deleted = $2 "+
			"WHERE is_deleted = false AND id = $3 AND organization_id = $4",
		imageTable,
	)

	_, err := r.db.Exec(query, time.Now().Format("2006-01-02 15:04:05"), userID, imageID, organizationID)

	return errors.Wrap(err, "image delete query error")
}
