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

// @ID 			 create_user
// @Router       /food/api/v1/createuser [POST]
// @Summary      Create User
// @Description  Create a new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        User body models.CreateUser true "User"
// @Success      200 {object} models.User
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) CreateUser(c *gin.Context) {
	var (
		user = models.User{}
	)

	if err := c.ShouldBindJSON(&user); err != nil {
		h.log.Error(err.Error() + " : " + "error User Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	resp, err := h.storage.User().Create(c.Request.Context(), &user)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error User Create")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("User created successfully!")
	c.JSON(http.StatusCreated, resp)
}

// @ID 			 update_user
// @Router       /food/api/v1/updateuser/{id} [PUT]
// @Summary      Update User
// @Description  Update an existing user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        User body models.UpdateUser true "UpdateUserRequest"
// @Success      200 {object} models.User
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) UpdateUser(c *gin.Context) {
	var updateUser models.UpdateUser

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		h.log.Error(err.Error() + " : " + "error User Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	id := c.Param("id")
	user, err := h.storage.User().GetByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error User Not Found")
		c.JSON(http.StatusBadRequest, "User not found!")
		return
	}

	user.Name = updateUser.Name
	user.Email = updateUser.Email

	resp, err := h.storage.User().Update(c.Request.Context(), user)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error User Update")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("User updated successfully!")
	c.JSON(http.StatusOK, resp)
}

// @ID 			 get_user
// @Router		 /food/api/v1/getbyiduser/{id} [GET]
// @Summary		 Get User by ID
// @Description  Retrieve a user by their ID
// @Tags		 user
// @Accept		 json
// @Produce		 json
// @Param		 id path string true "User ID"
// @Success		 200  {object}  models.User
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.log.Error("missing user id")
		c.JSON(http.StatusBadRequest, "you must fill")
		return
	}

	user, err := h.storage.User().GetByID(context.Background(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error while getting user by ID")
		c.JSON(http.StatusInternalServerError, "Server Error")
		return
	}

	h.log.Info("User retrieved successfully by ID")
	c.JSON(http.StatusOK, user)
}

// @ID 			    get_all_users
// @Router 			/food/api/v1/getallusers [GET]
// @Summary 		Get All Users
// @Description		Retrieve all users
// @Tags 			user
// @Accept 			json
// @Produce 		json
// @Param 			search query string false "Search users by name or email"
// @Param 			page   query uint64 false "Page number"
// @Param 			limit  query uint64 false "Limit number of results per page"
// @Success 		200  {object} models.GetAllUsersResponse
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) GetAllUsers(c *gin.Context) {
	var req = &models.GetAllUsersRequest{}

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

	users, err := h.storage.User().GetAll(context.Background(), req)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error while getting all users")
		c.JSON(http.StatusInternalServerError, "Error while getting all users")
		return
	}

	h.log.Info("Users retrieved successfully")
	c.JSON(http.StatusOK, users)
}

// @ID 			delete_user
// @Router		/food/api/v1/deleteuser/{id} [DELETE]
// @Summary		Delete User by ID
// @Description Delete a user by their ID
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		id path string true "User ID"
// @Success     200 {object} Response{data=string} "Success Request"
// @Response    400 {object} Response{data=string} "Bad Request"
// @Failure     500 {object} Response{data=string} "Server error"
func (h *Handler) DeleteUser(c *gin.Context) {
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
