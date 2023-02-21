package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// CommentPostgresql repository.
type CommentPostgresql struct {
	db *sqlx.DB
}

// NewCommentPostgresql is a constructor for CommentPostgresql.
func NewCommentPostgresql(db *sqlx.DB) *CommentPostgresql {
	return &CommentPostgresql{db: db}
}

// GetAll selects all comments from database.
func (r *CommentPostgresql) GetAll(userID string, organizationID string) ([]model.Comment, error) {
	var comments []model.Comment

	query := fmt.Sprintf(
		"SELECT id, item_id, organization_id, content, status_id, user_created, created_at FROM %s "+
			"WHERE is_deleted = false AND organization_id = $1 ",
		commentTable,
	)

	err := r.db.Select(&comments, query, organizationID)

	return comments, errors.Wrap(err, "comments select query error")
}

// GetOne select comment by id from database.
func (r *CommentPostgresql) GetOne(userID string, organizationID string, commentID string) (model.Comment, error) {
	var comment model.Comment

	query := fmt.Sprintf(
		"SELECT id, item_id, organization_id, content, status_id, user_created, created_at FROM %s "+
			"WHERE is_deleted = false AND organization_id = $1 AND id = $2 ",
		commentTable,
	)
	err := r.db.Get(&comment, query, organizationID, commentID)

	return comment, errors.Wrap(err, "comment select query error")
}

// Create insert comment into database.
func (r *CommentPostgresql) Create(userID string, comment model.Comment) (string, error) {
	var commentID string

	query := fmt.Sprintf(
		"INSERT INTO %s (item_id, organization_id, content, status_id, user_created) "+
			"VALUES ($1, $2, $3, $4, $5) RETURNING id",
		commentTable,
	)

	row := r.db.QueryRow(
		query,
		comment.ItemID,
		comment.OrganizationID,
		comment.Content,
		comment.Status,
		userID,
	)

	err := row.Scan(&commentID)

	return commentID, errors.Wrap(err, "comment create query error")
}

// Update updates comment by id in database.
func (r *CommentPostgresql) Update(userID string, input model.UpdateCommentInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.ItemID != nil {
		setValues = append(setValues, fmt.Sprintf("item_id=$%d", argID))
		args = append(args, *input.ItemID)
		argID++
	}

	if input.Content != nil {
		setValues = append(setValues, fmt.Sprintf("content=$%d", argID))
		args = append(args, *input.Content)
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

	setValues = append(setValues, fmt.Sprintf("user_updated=$%d", argID))
	args = append(args, userID)
	argID++

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE is_deleted = false AND organization_id = '%s' AND id = '%s'",
		commentTable, setQuery, *input.OrganizationID, *input.ID)

	_, err := r.db.Exec(query, args...)

	return errors.Wrap(err, "comment update query error")
}

// Delete deletes comment by id from database.
func (r *CommentPostgresql) Delete(userID string, organizationID string, commentID string) error {
	query := fmt.Sprintf(
		"UPDATE %s SET is_deleted = true, deleted_at = $1, user_deleted = $2 "+
			"WHERE is_deleted = false AND id = $3 AND organization_id = $4",
		commentTable,
	)

	_, err := r.db.Exec(query, time.Now().Format("2006-01-02 15:04:05"), userID, commentID, organizationID)

	return errors.Wrap(err, "comment delete query error")
}
