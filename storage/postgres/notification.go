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

type NotificationRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewNotification(db *pgxpool.Pool, log logger.LoggerI) NotificationRepo {
	return NotificationRepo{
		db:  db,
		log: log,
	}
}

func (n *NotificationRepo) Create(ctx context.Context, notification *models.Notification) (*models.Notification, error) {

	id := uuid.New()
	query := `INSERT INTO "notifications" (
		id,
		user_id,
		message,
		status,
		created_at,
		updated_at)
		VALUES($1,$2,$3,$4,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) 
	`

	_, err := n.db.Exec(context.Background(), query,
		id.String(),
		notification.UserId,
		notification.Message,
		notification.IsRead,
	)

	if err != nil {
		return &models.Notification{}, err
	}
	return &models.Notification{
		Id:        id.String(),
		UserId:    notification.UserId,
		Message:   notification.Message,
		IsRead:    notification.IsRead,
		CreatedAt: notification.CreatedAt,
	}, nil
}

func (n *NotificationRepo) Update(ctx context.Context, notification *models.Notification) (*models.Notification, error) {
	query := `UPDATE "notifications" SET 
		user_id=$1,
		message=$2,
		status=$3,
		updated_at=CURRENT_TIMESTAMP
		WHERE id = $4
	`
	_, err := n.db.Exec(context.Background(), query,
		notification.UserId,
		notification.Message,
		notification.IsRead,
		notification.Id,
	)
	if err != nil {
		return &models.Notification{}, err
	}
	return &models.Notification{
		Id:        notification.Id,
		UserId:    notification.UserId,
		Message:   notification.Message,
		IsRead:    notification.IsRead,
		CreatedAt: notification.CreatedAt,
	}, nil
}

func (n *NotificationRepo) GetAll(ctx context.Context, req *models.GetAllNotificationsRequest) (*models.GetAllNotificationsResponse, error) {
	var (
		resp   = &models.GetAllNotificationsResponse{}
		filter = ""
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(` AND (message ILIKE '%%%v%%') `, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter: ", filter)

	rows, err := n.db.Query(context.Background(), `SELECT count(id) OVER(),
        id,
        user_id,
        message,
        status,
        created_at,
        updated_at FROM "notifications"`+filter)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var (
			notification models.Notification
			userId       sql.NullString
			message      sql.NullString
			is_read      sql.NullBool
			createdAt    sql.NullString
			updatedAt    sql.NullString
		)
		if err := rows.Scan(
			&resp.Count,
			&notification.Id,
			&userId,
			&message,
			&is_read,
			&createdAt,
			&updatedAt); err != nil {
			return resp, err
		}
		resp.Notifications = append(resp.Notifications, models.Notification{
			Id:        notification.Id,
			UserId:    userId.String,
			Message:   message.String,
			IsRead:    is_read.Bool,
			CreatedAt: createdAt.String,
		})
	}
	return resp, nil
}

func (n *NotificationRepo) GetByID(ctx context.Context, id string) (*models.Notification, error) {
	var (
		notification models.Notification
		userId       sql.NullString
		message      sql.NullString
		is_read      sql.NullBool
		createdAt    sql.NullString
		updatedAt    sql.NullString
	)
	if err := n.db.QueryRow(context.Background(), `SELECT id, user_id, message, status, created_at, updated_at FROM "notifications" WHERE id = $1`, id).Scan(
		&notification.Id,
		&userId,
		&message,
		&is_read,
		&createdAt,
		&updatedAt,
	); err != nil {
		return &models.Notification{}, err
	}
	return &models.Notification{
		Id:        notification.Id,
		UserId:    userId.String,
		Message:   message.String,
		IsRead:    is_read.Bool,
		CreatedAt: createdAt.String,
	}, nil
}

func (n *NotificationRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "notifications" WHERE id = $1`
	_, err := n.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
