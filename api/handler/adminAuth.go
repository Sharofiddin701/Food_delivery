package handler

import (
	"fmt"
	"food/api/models"
	check "food/pkg/validation"

	// check "food/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminLogin godoc
// @Router       /food/api/v1/admin/login [POST]
// @Summary      Admin login
// @Description  Login to Food_delivery
// @Tags         admin_auth
// @Accept       json
// @Produce      json
// @Param        login body models.AdminLoginRequest true "login"
// @Success      201  {object}  models.AdminLoginResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) AdminLogin(c *gin.Context) {
	loginReq := models.AdminLoginRequest{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		handleResponseLog(c, h.log, "error while binding body", http.StatusInternalServerError, err)
		return
	}

	fmt.Println("loginReq: ", loginReq)

	if err := check.ValidatePhoneNumber(loginReq.Login); err != nil {
		handleResponseLog(c, h.log, "error while validating phone number", http.StatusBadRequest, err.Error())
		return
	}

	password, err := h.storage.Admin().GetByPhone(c.Request.Context(), loginReq.Login)
	if err != nil {
		handleResponseLog(c, h.log, "error while gettingByPhone", http.StatusInternalServerError, err.Error())
		return
	}

	if loginReq.Password != password.Password {
		handleResponseLog(c, h.log, "password does not match", http.StatusBadRequest, "")
		return
	}

	loginResp, err := h.service.AdminAuth().AdminLogin(c.Request.Context(), loginReq)
	if err != nil {
		handleResponseLog(c, h.log, "unauthorized", http.StatusUnauthorized, err.Error())
		return
	}

	handleResponseLog(c, h.log, "Success", http.StatusOK, loginResp)
}

// // AdminRegister godoc
// // @Router       /food/api/v1/admin/sendcode [POST]
// // @Summary      Admin register
// // @Description  Registering to Food_delivery
// // @Tags         admin_auth
// // @Accept       json
// // @Produce      json
// // @Param        register body models.AdminRegisterRequest true "register"
// // @Success      201  {object}  models.Response
// // @Failure      400  {object}  models.Response
// // @Failure      404  {object}  models.Response
// // @Failure      500  {object}  models.Response
// func (h *Handler) AdminRegister(c *gin.Context) {
// 	loginReq := models.AdminRegisterRequest{}

// 	if err := c.ShouldBindJSON(&loginReq); err != nil {
// 		handleResponseLog(c, h.log, "error while binding body", http.StatusBadRequest, err)
// 		return
// 	}
// 	fmt.Println("loginReq: ", loginReq)

// 	// if err := check.ValidateEmailAddress(loginReq.MobilePhone); err != nil {
// 	// 	handleResponseLog(c, h.log, "error while validating email" + loginReq.MobilePhone, http.StatusBadRequest, err.Error())
// 	// 	return
// 	// }

// 	// if err := check.CheckEmailExists(loginReq.MobilePhone); err != nil {
// 	// 	handleResponseLog(c, h.log, "error email does not exist" + loginReq.MobilePhone, http.StatusBadRequest, err.Error())
// 	// }

// 	err := h.service.AdminAuth().AdminRegister(c.Request.Context(), loginReq)
// 	if err != nil {
// 		handleResponseLog(c, h.log, "", http.StatusInternalServerError, err)
// 		return
// 	}

// 	handleResponseLog(c, h.log, "Otp sent successfull", http.StatusOK, "")
// }

// // AdminRegisterConfirm godoc
// // @Router       /food/api/v1/admin/verifycode [POST]
// // @Summary      Admin register
// // @Description  Registering to Food_delivery
// // @Tags         admin_auth
// // @Accept       json
// // @Produce      json
// // @Param        register body models.AdminRegisterConfRequest true "register"
// // @Success      201  {object}  models.AdminLoginResponse
// // @Failure      400  {object}  models.Response
// // @Failure      404  {object}  models.Response
// // @Failure      500  {object}  models.Response
// func (h *Handler) AdminRegisterConfirm(c *gin.Context) {
// 	req := models.UserRegisterConfRequest{}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		handleResponseLog(c, h.log, "error while binding body", http.StatusBadRequest, err)
// 		return
// 	}
// 	fmt.Println("req: ", req)

// 	//TODO: need validate login & password

// 	confResp, err := h.service.AdminAuth().AdminRegisterConfirm(c.Request.Context(), req)
// 	if err != nil {
// 		handleResponseLog(c, h.log, "error while confirming", http.StatusUnauthorized, err.Error())
// 		return
// 	}

// 	handleResponseLog(c, h.log, "Success", http.StatusOK, confResp)

// }

// // AdminLoginByPhoneConfirm godoc
// // @Router       /food/api/v1/admin/byphoneconfirm [POST]
// // @Summary      Admin login by phone confirmation
// // @Description  Login to the system using phone number and OTP
// // @Tags         admin_auth
// // @Accept       json
// // @Produce      json
// // @Param        login body models.AdminLoginPhoneConfirmRequest true "login"
// // @Success      200  {object}  models.AdminLoginResponse
// // @Failure      400  {object}  models.Response
// // @Failure      401  {object}  models.Response
// // @Failure      500  {object}  models.Response
// func (h *Handler) AdminLoginByPhoneConfirm(c *gin.Context) {
// 	var req models.UserLoginPhoneConfirmRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		h.log.Error("error while binding request body: " + err.Error())
// 		c.JSON(http.StatusBadRequest, models.Response{
// 			StatusCode:  http.StatusBadRequest,
// 			Description: err.Error(),
// 		})
// 		return
// 	}
// 	resp, err := h.service.AdminAuth().AdminLoginByPhoneConfirm(c.Request.Context(), req)
// 	if err != nil {
// 		statusCode := http.StatusInternalServerError
// 		message := "INTERNALSERVERERROR"

// 		if err.Error() == "OTP code not found or expired time" || err.Error() == "Incorrect OTP code" {
// 			statusCode = http.StatusUnauthorized
// 			message = err.Error()
// 		}

// 		h.log.Error("error in UserLoginByPhoneConfirm: " + err.Error())
// 		c.JSON(statusCode, models.Response{
// 			StatusCode:  statusCode,
// 			Description: message,
// 		})
// 		return
// 	}

// 	h.log.Info("Successfully logged in by phone")
// 	c.JSON(http.StatusOK, resp)
// }
