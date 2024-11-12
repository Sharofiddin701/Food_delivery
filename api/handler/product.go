package handler

import (
	"context"
	"food/api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @ID 			create_product
// @Router 		/food/api/v1/createproduct [POST]
// @Summary 	Create Product
// @Description Create a new product
// @Tags 		product
// @Accept 		json
// @Produce 	json
// @Param 		Product body models.CreateProduct true "Product"
// @Success 	200 {object} models.Product
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		h.log.Error(err.Error() + " : " + "error Product Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	resp, err := h.storage.Product().Create(c.Request.Context(), &product)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error Product Create")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("Product created successfully!")
	c.JSON(http.StatusCreated, resp)
}

// @ID 			update_product
// @Router 		/food/api/v1/updateproduct/{id} [PUT]
// @Summary 	Update Product
// @Description Update an existing product
// @Tags 		product
// @Accept 		json
// @Produce 	json
// @Param 		id path string true "Product ID"
// @Param 		Product body models.UpdateProduct true "UpdateProductRequest"
// @Success 	200 {object} models.Product
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) UpdateProduct(c *gin.Context) {
	var updateProduct models.UpdateProduct

	if err := c.ShouldBindJSON(&updateProduct); err != nil {
		h.log.Error(err.Error() + " : " + "error Product Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	id := c.Param("id")
	product, err := h.storage.Product().GetByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error Product Not Found")
		c.JSON(http.StatusBadRequest, "Product not found!")
		return
	}

	product.Name = updateProduct.Name
	product.Description = updateProduct.Description
	product.Price = updateProduct.Price
	product.ImageURL = updateProduct.ImageURL
	product.CategoryId = updateProduct.CategoryId

	resp, err := h.storage.Product().Update(c.Request.Context(), product)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error Product Update")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("Product updated successfully!")
	c.JSON(http.StatusOK, resp)
}

// @ID 			get_product
// @Router 		/food/api/v1/getproduct/{id} [GET]
// @Summary 	Get Product by ID
// @Description Retrieve a product by its ID
// @Tags 		product
// @Accept 		json
// @Produce 	json
// @Param 		id path string true "Product ID"
// @Success 	200 {object} models.Product
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.log.Error("missing product id")
		c.JSON(http.StatusBadRequest, "you must fill the ID")
		return
	}

	product, err := h.storage.Product().GetByID(context.Background(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error while getting product by ID")
		c.JSON(http.StatusInternalServerError, "Server Error")
		return
	}

	h.log.Info("Product retrieved successfully by ID")
	c.JSON(http.StatusOK, product)
}

// @ID 			get_all_products
// @Router 		/food/api/v1/getallproducts [GET]
// @Summary 	Get All Products
// @Description Retrieve all products
// @Tags 		product
// @Accept 		json
// @Produce 	json
// @Param       category_id query string false "get by category_id"
// @Param 		search query string false "Search products by name or description"
// @Param 		page   query uint64 false "Page number"
// @Param 		limit  query uint64 false "Limit number of results per page"
// @Success 	200 {object} models.GetAllProductsResponse
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) GetAllProducts(c *gin.Context) {
	var req = &models.GetAllProductsRequest{}

	req.Search = c.Query("search")
	req.CategoryId = c.Query("category_id")

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

	products, err := h.storage.Product().GetAll(context.Background(), req)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error while getting all products")
		c.JSON(http.StatusInternalServerError, "Error while getting all products")
		return
	}

	h.log.Info("Products retrieved successfully")
	c.JSON(http.StatusOK, products)
}

// @ID 			delete_product
// @Router 		/food/api/v1/deleteproduct/{id} [DELETE]
// @Summary 	Delete Product by ID
// @Description Delete a product by its ID
// @Tags 		product
// @Accept 		json
// @Produce 	json
// @Param 		id path string true "Product ID"
// @Success 	200 {object} Response{data=string} "Success Request"
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.log.Error("missing product id")
		c.JSON(http.StatusBadRequest, "fill the gap with id")
		return
	}

	err := uuid.Validate(id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while validating id")
		c.JSON(http.StatusBadRequest, "please enter a valid id")
		return
	}

	err = h.storage.Product().Delete(context.Background(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while deleting product")
		c.JSON(http.StatusBadRequest, "please input valid data")
		return
	}

	h.log.Info("Product deleted successfully!")
	c.JSON(http.StatusOK, id)
}
