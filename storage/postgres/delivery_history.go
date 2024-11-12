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

type DeliveryHistoryRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewDeliveryHistory(db *pgxpool.Pool, log logger.LoggerI) DeliveryHistoryRepo {
	return DeliveryHistoryRepo{
		db: db,
		log: log,
	}
}

func (r *DeliveryHistoryRepo) Create(ctx context.Context, deliveryHistory *models.DeliveryHistory) (*models.DeliveryHistory, error) {

	id := uuid.New()
	query := `INSERT INTO "delivery_history" (
		id,
		courier_id,
		order_id,
		earnings,
		delivered_at)
		VALUES($1,$2,$3,$4,$5) 
	`

	_, err := r.db.Exec(context.Background(), query,
		id.String(),
		deliveryHistory.CourierId,
		deliveryHistory.OrderId,
		deliveryHistory.Earnings,
		deliveryHistory.DeliveredAt,
	)

	if err != nil {
		return &models.DeliveryHistory{}, err
	}
	return &models.DeliveryHistory{
		Id:          id.String(),
		CourierId:   deliveryHistory.CourierId,
		OrderId:     deliveryHistory.OrderId,
		Earnings:    deliveryHistory.Earnings,
		DeliveredAt: deliveryHistory.DeliveredAt,
	}, nil
}

func (r *DeliveryHistoryRepo) Update(ctx context.Context, deliveryHistory *models.DeliveryHistory) (*models.DeliveryHistory, error) {
	query := `UPDATE "delivery_history" SET 
		courier_id=$1,
		order_id=$2,
		earnings=$3,
		delivered_at=$4
		WHERE id = $5
	`
	_, err := r.db.Exec(context.Background(), query,
		deliveryHistory.CourierId,
		deliveryHistory.OrderId,
		deliveryHistory.Earnings,
		deliveryHistory.DeliveredAt,
		deliveryHistory.Id,
	)
	if err != nil {
		return &models.DeliveryHistory{}, err
	}
	return &models.DeliveryHistory{
		Id:          deliveryHistory.Id,
		CourierId:   deliveryHistory.CourierId,
		OrderId:     deliveryHistory.OrderId,
		Earnings:    deliveryHistory.Earnings,
		DeliveredAt: deliveryHistory.DeliveredAt,
	}, nil
}

func (r *DeliveryHistoryRepo) GetAll(ctx context.Context, req *models.GetAllDeliveryHistoriesRequest) (*models.GetAllDeliveryHistoriesResponse, error) {
	var (
		resp   = &models.GetAllDeliveryHistoriesResponse{}
		filter = ""
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(` AND (courier_id ILIKE '%%%v%%' OR order_id ILIKE '%%%v%%') `, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter: ", filter)

	rows, err := r.db.Query(context.Background(), `SELECT count(id) OVER(),
        id,
        courier_id,
        order_id,
        earnings,
        delivered_at FROM "delivery_history"`+filter)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var (
			deliveryHistory = models.DeliveryHistory{}
			courierId       sql.NullString
			orderId         sql.NullString
			earnings        sql.NullFloat64
			deliveredAt     sql.NullString
		)
		if err := rows.Scan(
			&resp.Count,
			&deliveryHistory.Id,
			&courierId,
			&orderId,
			&earnings,
			&deliveredAt); err != nil {
			return resp, err
		}
		resp.DeliveryHistories = append(resp.DeliveryHistories, models.DeliveryHistory{
			Id:          deliveryHistory.Id,
			CourierId:   courierId.String,
			OrderId:     orderId.String,
			Earnings:    earnings.Float64,
			DeliveredAt: deliveredAt.String,
		})
	}
	return resp, nil
}

func (r *DeliveryHistoryRepo) GetByID(ctx context.Context, id string) (*models.DeliveryHistory, error) {
	var (
		deliveryHistory = models.DeliveryHistory{}
		courierId       sql.NullString
		orderId         sql.NullString
		earnings        sql.NullFloat64
		deliveredAt     sql.NullString
	)
	if err := r.db.QueryRow(context.Background(), `SELECT id, courier_id, order_id, earnings, delivered_at FROM "delivery_history" WHERE id = $1`, id).Scan(
		&deliveryHistory.Id,
		&courierId,
		&orderId,
		&earnings,
		&deliveredAt,
	); err != nil {
		return &models.DeliveryHistory{}, err
	}
	return &models.DeliveryHistory{
		Id:          deliveryHistory.Id,
		CourierId:   courierId.String,
		OrderId:     orderId.String,
		Earnings:    earnings.Float64,
		DeliveredAt: deliveredAt.String,
	}, nil
}

func (r *DeliveryHistoryRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "delivery_history" WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
