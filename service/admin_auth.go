package service

import (
	"context"
	"errors"
	"fmt"
	"food/api/models"
	"food/config"
	"food/pkg"
	"food/pkg/jwt"
	"food/pkg/logger"

	// "food/pkg/password"

	"food/storage"
	"time"

	"github.com/go-redis/redis"
)

type adminAuthService struct {
	storage storage.IStorage
	log     logger.LoggerI
	redis   storage.IRedisStorage
}

func NewAuthAdminService(storage storage.IStorage, log logger.LoggerI, redis storage.IRedisStorage) adminAuthService {
	return adminAuthService{
		storage: storage,
		log:     log,
		redis:   redis,
	}
}

func (a adminAuthService) AdminLogin(ctx context.Context, loginRequest models.AdminLoginRequest) (models.AdminLoginResponse, error) {
	fmt.Println(" loginRequest.Login: ", loginRequest.Login)
	admin, err := a.storage.Admin().GetByLogin(ctx, loginRequest.Login)
	if err != nil {
		a.log.Error("error while getting user credentials by login", logger.Error(err))
		return models.AdminLoginResponse{}, err
	}

	// if err = password.CompareHashAndPassword(user.Password, loginRequest.Password); err != nil {
	// 	a.log.Error("error while comparing password", logger.Error(err))
	// 	return models.UserLoginResponse{}, err
	// }

	m := make(map[interface{}]interface{})

	m["user_id"] = admin.Id
	m["user_role"] = config.USER_ROLE

	// accessToken, refreshToken, err := jwt.GenJWT(m)
	// if err != nil {
	// 	a.log.Error("error while generating tokens for user login", logger.Error(err))
	// 	return models.AdminLoginResponse{}, err
	// }

	return models.AdminLoginResponse{
		// AccessToken:  accessToken,
		// RefreshToken: refreshToken,
		Id:           admin.Id,
		Phone:        admin.Phone,
	}, nil
}

func (a adminAuthService) AdminRegister(ctx context.Context, loginRequest models.AdminRegisterRequest) error {
	fmt.Println(" loginRequest.Login: ", loginRequest.MobilePhone)

	otpCode := pkg.GenerateOTP()

	msg := fmt.Sprintf("iBron ilovasi ro‘yxatdan o‘tish uchun tasdiqlash kodi: %v", otpCode)

	err := a.redis.SetX(ctx, loginRequest.MobilePhone, otpCode, time.Minute*2)
	if err != nil {
		a.log.Error("error while setting otpCode to redis user register", logger.Error(err))
		return err
	}

	err = pkg.SendSms(loginRequest.MobilePhone, msg)
	if err != nil {
		a.log.Error("error while sending otp code to user register", logger.Error(err))
		return err
	}
	return nil
}

func (a adminAuthService) AdminRegisterConfirm(ctx context.Context, req models.UserRegisterConfRequest) (models.UserLoginResponse, error) {
	resp := models.UserLoginResponse{}

	otp, err := a.redis.Get(ctx, req.MobilePhone)
	if err != nil {
		a.log.Error("error while getting otp code for user register confirm", logger.Error(err))
		return resp, err
	}
	if req.Otp != otp {
		a.log.Error("incorrect otp code for user register confirm", logger.Error(err))
		return resp, errors.New("incorrect otp code")
	}
	req.User.Phone = req.MobilePhone
	id, err := a.storage.User().Create(ctx, req.User)
	if err != nil {
		a.log.Error("error while creating user", logger.Error(err))
		return resp, err
	}
	var m = make(map[interface{}]interface{})

	m["user_id"] = id
	m["user_role"] = config.USER_ROLE

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		a.log.Error("error while generating tokens for user register confirm", logger.Error(err))
		return resp, err
	}
	resp.AccessToken = accessToken
	resp.RefreshToken = refreshToken

	return resp, nil
}

func (a adminAuthService) AdminLoginByPhoneConfirm(ctx context.Context, req models.UserLoginPhoneConfirmRequest) (models.UserLoginResponse, error) {
	resp := models.UserLoginResponse{}

	storedOTP, err := a.redis.Get(ctx, req.MobilePhone)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			a.log.Error("OTP code not found or expired", logger.Error(err))
			return resp, errors.New("OTP kod topilmadi yoki muddati tugagan")
		}
		a.log.Error("error while getting OTP code from redis", logger.Error(err))
		return resp, errors.New("tizim xatosi yuz berdi")
	}

	if req.SmsCode != storedOTP {
		a.log.Error("incorrect OTP code", logger.Error(errors.New("OTP code mismatch")))
		return resp, errors.New("noto'g'ri OTP kod")
	}

	err = a.redis.Del(ctx, req.MobilePhone)
	if err != nil {
		a.log.Error("error while deleting OTP from redis", logger.Error(err))
		return resp, err
	}
	user, err := a.storage.Admin().CheckPhoneNumberExist(ctx, req.MobilePhone)
	if err != nil {
		a.log.Error("error while getting user by phone number", logger.Error(err))
		return resp, err
	}

	resp.Phone = req.MobilePhone
	resp.Id = user.Id

	return resp, nil
}
