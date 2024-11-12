package handler

import (
	"context"
	_ "food/api/docs"
	"food/api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @ID 			 create_admin
// @Router       /food/api/v1/createadmin [POST]
// @Summary      Create Admin
// @Description  Create a new admin
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        Admin body models.CreateAdmin true "Admin"
// @Success      200 {object} models.Admin
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) CreateAdmin(c *gin.Context) {
	var (
		admin = models.Admin{}
	)

	if err := c.ShouldBindJSON(&admin); err != nil {
		h.log.Error(err.Error() + " : " + "error User Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	resp, err := h.storage.Admin().Create(c.Request.Context(), &admin)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error Admin Create")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("User created successfully!")
	c.JSON(http.StatusCreated, resp)
}

// @ID 			 update_admin
// @Router       /food/api/v1/updateadmin/{id} [PUT]
// @Summary      Update Admin
// @Description  Update an existing Admin
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id path string true "Admin ID"
// @Param        User body models.UpdateAdmin true "UpdateAdminRequest"
// @Success      200 {object} models.Admin
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) UpdateAdmin(c *gin.Context) {
	var updateAdmin models.UpdateAdmin

	if err := c.ShouldBindJSON(&updateAdmin); err != nil {
		h.log.Error(err.Error() + " : " + "error User Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	id := c.Param("id")
	user, err := h.storage.Admin().GetByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error User Not Found")
		c.JSON(http.StatusBadRequest, "User not found!")
		return
	}

	user.Name = updateAdmin.Name
	user.Email = updateAdmin.Email

	resp, err := h.storage.Admin().Update(c.Request.Context(), user)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error User Update")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("User updated successfully!")
	c.JSON(http.StatusOK, resp)
}

// @ID 			 get_admin
// @Router		 /food/api/v1/getbyidadmin/{id} [GET]
// @Summary		 Get User by ID
// @Description  Retrieve a user by their ID
// @Tags		 admin
// @Accept		 json
// @Produce		 json
// @Param		 id path string true "Admin ID"
// @Success		 200  {object}  models.Admin
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) GetAdminByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.log.Error("missing user id")
		c.JSON(http.StatusBadRequest, "you must fill")
		return
	}

	user, err := h.storage.Admin().GetByID(context.Background(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error while getting user by ID")
		c.JSON(http.StatusInternalServerError, "Server Error")
		return
	}

	h.log.Info("User retrieved successfully by ID")
	c.JSON(http.StatusOK, user)
}

// @ID 			    get_all_admins
// @Router 			/food/api/v1/getalladmins [GET]
// @Summary 		Get All Admins
// @Description		Retrieve all admins
// @Tags 			admin
// @Accept 			json
// @Produce 		json
// @Param 			search query string false "Search users by name or email"
// @Param 			page   query uint64 false "Page number"
// @Param 			limit  query uint64 false "Limit number of results per page"
// @Success 		200 {object} models.GetAllAdminsResponse
// @Response        400 {object} Response{data=string} "Bad Request"
// @Failure         500 {object} Response{data=string} "Server error"
func (h *Handler) GetAllAdmins(c *gin.Context) {
	var req = &models.GetAllAdminsRequest{}

	req.Search = c.Query("search")

	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while parsing page")
		c.JSON(http.StatusBadRequest, "BadRequest at paging")
		return
	}

	limit, err := strconv.ParseUint(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while parsing limit")
		c.JSON(http.StatusInternalServerError, "Internal server error while parsing limit")
		return
	}

	req.Page = page
	req.Limit = limit

	users, err := h.storage.Admin().GetAll(context.Background(), req)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error while getting all admins")
		c.JSON(http.StatusInternalServerError, "Error while getting all users")
		return
	}

	h.log.Info("Users retrieved successfully")
	c.JSON(http.StatusOK, users)
}

// @ID 			delete_admin
// @Router		/food/api/v1/deleteadmin/{id} [DELETE]
// @Summary		Delete Admin by ID
// @Description Delete a admin by their ID
// @Tags		admin
// @Accept		json
// @Produce		json
// @Param		id path string true "Admin ID"
// @Success     200 {object} Response{data=string} "Success Request"
// @Response    400 {object} Response{data=string} "Bad Request"
// @Failure     500 {object} Response{data=string} "Server error"
func (h *Handler) DeleteAdmin(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.log.Error("missing user id")
		c.JSON(http.StatusBadRequest, "fill the gap with id")
		return
	}

	err := uuid.Validate(id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while validating id")
		c.JSON(http.StatusBadRequest, "please enter a valid id")
		return
	}

	err = h.storage.User().Delete(context.Background(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while deleting user")
		c.JSON(http.StatusBadRequest, "please input valid data")
		return
	}

	h.log.Info("User deleted successfully!")
	c.JSON(http.StatusOK, id)
}
