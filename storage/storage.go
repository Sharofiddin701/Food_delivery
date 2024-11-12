package storage

import (
	"context"
	"food/api/models"
	"time"
)

type IStorage interface {
	CloseDB()
	Admin() IAdminStorage
	User() IUserStorage
	Combo() IComboStorage
	Branch() IBranchStorage
	Banner() IBannerStorage
	Category() ICategoryStorage
	Product() IProductStorage
	Payment() IPaymentStorage
	Order() IOrderStorage
	CourierAssignment() ICourierAssignmentStorage
	Notification() INotificationStorage
	DeliveryHistory() IDeliveryHistoryStorage
	Redis() IRedisStorage
}

type IPaymentStorage interface {
	Create(context.Context, *models.Payment) (*models.Payment, error)
	GetByID(ctx context.Context, id string) (*models.Payment, error)
	GetAll(ctx context.Context, request *models.GetAllPaymentsRequest) (*models.GetAllPaymentsResponse, error)
	Update(context.Context, *models.Payment) (*models.Payment, error)
	Delete(context.Context, string) error
}

type IUserStorage interface {
	Create(context.Context, *models.User) (*models.User, error)
	GetAll(ctx context.Context, request *models.GetAllUsersRequest) (*models.GetAllUsersResponse, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(context.Context, *models.User) (*models.User, error)
	Delete(context.Context, string) error
	GetByLogin(ctx context.Context, login string) (models.User, error)
	CheckPhoneNumberExist(ctx context.Context, id string) (models.User, error)
	GetByPhone(ctx context.Context, number string) (*models.User, error)
}

type IComboStorage interface {
	Create(context.Context, *models.ComboCreateRequest) (*models.ComboCreateRequest, error)
	GetAll(ctx context.Context, request *models.GetAllCombosRequest) (*[]models.ComboCreateRequest, error)
	GetCombo(ctx context.Context, id string) (*models.ComboCreateRequest, error)
	Update(ctx context.Context, id string, updatedCombo *models.Combo) (*models.ComboCreateRequest, error)
	// Delete(context.Context, string) error
}

type IAdminStorage interface {
	Create(context.Context, *models.Admin) (*models.Admin, error)
	GetAll(ctx context.Context, request *models.GetAllAdminsRequest) (*models.GetAllAdminsResponse, error)
	GetByID(ctx context.Context, id string) (*models.Admin, error)
	Update(context.Context, *models.Admin) (*models.Admin, error)
	Delete(context.Context, string) error
	GetByLogin(ctx context.Context, login string) (models.Admin, error)
	CheckPhoneNumberExist(ctx context.Context, id string) (models.Admin, error)
	GetByPhone(ctx context.Context, number string) (*models.Admin, error)
}

type IBannerStorage interface {
	Create(context.Context, *models.Banner) (*models.Banner, error)
	GetAll(ctx context.Context, request *models.GetAllBannerRequest) (*models.GetAllBannerResponse, error)
	Delete(ctx context.Context, id string) error
}

type IBranchStorage interface {
	Create(ctx context.Context, branch *models.Branch) (*models.Branch, error)
	GetAll(ctx context.Context, request *models.GetAllBranchesRequest) (*models.GetAllBranchesResponse, error)
	GetByID(ctx context.Context, id string) (*models.Branch, error)
	Update(ctx context.Context, branch *models.Branch) (*models.Branch, error)
	Delete(ctx context.Context, id string) error
}

type ICategoryStorage interface {
	Create(context.Context, *models.Category) (*models.Category, error)
	GetAll(ctx context.Context, request *models.GetAllCategoriesRequest) (*models.GetAllCategoriesResponse, error)
	GetByID(ctx context.Context, id string) (*models.Category, error)
	Update(context.Context, *models.Category) (*models.Category, error)
	Delete(context.Context, string) error
}

type IProductStorage interface {
	Create(context.Context, *models.Product) (*models.Product, error)
	GetAll(ctx context.Context, request *models.GetAllProductsRequest) (*models.GetAllProductsResponse, error)
	GetByID(ctx context.Context, id string) (*models.Product, error)
	Update(context.Context, *models.Product) (*models.Product, error)
	Delete(context.Context, string) error
}

type IOrderStorage interface {
	Create(ctx context.Context, request *models.OrderCreateRequest) (*models.OrderCreateRequest, error)
	GetAll(ctx context.Context, request *models.GetAllOrdersRequest) (*[]models.OrderCreateRequest, error)
	GetOrder(ctx context.Context, id string) (*models.OrderCreateRequest, error)
	Update(ctx context.Context, id string, updatedOrder *models.Order) (*models.OrderCreateRequest, error)
	Delete(ctx context.Context, id string) error
	ChangeOrderStatus(ctx context.Context, req *models.PatchOrderStatusRequest, orderId string) (string, error)
}

type ICourierAssignmentStorage interface {
	Create(context.Context, *models.CourierAssignment) (*models.CourierAssignment, error)
	GetAll(ctx context.Context, request *models.GetAllCourierAssignmentsRequest) (*models.GetAllCourierAssignmentsResponse, error)
	GetByID(ctx context.Context, id string) (*models.CourierAssignment, error)
	Update(context.Context, *models.CourierAssignment) (*models.CourierAssignment, error)
	Delete(context.Context, string) error
}

type INotificationStorage interface {
	Create(context.Context, *models.Notification) (*models.Notification, error)
	GetAll(ctx context.Context, request *models.GetAllNotificationsRequest) (*models.GetAllNotificationsResponse, error)
	GetByID(ctx context.Context, id string) (*models.Notification, error)
	Update(context.Context, *models.Notification) (*models.Notification, error)
	Delete(context.Context, string) error
}

type IDeliveryHistoryStorage interface {
	Create(context.Context, *models.DeliveryHistory) (*models.DeliveryHistory, error)
	GetAll(ctx context.Context, request *models.GetAllDeliveryHistoriesRequest) (*models.GetAllDeliveryHistoriesResponse, error)
	GetByID(ctx context.Context, id string) (*models.DeliveryHistory, error)
	Update(context.Context, *models.DeliveryHistory) (*models.DeliveryHistory, error)
	Delete(context.Context, string) error
}

type IRedisStorage interface {
	SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Del(ctx context.Context, key string) error
}

// type IAuthStorage interface {
// 	UserRegister(ctx context.Context, loginRequest models.UserRegisterRequest) error
// 	UserRegisterConfirm(ctx context.Context, req models.UserRegisterConfRequest) (models.UserLoginResponse, error)
// }
