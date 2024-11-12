package postgres

import (
	"context"
	"database/sql"

	// "errors"
	"fmt"
	"food/api/models"
	"food/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	// "golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewUser(db *pgxpool.Pool, log logger.LoggerI) UserRepo {
	return UserRepo{
		db:  db,
		log: log,
	}
}

func (c *UserRepo) GetByLogin(ctx context.Context, login string) (models.User, error) {
	var (
		name      sql.NullString
		phone     sql.NullString
		sex       sql.NullString
		email     sql.NullString
		createdat sql.NullString
		updatedat sql.NullString
	)

	query := `SELECT 
		id, 
		name, 
		phone,
		sex,
		email,
		created_at, 
		updated_at
		FROM "user" WHERE email = $1`

	row := c.db.QueryRow(ctx, query, login)

	user := models.User{}

	err := row.Scan(
		&user.Id,
		&name,
		&sex,
		&phone,
		&email,
		&createdat,
		&updatedat,
	)

	if err != nil {
		c.log.Error("failed to scan user by LOGIN from database", logger.Error(err))
		return models.User{}, err
	}

	user.Name = name.String
	user.Sex = sex.String
	user.Phone = phone.String
	user.Email = email.String
	user.Created_at = createdat.String
	user.Updated_at = updatedat.String

	return user, nil
}

// func (c *UserRepo) Login(ctx context.Context, login models.User) (string, error) {
// 	var hashedPass string

// 	query := `SELECT password
// 	FROM "user"
// 	WHERE phone = $1`

// 	err := c.db.QueryRow(ctx, query,
// 		login.Phone,
// 	).Scan(&hashedPass)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return "", errors.New("incorrect login")
// 		}
// 		c.log.Error("failed to get user password from database", logger.Error(err))
// 		return "", err
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(login.Password))
// 	if err != nil {
// 		return "", errors.New("password mismatch")
// 	}

// 	return "Logged in successfully", nil
// }

func (c *UserRepo) CheckPhoneNumberExist(ctx context.Context, id string) (models.User, error) {

	resp := models.User{}

	query := ` SELECT id FROM "user" WHERE phone = $1 `

	err := c.db.QueryRow(ctx, query, id).Scan(&resp.Id)
	if err != nil {
		return models.User{}, err
	}

	return resp, nil
}

func (u *UserRepo) Create(ctx context.Context, user *models.User) (*models.User, error) {

	id := uuid.New()
	query := `INSERT INTO "user" (
		id,
		email,
		name,
		sex,
		phone,
		created_at,
		updated_at)
		VALUES($1,$2,$3,$4,$5, CURRENT_TIMESTAMP,CURRENT_TIMESTAMP) 
	`

	_, err := u.db.Exec(context.Background(), query,
		id.String(),
		user.Email,
		user.Name,
		user.Sex,
		user.Phone,
	)

	if err != nil {
		return &models.User{}, err
	}
	return &models.User{
		Id:         id.String(),
		Email:      user.Email,
		Name:       user.Name,
		Sex:        user.Sex,
		Phone:      user.Phone,
		Created_at: user.Created_at,
		Updated_at: user.Updated_at,
	}, nil
}

func (u *UserRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	query := `UPDATE "user" SET 
		email=$1,
		name=$2,
		phone=$3,
		sex=$4,
		role=$5,
		updated_at=CURRENT_TIMESTAMP
		WHERE id = $6
	`
	_, err := u.db.Exec(context.Background(), query,
		user.Name,
		user.Email,
		user.Phone,
		user.Sex,
		user.Id,
	)
	if err != nil {
		return &models.User{}, err
	}
	return &models.User{
		Id:         user.Id,
		Name:       user.Name,
		Sex:        user.Sex,
		Email:      user.Email,
		Phone:      user.Phone,
		Created_at: user.Created_at,
		Updated_at: user.Updated_at,
	}, nil
}

func (u *UserRepo) GetAll(ctx context.Context, req *models.GetAllUsersRequest) (*models.GetAllUsersResponse, error) {
	var (
		resp   = &models.GetAllUsersResponse{}
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
		sex,
        email,
        phone,
        created_at,
        updated_at FROM "user"`+filter)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var (
			user       = models.User{}
			name       sql.NullString
			sex        sql.NullString
			email      sql.NullString
			phone      sql.NullString
			created_at sql.NullString
			updated_at sql.NullString
		)
		if err := rows.Scan(
			&resp.Count,
			&user.Id,
			&name,
			&sex,
			&email,
			&phone,
			&created_at,
			&updated_at); err != nil {
			return resp, err
		}

		resp.Users = append(resp.Users, models.User{
			Id:         user.Id,
			Name:       name.String,
			Sex:        user.Sex,
			Email:      email.String,
			Phone:      phone.String,
			Created_at: created_at.String,
			Updated_at: updated_at.String,
		})
	}
	return resp, nil
}

func (u *UserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	var (
		user       = models.User{}
		name       sql.NullString
		sex        sql.NullString
		email      sql.NullString
		phone      sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)
	if err := u.db.QueryRow(context.Background(), `SELECT id, name, sex, email, phone, created_at, updated_at FROM "user" WHERE id = $1`, id).Scan(
		&user.Id,
		&name,
		&sex,
		&email,
		&phone,
		&created_at,
		&updated_at,
	); err != nil {
		return &models.User{}, err
	}
	return &models.User{
		Id:         user.Id,
		Name:       name.String,
		Sex:        sex.String,
		Email:      email.String,
		Phone:      phone.String,
		Created_at: created_at.String,
		Updated_at: updated_at.String,
	}, nil
}

func (u *UserRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "user" WHERE id = $1`
	_, err := u.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) GetByPhone(ctx context.Context, number string) (*models.User, error) {
	var (
		admin      = models.User{}
		name       sql.NullString
		sex        sql.NullString
		email      sql.NullString
		phone      sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)
	if err := u.db.QueryRow(context.Background(), `SELECT id, name, sex, email, phone, created_at, updated_at FROM "user" WHERE phone = $1`, phone).Scan(
		&admin.Id,
		&name,
		&sex,
		&email,
		&phone,
		&created_at,
		&updated_at,
	); err != nil {
		return &models.User{}, err
	}
	return &models.User{
		Id:         admin.Id,
		Name:       name.String,
		Sex:        sex.String,
		Email:      email.String,
		Phone:      phone.String,
		Created_at: created_at.String,
		Updated_at: updated_at.String,
	}, nil
}
