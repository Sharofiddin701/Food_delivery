package postgres

import (
	"context"
	"food/api/models"
	"food/pkg/logger"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PaymentRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewPayment(db *pgxpool.Pool, log logger.LoggerI) PaymentRepo {
	return PaymentRepo{
		db:  db,
		log: log,
	}
}

// Create implements storage.IPaymentStorage.
func (p *PaymentRepo) Create(ctx context.Context, payment *models.Payment) (*models.Payment, error) {

	query := `INSERT INTO "payment" (
		id,
		user_id,
        order_id,
        is_paid,
        payment_method,
        created_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
	`
	_, err := p.db.Exec(ctx, query,
		payment.Id,
		payment.UserId,
		payment.OrderId,
		payment.IsPaid,
		payment.PaymentMethod, // Pass payment method to the query
	)
	if err != nil {
		p.log.Error("Error creating payment: " + err.Error())
		return nil, err
	}

	// Return the newly created payment model
	return &models.Payment{
		Id:            payment.Id,
		UserId:        payment.UserId,
		OrderId:       payment.OrderId,
		IsPaid:        payment.IsPaid,
		PaymentMethod: payment.PaymentMethod,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}, nil
}

/*
CREATE TABLE IF NOT EXISTS "payment" (
  id UUID PRIMARY KEY DEFAULT,
  user_id UUID NOT NULL REFERENCES "user"(id),
  order_id UUID NOT NULL REFERENCES "order"(id),
  is_paid BOOLEAN NOT NULL DEFAULT false,
  payment_method VARCHAR NOT NULL CHECK (payment_method IN ('click', 'payme', 'naxt pul')),
  created_at TIMESTAMP DEFAULT now()
);

*/

// GetAll implements storage.IPaymentStorage.
func (p *PaymentRepo) GetAll(ctx context.Context, request *models.GetAllPaymentsRequest) (*models.GetAllPaymentsResponse, error) {
	panic("unimplemented")
}

// GetByID implements storage.IPaymentStorage.
func (p *PaymentRepo) GetByID(ctx context.Context, id string) (*models.Payment, error) {
	
	query := `SELECT id, user_id, order_id, is_paid, payment_method, created_at 
	          FROM "payment" WHERE id = $1`

	var payment models.Payment

	err := p.db.QueryRow(ctx, query, id).Scan(
		&payment.Id,
		&payment.UserId,
		&payment.OrderId,
		&payment.IsPaid,
		&payment.PaymentMethod,
		&payment.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			p.log.Warn("Payment not found with ID: " + id)
			return nil, nil
		}
		p.log.Error("Error retrieving payment by ID: " + err.Error())
		return nil, err
	}

	return &payment, nil
}

// Update implements storage.IPaymentStorage.
func (p *PaymentRepo) Update(context.Context, *models.Payment) (*models.Payment, error) {
	panic("unimplemented")
}

// Delete implements storage.IPaymentStorage.
func (p *PaymentRepo) Delete(context.Context, string) error {
	panic("unimplemented")
}
