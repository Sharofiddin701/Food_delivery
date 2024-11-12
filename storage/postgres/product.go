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

type ProductRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewProduct(db *pgxpool.Pool, log logger.LoggerI) ProductRepo {
	return ProductRepo{
		db:  db,
		log: log,
	}
}

func (p *ProductRepo) Create(ctx context.Context, product *models.Product) (*models.Product, error) {

	id := uuid.New()
	query := `INSERT INTO "product" (
		id,
		category_id,
		description,
		price,
		image_url,
		name,
		created_at,
		updated_at)
		VALUES($1,$2,$3,$4,$5,$6, CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) 
	`

	_, err := p.db.Exec(context.Background(), query,
		id.String(),
		product.CategoryId,
		product.Description,
		product.Price,
		product.ImageURL,
		product.Name,
	)

	if err != nil {
		return &models.Product{}, err
	}
	return &models.Product{
		Id:          id.String(),
		Name:        product.Name,
		CategoryId:  product.CategoryId,
		Description: product.Description,
		Price:       product.Price,
		ImageURL:    product.ImageURL,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

func (p *ProductRepo) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	query := `UPDATE "product" SET 
		name=$1,
		category_id=$2,
		description=$3,
		price=$4,
		image_url=$5,
		updated_at=CURRENT_TIMESTAMP
		WHERE id = $6
		`
	_, err := p.db.Exec(context.Background(), query,
		product.Name,
		product.CategoryId,
		product.Description,
		product.Price,
		product.ImageURL,
		product.Id,
	)
	if err != nil {
		return &models.Product{}, err
	}
	return &models.Product{
		Id:          product.Id,
		Name:        product.Name,
		CategoryId:  product.CategoryId,
		Description: product.Description,
		Price:       product.Price,
		ImageURL:    product.ImageURL,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

func (p *ProductRepo) GetAll(ctx context.Context, req *models.GetAllProductsRequest) (*models.GetAllProductsResponse, error) {
	var (
		resp   = &models.GetAllProductsResponse{}
		filter = ""
		args   []interface{}
		argIdx = 1
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(" WHERE (name ILIKE $%d) ", argIdx)
		args = append(args, "%"+req.Search+"%")
		argIdx++
	}

	if req.CategoryId != "" {
		if filter == "" {
			filter += fmt.Sprintf(" WHERE category_id = $%d::uuid ", argIdx)
		} else {
			filter += fmt.Sprintf(" AND category_id = $%d::uuid ", argIdx)
		}
		args = append(args, req.CategoryId)
		argIdx++
	}

	filter += fmt.Sprintf(" OFFSET %d LIMIT %d", offset, req.Limit)
	fmt.Println("filter: ", filter)

	query := `SELECT count(id) OVER(),
		id,
		category_id,
		name,
		description,
		price,
		image_url,
		created_at,
		updated_at 
		FROM "product"` + filter

	rows, err := p.db.Query(context.Background(), query, args...)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var (
			product     = models.Product{}
			category_id sql.NullString
			name        sql.NullString
			description sql.NullString
			price       sql.NullFloat64
			image_url   sql.NullString
			created_at  sql.NullString
			updated_at  sql.NullString
		)
		if err := rows.Scan(
			&resp.Count,
			&product.Id,
			&category_id,
			&name,
			&description,
			&price,
			&image_url,
			&created_at,
			&updated_at); err != nil {
			return resp, err
		}

		resp.Products = append(resp.Products, models.Product{
			Id:          product.Id,
			CategoryId:  category_id.String,
			Name:        name.String,
			Description: description.String,
			Price:       price.Float64,
			ImageURL:    image_url.String,
			CreatedAt:   created_at.String,
			UpdatedAt:   updated_at.String,
		})
	}
	return resp, nil
}

func (p *ProductRepo) GetByID(ctx context.Context, id string) (*models.Product, error) {
	var (
		product     = models.Product{}
		category_id sql.NullString
		name        sql.NullString
		description sql.NullString
		price       sql.NullFloat64
		image_url   sql.NullString
		created_at  sql.NullString
		updated_at  sql.NullString
	)
	if err := p.db.QueryRow(context.Background(), `SELECT id, category_id, name, description, price, image_url, created_at, updated_at FROM "product" WHERE id = $1`, id).Scan(
		&product.Id,
		&name,
		&category_id,
		&description,
		&price,
		&image_url,
		&created_at,
		&updated_at,
	); err != nil {
		return &models.Product{}, err
	}
	return &models.Product{
		Id:          product.Id,
		CategoryId:  category_id.String,
		Name:        name.String,
		Description: description.String,
		Price:       price.Float64,
		ImageURL:    image_url.String,
		CreatedAt:   created_at.String,
		UpdatedAt:   updated_at.String,
	}, nil
}

func (p *ProductRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "product" WHERE id = $1`
	_, err := p.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
