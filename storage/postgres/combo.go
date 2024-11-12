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

type ComboRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewCombo(db *pgxpool.Pool, log logger.LoggerI) ComboRepo {
	return ComboRepo{
		db:  db,
		log: log,
	}
}
func (c *ComboRepo) Create(ctx context.Context, combo *models.ComboCreateRequest) (*models.ComboCreateRequest, error) {
	tx, err := c.db.Begin(context.Background())
	if err != nil {
		return &models.ComboCreateRequest{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	// Generate a new UUID for the combo
	comboId := uuid.New().String()

	// var totalSum float64
	for i, item := range combo.Items {
		if item.Quantity <= 0 {
			return &models.ComboCreateRequest{}, fmt.Errorf("quantity must be greater than 0 for product %s", item.ProductId)
		}

		var productPrice float64
		productQuery := `SELECT price FROM "product" WHERE id = $1`
		err = c.db.QueryRow(context.Background(), productQuery, item.ProductId).Scan(&productPrice)
		if err != nil {
			return &models.ComboCreateRequest{}, fmt.Errorf("failed to retrieve price for product %s: %w", item.ProductId, err)
		}

		combo.Items[i].Price = productPrice
		combo.Items[i].TotalPrice = productPrice * float64(item.Quantity)
		// totalSum += combo.Items[i].TotalPrice
		combo.Items[i].ComboId = comboId
		combo.Items[i].CreatedAt = item.CreatedAt
	}

	// Insert the combo
	comboQuery := `INSERT INTO "combo" (id, name, price, description, created_at) 
					  VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) RETURNING id`

	_, err = tx.Exec(context.Background(), comboQuery, comboId, combo.Combo.Name, combo.Combo.Price, combo.Combo.Description)
	if err != nil {
		return &models.ComboCreateRequest{}, err
	}

	// Insert the combo items
	itemQuery := `INSERT INTO "combo_items" (id, quantity, combo_id, product_id, price, total_price, created_at) 
					 VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)`

	for _, item := range combo.Items {
		itemId := uuid.New().String()
		_, err = tx.Exec(context.Background(), itemQuery, itemId, item.Quantity, comboId, item.ProductId, item.Price, item.TotalPrice)
		for i := range combo.Items {
			combo.Items[i].Id = itemId
			combo.Items[i].CreatedAt = item.CreatedAt
		}
		if err != nil {
			return &models.ComboCreateRequest{}, err
		}
	}

	combo.Combo.Id = comboId
	// combo.Combo.TotalPrice = totalSum

	return combo, tx.Commit(context.Background())
}

func (c *ComboRepo) GetAll(ctx context.Context, request *models.GetAllCombosRequest) (*[]models.ComboCreateRequest, error) {
	var (
		combos     []models.ComboCreateRequest
		created_at sql.NullString
		updated_at sql.NullString
	)

	comboQuery := `
		SELECT id, name, description, price, created_at, updated_at
		FROM "combo"
		ORDER BY created_at DESC
	`
	rows, err := c.db.Query(ctx, comboQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve combos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var combo models.Combo
		err = rows.Scan(&combo.Id, &combo.Name, &combo.Description, &combo.Price, &created_at, &updated_at)
		if err != nil {
			return nil, fmt.Errorf("failed to scan combo: %w", err)
		}

		comboItemQuery := `
			SELECT id, product_id, combo_id, quantity, price, total_price, created_at, updated_at
			FROM "combo_items"
			WHERE combo_id = $1
		`
		itemRows, err := c.db.Query(ctx, comboItemQuery, combo.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve items for combo %s: %w", combo.Id, err)
		}
		defer itemRows.Close()

		var comboItems []models.ComboItem
		for itemRows.Next() {
			var item models.ComboItem
			err = itemRows.Scan(&item.Id, &item.ProductId, &item.ComboId, &item.Quantity, &item.Price, &item.TotalPrice, &created_at, &updated_at)
			if err != nil {
				return nil, fmt.Errorf("failed to scan combo item: %w", err)
			}
			comboItems = append(comboItems, models.ComboItem{
				Id:         item.Id,
				ProductId:  item.ProductId,
				ComboId:    item.ComboId,
				Quantity:   item.Quantity,
				Price:      item.Price,
				TotalPrice: item.TotalPrice,
				CreatedAt:  created_at.String,
				UpdatedAt:  updated_at.String,
			})
		}

		combos = append(combos, models.ComboCreateRequest{
			Combo: models.Combo{
				Id:          combo.Id,
				Name:        combo.Name,
				Description: combo.Description,
				Price:       combo.Price,
				CreatedAt:   created_at.String,
				UpdatedAt:   updated_at.String,
			},
			Items: comboItems,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &combos, nil
}

func (r *ComboRepo) GetCombo(ctx context.Context, id string) (*models.ComboCreateRequest, error) {
	var (
		created_at sql.NullString
		updated_at sql.NullString
	)

	comboQuery := `
		SELECT id, name, description, price, created_at, updated_at
		FROM "combo"
		WHERE id = $1
	`

	var combo models.Combo
	err := r.db.QueryRow(ctx, comboQuery, id).Scan(&combo.Id, &combo.Name, &combo.Description, &combo.Price, &created_at, &updated_at)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("combo not found")
		}
		return nil, fmt.Errorf("failed to retrieve combo: %w", err)
	}

	combo.CreatedAt = created_at.String
	combo.UpdatedAt = updated_at.String

	comboItemQuery := `
		SELECT id, combo_id, product_id, quantity, price, total_price, created_at, updated_at
		FROM "combo_items"
		WHERE combo_id = $1
	`
	itemRows, err := r.db.Query(ctx, comboItemQuery, combo.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve items for combo %s: %w", combo.Id, err)
	}
	defer itemRows.Close()

	var comboItems []models.ComboItem
	for itemRows.Next() {
		var item models.ComboItem
		err = itemRows.Scan(&item.Id, &item.ComboId, &item.ProductId, &item.Quantity, &item.Price, &item.TotalPrice, &created_at, &updated_at)
		if err != nil {
			return nil, fmt.Errorf("failed to scan combo item: %w", err)
		}

		item.CreatedAt = created_at.String
		item.UpdatedAt = updated_at.String

		comboItems = append(comboItems, item)
	}

	if err = itemRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating combo items: %w", err)
	}

	return &models.ComboCreateRequest{
		Combo: combo,
		Items: comboItems,
	}, nil
}

func (r *ComboRepo) Update(ctx context.Context, id string, updatedCombo *models.Combo) (*models.ComboCreateRequest, error) {

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

	// Update the combo
	comboUpdateQuery := `
		UPDATE "combo"
		SET name = $1, description = $2, price = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 
		RETURNING id, name, description, price, created_at, updated_at
	`
	var combo models.Combo
	err = tx.QueryRow(ctx, comboUpdateQuery, updatedCombo.Name, updatedCombo.Description, updatedCombo.Price, id).Scan(
		&combo.Id, &combo.Name, &combo.Description, &combo.Price, &combo.CreatedAt, &combo.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update and retrieve combo: %w", err)
	}

	// Update associated combo items
	var comboItems []models.ComboItem
	for _, item := range updatedCombo.ComboItems {
		if item.Id != "" {
			// Update existing combo item
			itemUpdateQuery := `
				UPDATE "comboitem"
				SET product_id = $1, quantity = $2, price = $3, total_price = $4, updated_at = CURRENT_TIMESTAMP
				WHERE combo_id = $5 AND id = $6 
				RETURNING id, product_id, combo_id, quantity, price, total_price, created_at, updated_at
			`
			var updatedItem models.ComboItem
			err := tx.QueryRow(ctx, itemUpdateQuery, item.ProductId, item.Quantity, item.Price, item.TotalPrice, id, item.Id).Scan(
				&updatedItem.Id, &updatedItem.ProductId, &updatedItem.ComboId, &updatedItem.Quantity, &updatedItem.Price, &updatedItem.TotalPrice, &updatedItem.CreatedAt, &updatedItem.UpdatedAt,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to update and retrieve combo item with ID %s: %w", item.Id, err)
			}
			comboItems = append(comboItems, updatedItem)
		} else {
			// Insert new combo item
			itemInsertQuery := `
				INSERT INTO "comboitem" (id, combo_id, product_id, quantity, price, total_price, created_at, updated_at)
				VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
				RETURNING id, product_id, combo_id, quantity, price, total_price, created_at, updated_at
			`
			var newItem models.ComboItem
			err := tx.QueryRow(ctx, itemInsertQuery, id, item.ProductId, item.Quantity, item.Price, item.TotalPrice).Scan(
				&newItem.Id, &newItem.ProductId, &newItem.ComboId, &newItem.Quantity, &newItem.Price, &newItem.TotalPrice, &newItem.CreatedAt, &newItem.UpdatedAt,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to insert new combo item: %w", err)
			}
			comboItems = append(comboItems, newItem)
		}
	}

	comboCreateRequest := &models.ComboCreateRequest{
		Combo: models.Combo{
			Id:          combo.Id,
			Name:        combo.Name,
			Description: combo.Description,
			Price:       combo.Price,
			CreatedAt:   combo.CreatedAt,
			UpdatedAt:   combo.UpdatedAt,
		},
		Items: comboItems,
	}

	return comboCreateRequest, nil
}
