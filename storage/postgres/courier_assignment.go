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

type CourierAssignmentRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewCourierAssignment(db *pgxpool.Pool, log logger.LoggerI) CourierAssignmentRepo {
	return CourierAssignmentRepo{
		db:  db,
		log: log,
	}
}

func (c *CourierAssignmentRepo) Create(ctx context.Context, assignment *models.CourierAssignment) (*models.CourierAssignment, error) {

	id := uuid.New()
	query := `INSERT INTO "courier_assignments" (
		id,
		order_id,
		courier_id,
		assigned_at,
		status,
		updated_at)
		VALUES($1,$2,$3,CURRENT_TIMESTAMP,$4,CURRENT_TIMESTAMP) 
	`

	_, err := c.db.Exec(context.Background(), query,
		id.String(),
		assignment.OrderId,
		assignment.CourierId,
		assignment.Status,
	)

	if err != nil {
		return &models.CourierAssignment{}, err
	}
	return &models.CourierAssignment{
		Id:         id.String(),
		OrderId:    assignment.OrderId,
		CourierId:  assignment.CourierId,
		AssignedAt: assignment.AssignedAt,
		Status:     assignment.Status,
		UpdatedAt:  assignment.UpdatedAt,
	}, nil
}

func (c *CourierAssignmentRepo) Update(ctx context.Context, assignment *models.CourierAssignment) (*models.CourierAssignment, error) {
	query := `UPDATE "courier_assignments" SET 
		order_id=$1,
		courier_id=$2,
		status=$3,
		updated_at=CURRENT_TIMESTAMP
		WHERE id = $4
	`
	_, err := c.db.Exec(context.Background(), query,
		assignment.OrderId,
		assignment.CourierId,
		assignment.Status,
		assignment.Id,
	)
	if err != nil {
		return &models.CourierAssignment{}, err
	}
	return &models.CourierAssignment{
		Id:         assignment.Id,
		OrderId:    assignment.OrderId,
		CourierId:  assignment.CourierId,
		Status:     assignment.Status,
		AssignedAt: assignment.AssignedAt,
		UpdatedAt:  assignment.UpdatedAt,
	}, nil
}

func (c *CourierAssignmentRepo) GetAll(ctx context.Context, req *models.GetAllCourierAssignmentsRequest) (*models.GetAllCourierAssignmentsResponse, error) {
	var (
		resp   = &models.GetAllCourierAssignmentsResponse{}
		filter = ""
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(` AND (order_id ILIKE '%%%v%%' OR courier_id ILIKE '%%%v%%' OR status ILIKE '%%%v%%') `, req.Search, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter: ", filter)

	rows, err := c.db.Query(context.Background(), `SELECT count(id) OVER(),
        id,
        order_id,
        courier_id,
        assigned_at,
        status,
        updated_at FROM courier_assignments FROM "courier_assignment"`+filter)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var (
			assignment models.CourierAssignment
			orderID    sql.NullString
			courierID  sql.NullString
			assignedAt sql.NullString
			status     sql.NullString
			updatedAt  sql.NullString
		)
		if err := rows.Scan(
			&resp.Count,
			&assignment.Id,
			&orderID,
			&courierID,
			&assignedAt,
			&status,
			&updatedAt); err != nil {
			return resp, err
		}

		resp.CourierAssignments = append(resp.CourierAssignments, models.CourierAssignment{
			Id:         assignment.Id,
			OrderId:    orderID.String,
			CourierId:  courierID.String,
			AssignedAt: assignedAt.String,
			Status:     status.String,
			UpdatedAt:  updatedAt.String,
		})
	}
	return resp, nil
}

func (c *CourierAssignmentRepo) GetByID(ctx context.Context, id string) (*models.CourierAssignment, error) {
	var (
		assignment models.CourierAssignment
		orderID    sql.NullString
		courierID  sql.NullString
		assignedAt sql.NullString
		status     sql.NullString
		createdAt  sql.NullString
		updatedAt  sql.NullString
	)
	if err := c.db.QueryRow(context.Background(), `SELECT id, order_id, courier_id, assigned_at, status, created_at, updated_at FROM "courier_assignments" WHERE id = $1`, id).Scan(
		&assignment.Id,
		&orderID,
		&courierID,
		&assignedAt,
		&status,
		&createdAt,
		&updatedAt,
	); err != nil {
		return &models.CourierAssignment{}, err
	}
	return &models.CourierAssignment{
		Id:         assignment.Id,
		OrderId:    orderID.String,
		CourierId:  courierID.String,
		AssignedAt: assignedAt.String,
		Status:     status.String,
		UpdatedAt:  updatedAt.String,
	}, nil
}

func (c *CourierAssignmentRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM courier_assignments WHERE id = $1`
	_, err := c.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
