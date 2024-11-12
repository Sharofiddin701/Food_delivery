package handler

import (
	"food/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @ID 			create_payment
// @Router 		/food/api/v1/createpayment [POST]
// @Summary 	Create payment
// @Description Create a new payment
// @Tags 		payment
// @Accept 		json
// @Produce 	json
// @Param 		payment body models.CreatePayment true "Payment"
// @Success 	200 {object} Response{data=string} "Success"
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) CreatePayment(c *gin.Context) {
	var payment models.Payment

	if err := c.ShouldBindJSON(&payment); err != nil {
		h.log.Error(err.Error() + " : " + "error payment Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	resp, err := h.storage.Payment().Create(c.Request.Context(), &payment)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error payment Create")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("payment created successfully!")
	c.JSON(http.StatusCreated, Response{Data: resp})
}
