package handler

import (
	"context"
	"encoding/json"
	"food/api/models"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Create Order godoc
// @ID          create_order
// @Router      /food/api/v1/order [POST]
// @Summary     Create Order
// @Description Create Order
// @Tags        order
// @Accept      json
// @Order       json
// @Param       Order body models.SwaggerOrderCreateRequest true "CreateOrderRequest"
// @Success     201 {object} Response{data=string} "Success Request"
// @Response    400 {object} Response{data=string} "Bad Request"
// @Failure     500 {object} Response{data=string} "Server error"
func (h *Handler) CreateOrder(c *gin.Context) {
	var (
		request models.OrderCreateRequest
	)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("error reading body: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Server Error!"})
		return
	}
	h.log.Info("Incoming JSON: " + string(body))

	err = json.Unmarshal(body, &request)
	if err != nil {
		h.log.Error("error unmarshalling JSON: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{Data: "Invalid JSON!"})
		return
	}

	if request.Order.UserId == "" {
		h.log.Error("Customer ID is empty!")
		c.JSON(http.StatusBadRequest, Response{Data: "Customer ID is required!"})
		return
	}
	for _, item := range request.Items {
		if item.ProductId == "" {
			h.log.Error("Product ID is empty for one of the items!")
			c.JSON(http.StatusBadRequest, Response{Data: "Product ID is required for each item!"})
			return
		}
	}

	order, err := h.storage.Order().Create(c.Request.Context(), &request)
	if err != nil {
		h.log.Error("error in Order.CreateOrder: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Server Error!"})
		return
	}

	h.log.Info("Order Created Successfully!")
	c.JSON(http.StatusCreated, Response{Data: order})
}

// @ID             get_order
// @Router         /food/api/v1/getorder/{id} [GET]
// @Summary        get_order
// @Description    get_order
// @Tags           order
// @Accept         json
// @Produces 	   json
// @Param          id path string true "Order Id or user Id"
// @Success 	   200 {object} models.OrderCreateRequest
// @Response 	   400 {object} Response{data=string} "Bad Request"
// @Failure 	   500 {object} Response{data=string} "Server error"
func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.log.Error("missing order id")
		c.JSON(http.StatusBadRequest, Response{Data: "you must fill the id"})
		return
	}

	order, err := h.storage.Order().GetOrder(c.Request.Context(), id)
	if err != nil {
		h.log.Error("error in Order.GetByID: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Server Error!"})
		return
	}

	h.log.Info("Order retrieved successfully!")
	c.JSON(http.StatusOK, Response{Data: order})
}

// @ID 			   get_all_orders
// @Router 		   /food/api/v1/getallorders [GET]
// @Summary 	   Get All Orders
// @Description    Retrieve all orders
// @Tags 		   order
// @Accept 		   json
// @Produce 	   json
// @Param 		   search query string false "Search orders by name or description"
// @Param 		   page   query uint64 false "Page number"
// @Param 		   limit  query uint64 false "Limit number of results per page"
// @Success 	   200 {object} Response{data=string} "Success"
// @Response 	   400 {object} Response{data=string} "Bad Request"
// @Failure 	   500 {object} Response{data=string} "Server error"
func (h *Handler) GetAllOrders(c *gin.Context) {
	var req = &models.GetAllOrdersRequest{}

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

	products, err := h.storage.Order().GetAll(context.Background(), req)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error while getting all products")
		c.JSON(http.StatusInternalServerError, "Error while getting all products")
		return
	}

	h.log.Info("Products retrieved successfully")
	c.JSON(http.StatusOK, Response{Data: products})
}

// @ID 			update_order
// @Router 		/food/api/v1/updateorder/{id} [PUT]
// @Summary 	Update Order
// @Description Update an existing order
// @Tags 		order
// @Accept 		json
// @Produce 	json
// @Param 		id path string true "Order ID"
// @Param 		Order body models.OrderUpdateS true "UpdateOrderRequest"
// @Success 	200 {object} Response{data=string} "Success"
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) UpdateOrder(c *gin.Context) {
	var updateOrder *models.Order

	if err := c.ShouldBindJSON(&updateOrder); err != nil {
		h.log.Error(err.Error() + " : " + "error Order Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	id := c.Param("id")

	resp, err := h.storage.Order().Update(c.Request.Context(), id, updateOrder)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error Order Update")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("Order updated successfully!")
	c.JSON(http.StatusOK, Response{Data: resp})
}

// @ID 			delete_order
// @Router 		/food/api/v1/deleteorder/{id} [DELETE]
// @Summary 	Delete order by Id
// @Description Delete a order by its Id
// @Tags 		order
// @Accept 		json
// @Produce 	json
// @Param 		id path string true "Order Id"
// @Success 	200 {object} Response{data=string} "Success Request"
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) DeleteOrder(c *gin.Context) {
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

	err = h.storage.Order().Delete(context.Background(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while deleting order")
		c.JSON(http.StatusBadRequest, "please input valid data")
		return
	}

	h.log.Info("Order deleted successfully!")
	c.JSON(http.StatusOK, Response{Data: id})
}

// @ID             change_order_status
// @Router         /food/api/v1/orderStatus/{id} [PATCH]
// @Summary        Change Order Status
// @Description    Change the status of an order
// @Tags           order
// @Accept         json
// @Produces       json
// @Param          id path string true "Order ID"
// @Param          status body models.PatchOrderStatusRequest true "New Order Status"
// @Success        200 {object} Response{data=string} "Order status updated successfully"
// @Response       400 {object} Response{data=string} "Bad Request"
// @Response       500 {object} Response{data=string} "Server error"
func (h *Handler) ChangeOrderStatus(c *gin.Context) {
	// Get order ID from the URL
	orderId := c.Param("id")

	// Initialize the request body struct
	req := models.PatchOrderStatusRequest{}

	// Check if order ID is provided
	if orderId == "" {
		h.log.Error("missing order id")
		c.JSON(http.StatusBadRequest, Response{Data: "Order ID is required"})
		return
	}

	// Bind the JSON request body to the req variable
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to bind request body: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{Data: "Invalid request body"})
		return
	}

	// Call the repository method to change the order status
	resp, err := h.storage.Order().ChangeOrderStatus(c.Request.Context(), &req, orderId)
	if err != nil {
		h.log.Error("failed to update order status: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Server Error!"})
		return
	}

	h.log.Info("Order status updated successfully")
	c.JSON(http.StatusOK, Response{Data: resp})
}
