package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"food/api/models"
	"food/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AdminRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewAdmin(db *pgxpool.Pool, log logger.LoggerI) AdminRepo {
	return AdminRepo{
		db:  db,
		log: log,
	}
}

func (c *AdminRepo) GetByLogin(ctx context.Context, login string) (models.Admin, error) {
	var (
		name      sql.NullString
		phone     sql.NullString
		email     sql.NullString
		createdat sql.NullString
		updatedat sql.NullString
	)

	query := `SELECT 
		id, 
		name, 
		phone,
		email,
		created_at, 
		updated_at,
		password
		FROM "admin" WHERE phone = $1`

	row := c.db.QueryRow(ctx, query, login)

	admin := models.Admin{}

	err := row.Scan(
		&admin.Id,
		&name,
		&phone,
		&email,
		&createdat,
		&updatedat,
		&admin.Password,
	)

	if err != nil {
		c.log.Error("failed to scan user by LOGIN from database", logger.Error(err))
		return models.Admin{}, err
	}

	admin.Name = name.String
	admin.Phone = phone.String
	admin.Email = email.String
	admin.Created_at = createdat.String
	admin.Updated_at = updatedat.String

	return admin, nil
}

func (c *AdminRepo) Login(ctx context.Context, login models.Admin) (string, error) {
	var hashedPass string

	query := `SELECT password
	FROM "admin"
	WHERE phone = $1`

	err := c.db.QueryRow(ctx, query,
		login.Phone,
	).Scan(&hashedPass)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("incorrect login")
		}
		c.log.Error("failed to get user password from database", logger.Error(err))
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(login.Password))
	if err != nil {
		return "", errors.New("password mismatch")
	}

	return "Logged in successfully", nil
}

func (c *AdminRepo) CheckPhoneNumberExist(ctx context.Context, id string) (models.Admin, error) {

	resp := models.Admin{}

	query := ` SELECT id FROM "admin" WHERE phone = $1 `

	err := c.db.QueryRow(ctx, query, id).Scan(&resp.Id)
	if err != nil {
		return models.Admin{}, err
	}

	return resp, nil
}

func (u *AdminRepo) Create(ctx context.Context, user *models.Admin) (*models.Admin, error) {

	id := uuid.New()
	query := `INSERT INTO "admin" (
		id,
		email,
		name,
		phone,
		password,
		created_at,
		updated_at)
		VALUES($1,$2,$3,$4,$5, CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) 
	`

	_, err := u.db.Exec(context.Background(), query,
		id.String(),
		user.Email,
		user.Name,
		user.Phone,
		user.Password,
	)

	if err != nil {
		return &models.Admin{}, err
	}
	return &models.Admin{
		Id:         id.String(),
		Email:      user.Email,
		Name:       user.Name,
		Phone:      user.Phone,
		Password:   user.Password,
		Created_at: user.Created_at,
		Updated_at: user.Updated_at,
	}, nil
}

func (u *AdminRepo) Update(ctx context.Context, admin *models.Admin) (*models.Admin, error) {
	query := `UPDATE "admin" SET 
		email=$1,
		name=$2,
		phone=$3,
		password=$4,
		updated_at=CURRENT_TIMESTAMP
		WHERE id = $5
	`
	_, err := u.db.Exec(context.Background(), query,
		admin.Name,
		admin.Email,
		admin.Phone,
		admin.Password,
		admin.Id,
	)
	if err != nil {
		return &models.Admin{}, err
	}
	return &models.Admin{
		Id:         admin.Id,
		Name:       admin.Name,
		Email:      admin.Email,
		Phone:      admin.Phone,
		Password:   admin.Password,
		Created_at: admin.Created_at,
		Updated_at: admin.Updated_at,
	}, nil
}

func (u *AdminRepo) GetAll(ctx context.Context, req *models.GetAllAdminsRequest) (*models.GetAllAdminsResponse, error) {
	var (
		resp   = &models.GetAllAdminsResponse{}
		filter = ""
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(` WHERE (email ILIKE '%%%v%%' OR phone ILIKE '%%%v%%') `, req.Search, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter: ", filter)

	rows, err := u.db.Query(context.Background(), `SELECT count(id) OVER(),
        id,
		name,
        email,
        phone,
        password,
        created_at,
        updated_at FROM "admin"`+filter)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var (
			admin      = models.Admin{}
			name       sql.NullString
			email      sql.NullString
			phone      sql.NullString
			password   sql.NullString
			created_at sql.NullString
			updated_at sql.NullString
		)
		if err := rows.Scan(
			&resp.Count,
			&admin.Id,
			&name,
			&email,
			&phone,
			&password,
			&created_at,
			&updated_at); err != nil {
			return resp, err
		}

		resp.Admins = append(resp.Admins, models.Admin{
			Id:         admin.Id,
			Name:       name.String,
			Email:      email.String,
			Phone:      phone.String,
			Password:   password.String,
			Created_at: created_at.String,
			Updated_at: updated_at.String,
		})
	}
	return resp, nil
}

func (u *AdminRepo) GetByID(ctx context.Context, id string) (*models.Admin, error) {
	var (
		admin      = models.Admin{}
		name       sql.NullString
		email      sql.NullString
		phone      sql.NullString
		password   sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)
	if err := u.db.QueryRow(context.Background(), `SELECT id, name, email, phone, password, created_at, updated_at FROM "admin" WHERE id = $1`, id).Scan(
		&admin.Id,
		&name,
		&email,
		&phone,
		&password,
		&created_at,
		&updated_at,
	); err != nil {
		return &models.Admin{}, err
	}
	return &models.Admin{
		Id:         admin.Id,
		Name:       name.String,
		Email:      email.String,
		Phone:      phone.String,
		Password:   password.String,
		Created_at: created_at.String,
		Updated_at: updated_at.String,
	}, nil
}

func (u *AdminRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "admin" WHERE id = $1`
	_, err := u.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *AdminRepo) GetByPhone(ctx context.Context, number string) (*models.Admin, error) {
	var (
		admin      = models.Admin{}
		name       sql.NullString
		email      sql.NullString
		phone      sql.NullString
		password   sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)
	if err := u.db.QueryRow(context.Background(), `SELECT id, name, email, phone, password, created_at, updated_at FROM "admin" WHERE phone = $1`, number).Scan(
		&admin.Id,
		&name,
		&email,
		&phone,
		&password,
		&created_at,
		&updated_at,
	); err != nil {
		return &models.Admin{}, err
	}
	return &models.Admin{
		Id:         admin.Id,
		Name:       name.String,
		Email:      email.String,
		Phone:      phone.String,
		Password:   password.String,
		Created_at: created_at.String,
		Updated_at: updated_at.String,
	}, nil
}
