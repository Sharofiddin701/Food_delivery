package handler

import (
	"context"
	"fmt"

	"food/api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @ID 			 create_category
// @Router       /food/api/v1/category [POST]
// @Summary      Create Category
// @Description  Create Category
// @Tags         category
// @Accept       json
// @Category     json
// @Param        Category body models.CreateCategory true "Category"
// @Success      200 {object} models.Category
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) CreateCategory(c *gin.Context) {
	var categoryCreate models.CreateCategory

	if err := c.ShouldBindJSON(&categoryCreate); err != nil {
		h.log.Error(err.Error() + " : " + "error Category Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	category := &models.Category{
		Name: categoryCreate.Name,
	}
	resp, err := h.storage.Category().Create(c.Request.Context(), category)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error Category Create")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("Category created successfully!")
	c.JSON(http.StatusCreated, resp)
}

// @ID 			 update_category
// @Router       /food/api/v1/category/{id} [PUT]
// @Summary      Update Category
// @Description  Update Category
// @Tags         category
// @Accept       json
// @Category     json
// @Param        id path string true "Category ID"
// @Param        Category body models.UpdateCategory true "UpdateCategoryRequest"
// @Success      200 {object} models.Category
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) UpdateCategory(c *gin.Context) {
	var updateCategory models.UpdateCategory

	if err := c.ShouldBindJSON(&updateCategory); err != nil {
		h.log.Error(err.Error() + " : " + "error Category Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	id := c.Param("id")
	category, err := h.storage.Category().GetByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error Category Not Found")
		c.JSON(http.StatusBadRequest, "Category not found!")
		return
	}

	category.Name = updateCategory.Name
	resp, err := h.storage.Category().Update(c.Request.Context(), category)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error Category Update")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("Category updated successfully!")
	c.JSON(http.StatusOK, resp)
}

// @ID 			 get_category
// @Router		 /food/api/v1/getbycategory/{id} [GET]
// @Summary		 get a category by its id
// @Description  This api gets a category by its id and returns its info
// @Tags		 category
// @Accept		 json
// @Produce		 json
// @Param		 id path string true "id"
// @Success		 200  {object}  models.Category
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.log.Error("missing category id")
		c.JSON(http.StatusBadRequest, "you must fill")
		return
	}

	category, err := h.storage.Category().GetByID(context.Background(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error while getting by id category")
		c.JSON(http.StatusInternalServerError, "Server Error")
		return
	}

	h.log.Info("Category was successfully gotten by Id")
	c.JSON(http.StatusOK, category)
}

// @ID 			    getall_category
// @Router 			/food/api/v1/getallcategory [GET]
// @Summary 		Get all category
// @Description		Retrieves information about all categories
// @Tags 			category
// @Accept 			json
// @Produce 		json
// @Param 			search query string false "categories"
// @Param 			page query uint64 false "page"
// @Param 			limit query uint64 false "limit"
// @Success 		200 {object} models.GetAllCategoriesResponse
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) GetAllCategories(c *gin.Context) {
	var (
		req = &models.GetAllCategoriesRequest{}
	)

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
		c.JSON(http.StatusInternalServerError, "internal server error while parsing limit")
		return
	}

	req.Page = page
	req.Limit = limit

	customers, err := h.storage.Category().GetAll(context.Background(), req)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while getting all categories")
		c.JSON(http.StatusInternalServerError, "error while getting all categories")
		return
	}

	h.log.Info("Category was successfully gotten by Id")
	c.JSON(http.StatusOK, customers)
}

// @ID 			delete_category
// @Router		/food/api/v1/deletecategory/{id} [DELETE]
// @Summary		delete a category by its id
// @Description This api deletes a category by its id
// @Tags		category
// @Accept		json
// @Produce		json
// @Param		id path string true "category ID"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h Handler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("id: ", id)

	if id == "" {
		h.log.Error("missing category id")
		c.JSON(http.StatusBadRequest, "fill the gap with id")
		return
	}

	err := uuid.Validate(id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while validating id")
		c.JSON(http.StatusBadRequest, "please enter valid id")
		return
	}

	err = h.storage.Category().Delete(context.Background(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while deleting customer")
		c.JSON(http.StatusBadRequest, "please input valid data")
		return
	}

	h.log.Info("Category deleted succesfully!")
	c.JSON(http.StatusOK, id)
}
