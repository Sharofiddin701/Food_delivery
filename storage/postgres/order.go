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

type OrderRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewOrder(db *pgxpool.Pool, log logger.LoggerI) OrderRepo {
	return OrderRepo{
		db:  db,
		log: log,
	}
}

func (o *OrderRepo) Create(ctx context.Context, order *models.OrderCreateRequest) (*models.OrderCreateRequest, error) {
	tx, err := o.db.Begin(context.Background())
	if err != nil {
		return &models.OrderCreateRequest{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	// Generate a new UUID for the order
	orderId := uuid.New().String()

	var totalSum float64
	for i, item := range order.Items {
		if item.Quantity <= 0 {
			return &models.OrderCreateRequest{}, fmt.Errorf("quantity must be greater than 0 for product %s", item.ProductId)
		}

		var productPrice float64
		productQuery := `SELECT price FROM "product" WHERE id = $1`
		err = o.db.QueryRow(context.Background(), productQuery, item.ProductId).Scan(&productPrice)
		if err != nil {
			return &models.OrderCreateRequest{}, fmt.Errorf("failed to retrieve price for product %s: %w", item.ProductId, err)
		}

		order.Items[i].Price = productPrice
		order.Items[i].TotalPrice = productPrice * float64(item.Quantity)
		totalSum += order.Items[i].TotalPrice
		order.Items[i].Id = item.ProductId
		order.Items[i].OrderId = orderId
		order.Items[i].CreatedAt = item.CreatedAt
	}

	// Insert the order
	orderQuery := `INSERT INTO "order" (id, user_id, total_price, delivery_status, longitude, latitude, address_name, created_at, updated_at) 
					  VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`

	_, err = tx.Exec(context.Background(), orderQuery, orderId, order.Order.UserId, totalSum, order.Order.DeliveryStatus, order.Order.Longitude, order.Order.Latitude, order.Order.AddressName)
	if err != nil {
		return &models.OrderCreateRequest{}, err
	}

	// Insert the order items
	itemQuery := `INSERT INTO "orderiteam" (id, quantity, order_id, product_id, price, total, created_at, updated_at) 
					 VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	for _, item := range order.Items {
		itemId := uuid.New().String()
		_, err = tx.Exec(context.Background(), itemQuery, itemId, item.Quantity, orderId, item.ProductId, item.Price, item.TotalPrice)
		if err != nil {
			return &models.OrderCreateRequest{}, err
		}
	}

	order.Order.Id = orderId
	order.Order.TotalPrice = totalSum

	return order, tx.Commit(context.Background())
}

func (r *OrderRepo) Update(ctx context.Context, id string, updatedOrder *models.Order) (*models.OrderCreateRequest, error) {
	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	// Update the order
	orderUpdateQuery := `
		UPDATE "order"
		SET user_id = $1, total_price = $2, status = $3, longitude = $4, latitude = $5, address_name = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7 AND deleted_at IS NULL
		RETURNING id, user_id, total_price, status, created_at, updated_at
	`
	var order models.Order
	err = tx.QueryRow(ctx, orderUpdateQuery, updatedOrder.UserId, updatedOrder.TotalPrice, updatedOrder.Status, updatedOrder.Longitude, updatedOrder.Latitude, updatedOrder.AddressName, id).Scan(
		&order.Id, &order.UserId, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update and retrieve order: %w", err)
	}

	// Update associated order items
	var orderItems []models.OrderItem
	for _, item := range updatedOrder.OrderItems {
		itemUpdateQuery := `
			UPDATE "orderiteam"
			SET product_id = $1, quantity = $2, price = $3, total = $4, updated_at = CURRENT_TIMESTAMP
			WHERE order_id = $5 AND id = $6 AND deleted_at IS NULL
			RETURNING id, product_id, order_id, quantity, price, total, created_at, updated_at
		`
		var updatedItem models.OrderItem
		err := tx.QueryRow(ctx, itemUpdateQuery, item.ProductId, item.Quantity, item.Price, item.TotalPrice, id, item.Id).Scan(
			&updatedItem.Id, &updatedItem.ProductId, &updatedItem.OrderId, &updatedItem.Quantity, &updatedItem.Price, &updatedItem.TotalPrice, &updatedItem.CreatedAt, &updatedItem.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update and retrieve order item with ID %s: %w", item.Id, err)
		}
		orderItems = append(orderItems, updatedItem)
	}

	// Return the updated order and its items
	orderCreateRequest := &models.OrderCreateRequest{
		Order: models.Order{
			Id:         order.Id,
			UserId:     order.UserId,
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			CreatedAt:  order.CreatedAt,
			UpdatedAt:  order.UpdatedAt,
		},
		Items: orderItems,
	}

	return orderCreateRequest, nil
}

func (o *OrderRepo) GetAll(ctx context.Context, request *models.GetAllOrdersRequest) (*[]models.OrderCreateRequest, error) {
	var (
		orders     []models.OrderCreateRequest
		created_at sql.NullString
		updated_at sql.NullString
	)

	// Query to retrieve all orders, sorted by the latest created orders at the top
	orderQuery := `
		SELECT id, user_id, total_price, delivery_status, status, longitude, latitude, address_name, created_at, updated_at
		FROM "order"
		ORDER BY created_at DESC
	`
	rows, err := o.db.Query(ctx, orderQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve orders: %w", err)
	}
	defer rows.Close()

	// Iterate over the retrieved orders
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.Id, &order.UserId, &order.TotalPrice, &order.DeliveryStatus, &order.Status, &order.Longitude, &order.Latitude, &order.AddressName, &created_at, &updated_at)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		// Query to retrieve order items for the current order
		orderItemQuery := `
			SELECT id, product_id, order_id, quantity, price, total, created_at, updated_at
			FROM "orderiteam"
			WHERE order_id = $1
		`
		itemRows, err := o.db.Query(ctx, orderItemQuery, order.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve items for order %s: %w", order.Id, err)
		}
		defer itemRows.Close()

		var orderItems []models.OrderItem
		for itemRows.Next() {
			var item models.OrderItem
			err = itemRows.Scan(&item.Id, &item.ProductId, &item.OrderId, &item.Quantity, &item.Price, &item.TotalPrice, &created_at, &updated_at)
			if err != nil {
				return nil, fmt.Errorf("failed to scan order item: %w", err)
			}
			orderItems = append(orderItems, models.OrderItem{
				Id:         item.Id,
				ProductId:  item.ProductId,
				OrderId:    item.OrderId,
				Quantity:   item.Quantity,
				Price:      item.Price,
				TotalPrice: item.TotalPrice,
				CreatedAt:  created_at.String,
				UpdatedAt:  updated_at.String,
			})
		}

		orders = append(orders, models.OrderCreateRequest{
			Order: models.Order{
				Id:             order.Id,
				UserId:         order.UserId,
				TotalPrice:     order.TotalPrice,
				DeliveryStatus: order.DeliveryStatus,
				Status:         order.Status,
				Longitude:      order.Longitude,
				Latitude:       order.Latitude,
				AddressName:    order.AddressName,
				CreatedAt:      created_at.String,
				UpdatedAt:      updated_at.String,
			},
			Items: orderItems,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &orders, nil
}

func (r *OrderRepo) GetOrder(ctx context.Context, id string) (*models.OrderCreateRequest, error) {
	var (
		created_at sql.NullString
		updated_at sql.NullString
	)

	orderQuery := `
		SELECT id, user_id, total_price, status, longitude, latitude, address_name, created_at, updated_at
		FROM "order"
		WHERE id = $1
	`

	var order models.Order
	err := r.db.QueryRow(ctx, orderQuery, id).Scan(&order.Id, &order.UserId, &order.TotalPrice, &order.Status, &order.Longitude, &order.Latitude, &order.AddressName, &created_at, &updated_at)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to retrieve order: %w", err)
	}

	order.CreatedAt = created_at.String
	order.UpdatedAt = updated_at.String

	orderItemQuery := `
		SELECT id, product_id, order_id, quantity, price, total, created_at, updated_at
		FROM "orderiteam"
		WHERE order_id = $1
	`
	itemRows, err := r.db.Query(ctx, orderItemQuery, order.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve items for order %s: %w", order.Id, err)
	}
	defer itemRows.Close()

	var orderItems []models.OrderItem
	for itemRows.Next() {
		var item models.OrderItem
		err = itemRows.Scan(&item.Id, &item.ProductId, &item.OrderId, &item.Quantity, &item.Price, &item.TotalPrice, &created_at, &updated_at)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		// Handle NullString for item timestamps
		item.CreatedAt = created_at.String
		item.UpdatedAt = updated_at.String

		orderItems = append(orderItems, item)
	}

	if err = itemRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order items: %w", err)
	}

	return &models.OrderCreateRequest{
		Order: order,
		Items: orderItems,
	}, nil
}

func (r *OrderRepo) Delete(ctx context.Context, id string) error {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	orderDeleteQuery := `
		UPDATE "order"
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`
	res, err := tx.Exec(ctx, orderDeleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("order not found or already deleted")
	}

	itemDeleteQuery := `
		UPDATE "orderiteam"
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE order_id = $1 AND deleted_at IS NULL
	`
	_, err = tx.Exec(ctx, itemDeleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to update order items: %w", err)
	}

	return nil
}

func (o *OrderRepo) ChangeOrderStatus(ctx context.Context, req *models.PatchOrderStatusRequest, orderId string) (string, error) {
	validStatuses := map[string]bool{
		"pending":   true,
		"confirmed": true,
		"picked_up": true,
		"delivered": true,
	}

	if !validStatuses[req.Status] {
		return "", fmt.Errorf("invalid status value: %s", req.Status)
	}

	tx, err := o.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	updateQuery := `UPDATE "order" SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err = tx.Exec(ctx, updateQuery, req.Status, orderId)
	if err != nil {
		return "", fmt.Errorf("failed to update order status: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return "Status changed successfully", nil
}
