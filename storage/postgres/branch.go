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

type BranchRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewBranchRepo(db *pgxpool.Pool, log logger.LoggerI) BranchRepo {
	return BranchRepo{
		db:  db,
		log: log,
	}
}

func (b *BranchRepo) Create(ctx context.Context, branch *models.Branch) (*models.Branch, error) {
	id := uuid.New()
	query := `INSERT INTO "branch" (
		id,
		name,
		address,
		latitude,
		longitude,
		created_at,
		updated_at)
		VALUES($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err := b.db.Exec(context.Background(), query,
		id.String(),
		branch.Name,
		branch.Address,
		branch.Latitude,
		branch.Longitude,
	)

	if err != nil {
		return &models.Branch{}, err
	}
	return &models.Branch{
		Id:        id.String(),
		Name:      branch.Name,
		Address:   branch.Address,
		Latitude:  branch.Latitude,
		Longitude: branch.Longitude,
		CreatedAt: branch.CreatedAt,
		UpdatedAt: branch.UpdatedAt,
	}, nil
}

func (b *BranchRepo) Update(ctx context.Context, branch *models.Branch) (*models.Branch, error) {
	query := `UPDATE "branch" SET 
		name=$1,
		address=$2,
		latitude=$3,
		longitude=$4,
		updated_at=CURRENT_TIMESTAMP
		WHERE id = $5
	`
	_, err := b.db.Exec(context.Background(), query,
		branch.Name,
		branch.Address,
		branch.Latitude,
		branch.Longitude,
		branch.Id,
	)
	if err != nil {
		return &models.Branch{}, err
	}
	return &models.Branch{
		Id:        branch.Id,
		Name:      branch.Name,
		Address:   branch.Address,
		Latitude:  branch.Latitude,
		Longitude: branch.Longitude,
		CreatedAt: branch.CreatedAt,
		UpdatedAt: branch.UpdatedAt,
	}, nil
}

func (b *BranchRepo) GetAll(ctx context.Context, req *models.GetAllBranchesRequest) (*models.GetAllBranchesResponse, error) {
	var (
		resp   = &models.GetAllBranchesResponse{}
		filter = ""
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(` AND (name ILIKE '%%%v%%' OR address ILIKE '%%%v%%') `, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter: ", filter)

	rows, err := b.db.Query(context.Background(), `SELECT count(id) OVER(),
        id,
        name,
        address,
        latitude,
        longitude,
        created_at,
        updated_at FROM "branch"`+filter)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var (
			branch     = models.Branch{}
			name       sql.NullString
			address    sql.NullString
			latitude   sql.NullFloat64
			longitude  sql.NullFloat64
			created_at sql.NullString
			updated_at sql.NullString
		)
		if err := rows.Scan(
			&resp.Count,
			&branch.Id,
			&name,
			&address,
			&latitude,
			&longitude,
			&created_at,
			&updated_at); err != nil {
			return resp, err
		}

		resp.Branches = append(resp.Branches, models.Branch{
			Id:        branch.Id,
			Name:      name.String,
			Address:   address.String,
			Latitude:  latitude.Float64,
			Longitude: longitude.Float64,
			CreatedAt: created_at.String,
			UpdatedAt: updated_at.String,
		})
	}
	return resp, nil
}

func (b *BranchRepo) GetByID(ctx context.Context, id string) (*models.Branch, error) {
	var (
		branch     = models.Branch{}
		name       sql.NullString
		address    sql.NullString
		latitude   sql.NullFloat64
		longitude  sql.NullFloat64
		created_at sql.NullString
		updated_at sql.NullString
	)
	if err := b.db.QueryRow(context.Background(), `SELECT id, name, address, latitude, longitude, created_at, updated_at FROM "branch" WHERE id = $1`, id).Scan(
		&branch.Id,
		&name,
		&address,
		&latitude,
		&longitude,
		&created_at,
		&updated_at,
	); err != nil {
		return &models.Branch{}, err
	}
	return &models.Branch{
		Id:        branch.Id,
		Name:      name.String,
		Address:   address.String,
		Latitude:  latitude.Float64,
		Longitude: longitude.Float64,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}, nil
}

func (b *BranchRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "branch" WHERE id = $1`
	_, err := b.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
