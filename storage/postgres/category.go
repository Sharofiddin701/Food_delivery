package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"food/api/models"
	"food/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CategoryRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewCategory(db *pgxpool.Pool, log logger.LoggerI) CategoryRepo {
	return CategoryRepo{
		db:  db,
		log: log,
	}
}

func (c *CategoryRepo) Create(ctx context.Context, category *models.Category) (*models.Category, error) {

	id := uuid.New()
	query := `INSERT INTO "category" (
		id,
		name,
		created_at,
		updated_at)
		VALUES($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err := c.db.Exec(context.Background(), query,
		id.String(),
		category.Name,
	)

	if err != nil {
		return &models.Category{}, err
	}
	return &models.Category{
		Id:        id.String(),
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}, nil
}

func (c *CategoryRepo) Update(ctx context.Context, category *models.Category) (*models.Category, error) {
	query := `UPDATE "category" SET 
		name=$1,
		updated_at=CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := c.db.Exec(context.Background(), query,
		category.Name,
		category.Id,
	)
	if err != nil {
		c.log.Error("error while updating category in strg" + err.Error())
		return &models.Category{}, err
	}
	return &models.Category{
		Id:        category.Id,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}, nil
}

func (c *CategoryRepo) GetAll(ctx context.Context, req *models.GetAllCategoriesRequest) (*models.GetAllCategoriesResponse, error) {
	var (
		resp   = &models.GetAllCategoriesResponse{}
		filter = ""
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
	    filter += fmt.Sprintf(` WHERE name ILIKE '%%%v%%'`, req.Search)
	}

	// Order by created_at DESC
	filter += " ORDER BY created_at DESC"

	// Append OFFSET and LIMIT to the filter
	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	query := `SELECT count(id) OVER(), id, name, created_at, updated_at FROM "category"` + filter
	rows, err := c.db.Query(context.Background(), query)
	if err != nil {
		return resp, err
	}
	for rows.Next() {
		var (
			category   = models.Category{}
			name       sql.NullString
			created_at sql.NullString
			updated_at sql.NullString
		)
		if err := rows.Scan(
			&resp.Count,
			&category.Id,
			&name,
			&created_at,
			&updated_at); err != nil {
			return resp, err
		}
		resp.Categories = append(resp.Categories, models.Category{
			Id:        category.Id,
			Name:      name.String,
			CreatedAt: created_at.String,
			UpdatedAt: updated_at.String,
		})
	}
	return resp, nil
}

func (c *CategoryRepo) GetByID(ctx context.Context, id string) (*models.Category, error) {
	var (
		category   = models.Category{}
		name       sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)
	if err := c.db.QueryRow(context.Background(), `SELECT id, name, created_at, updated_at FROM "category" WHERE id = $1`, id).Scan(
		&category.Id,
		&name,
		&created_at,
		&updated_at,
	); err != nil {
		return &models.Category{}, err
	}
	return &models.Category{
		Id:        category.Id,
		Name:      name.String,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}, nil
}

func (c *CategoryRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "category" WHERE id = $1`
	_, err := c.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
