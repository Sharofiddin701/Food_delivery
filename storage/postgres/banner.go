package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"food/api/models"
	"food/pkg/logger"

	"github.com/jackc/pgx/v4/pgxpool"
)

type BannerRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewBannerRepo(db *pgxpool.Pool, log logger.LoggerI) *BannerRepo {
	return &BannerRepo{
		db:  db,
		log: log,
	}
}

func (b *BannerRepo) Create(ctx context.Context, banner *models.Banner) (*models.Banner, error) {

	query := `INSERT INTO "banner" (
		image_url,
		created_at)
		VALUES ($1, CURRENT_TIMESTAMP)
	`

	_, err := b.db.Exec(ctx, query,
		banner.ImageUrl,
	)
	if err != nil {
		b.log.Error("Error creating banner: " + err.Error())
		return nil, err
	}

	b.log.Info("Banner created successfully!")
	return &models.Banner{
		ImageUrl:   banner.ImageUrl,
		Created_at: banner.Created_at,
	}, nil
}

func (b *BannerRepo) GetAll(ctx context.Context, req *models.GetAllBannerRequest) (*models.GetAllBannerResponse, error) {
	var (
		resp   = &models.GetAllBannerResponse{}
		filter = ""
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(` WHERE image_url ILIKE '%%%v%%' `, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)

	query := fmt.Sprintf(`SELECT count(image_url) OVER(), image_url, created_at FROM "banner" %s`, filter)
	rows, err := b.db.Query(ctx, query)
	if err != nil {
		b.log.Error("Error retrieving banners: " + err.Error())
		return nil, err
	}

	for rows.Next() {
		var (
			imageUrl   sql.NullString
			created_at sql.NullString
		)

		if err := rows.Scan(
			&resp.Count,
			&imageUrl,
			&created_at); err != nil {
			return nil, err
		}

		resp.Banners = append(resp.Banners, models.Banner{
			ImageUrl:   imageUrl.String,
			Created_at: created_at.String,
		})
	}

	b.log.Info("Banners retrieved successfully!")
	return resp, nil
}

func (b *BannerRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "banner" WHERE image_url = $1`
	_, err := b.db.Exec(ctx, query, id)
	if err != nil {
		b.log.Error("Error deleting banner: " + err.Error())
		return err
	}

	b.log.Info("Banner deleted successfully!")
	return nil
}
