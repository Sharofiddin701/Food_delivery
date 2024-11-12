package postgres

import (
	"context"
	"fmt"
	"food/config"
	"food/pkg/logger"
	"food/storage"
	"food/storage/redis"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	Pool    *pgxpool.Pool
	db      *pgxpool.Pool
	admin   *AdminRepo
	combo   *ComboRepo
	payment *PaymentRepo
	// redis              storage.IRedisStorage
	log  logger.LoggerI
	user *UserRepo
	// auth               *AuthRepo
	branch             *BranchRepo
	banner             *BannerRepo
	category           *CategoryRepo
	order              *OrderRepo
	product            *ProductRepo
	notification       *NotificationRepo
	delivery_history   *DeliveryHistoryRepo
	courier_assignment *CourierAssignmentRepo
	cfg                config.Config
}

// CloseDB implements storage.IStorage.
func (s Store) CloseDB() {
	s.Pool.Close()
}

func NewConnectionPostgres(cfg *config.Config) (storage.IStorage, error) {
	connect, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d ",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	))

	if err != nil {
		return nil, err
	}
	connect.MaxConns = 100

	pgxpool, err := pgxpool.ConnectConfig(context.Background(), connect)
	if err != nil {
		return nil, err
	}
	var loggerLevel = new(string)
	log := logger.NewLogger("app", *loggerLevel)
	defer func() {
		err := logger.Cleanup(log)
		if err != nil {
			return
		}
	}()
	return &Store{
		db:  pgxpool,
		log: logger.NewLogger("app", *loggerLevel),
	}, nil
}

func (s *Store) Payment() storage.IPaymentStorage {
	if s.db == nil {
		s.payment = &PaymentRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.payment
}

func (s *Store) User() storage.IUserStorage {
	if s.user == nil {
		s.user = &UserRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.user
}

func (s *Store) Combo() storage.IComboStorage {
	if s.combo == nil {
		s.combo = &ComboRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.combo
}

func (s *Store) Admin() storage.IAdminStorage {
	if s.admin == nil {
		s.admin = &AdminRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.admin
}

// Auth implements storage.IStorage.
// func (s *Store) Auth() storage.IAuthStorage {
// 	if s.auth == nil {
// 		s.auth = &AuthRepo{
// 			user:  s.user,
// 			db:    s.db,
// 			log:   s.log,
// 			redis: s.redis,
// 		}
// 	}
// 	return s.auth
// }

// Redis implements storage.IStorage.
func (s *Store) Redis() storage.IRedisStorage {
	return redis.New(s.cfg)
}

func (s *Store) Branch() storage.IBranchStorage {
	if s.branch == nil {
		s.branch = &BranchRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.branch
}

func (s *Store) Banner() storage.IBannerStorage {
	if s.banner == nil {
		s.banner = &BannerRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.banner
}

func (s *Store) Category() storage.ICategoryStorage {
	if s.category == nil {
		s.category = &CategoryRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.category
}

func (s *Store) Order() storage.IOrderStorage {
	if s.order == nil {
		s.order = &OrderRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.order
}

func (s *Store) Product() storage.IProductStorage {
	if s.product == nil {
		s.product = &ProductRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.product
}

// CourierAssignment implements storage.IStorage.
func (s *Store) CourierAssignment() storage.ICourierAssignmentStorage {
	if s.courier_assignment == nil {
		s.courier_assignment = &CourierAssignmentRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.courier_assignment
}

// DeliveryHistory implements storage.IStorage.
func (s *Store) DeliveryHistory() storage.IDeliveryHistoryStorage {
	if s.delivery_history == nil {
		s.delivery_history = &DeliveryHistoryRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.delivery_history
}

// Notification implements storage.IStorage.
func (s *Store) Notification() storage.INotificationStorage {
	if s.notification == nil {
		s.notification = &NotificationRepo{
			db:  s.db,
			log: s.log,
		}
	}
	return s.notification
}
