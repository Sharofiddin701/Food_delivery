package handler

import (
	"context"
	"encoding/json"
	"food/api/models"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Create Combo godoc
// @ID          create_combo
// @Router      /food/api/v1/combo [POST]
// @Summary     Create Combo
// @Description Create a new combo with a set of items
// @Tags        combo
// @Accept      json
// @Produce     json
// @Param       Combo body models.SwaggerComboCreateRequest true "CreateComboRequest"
// @Success     201 {object} Response{data=string} "Success Request"
// @Response    400 {object} Response{data=string} "Bad Request"
// @Failure     500 {object} Response{data=string} "Server error"
func (h *Handler) CreateCombo(c *gin.Context) {
	var (
		request models.ComboCreateRequest
	)

	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error("error reading body: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Server Error!"})
		return
	}
	h.log.Info("Incoming JSON: " + string(body))

	// Unmarshal the request body into the ComboCreateRequest struct
	err = json.Unmarshal(body, &request)
	if err != nil {
		h.log.Error("error unmarshalling JSON: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{Data: "Invalid JSON!"})
		return
	}

	// Validate the Combo data
	if request.Combo.Name == "" {
		h.log.Error("Combo name is empty!")
		c.JSON(http.StatusBadRequest, Response{Data: "Combo name is required!"})
		return
	}
	if request.Combo.Price <= 0 {
		h.log.Error("Invalid combo price!")
		c.JSON(http.StatusBadRequest, Response{Data: "Valid combo price is required!"})
		return
	}

	// Validate each item in the combo
	for _, item := range request.Combo.ComboItems {
		if item.ProductId == "" {
			h.log.Error("Product ID is empty for one of the items!")
			c.JSON(http.StatusBadRequest, Response{Data: "Product ID is required for each item!"})
			return
		}
		if item.Quantity <= 0 {
			h.log.Error("Invalid quantity for product: " + item.ProductId)
			c.JSON(http.StatusBadRequest, Response{Data: "Valid quantity is required for each item!"})
			return
		}
	}

	// Call the Create method in the repository to insert the combo into the database
	combo, err := h.storage.Combo().Create(c.Request.Context(), &request)
	if err != nil {
		h.log.Error("error in Combo.Create: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Server Error!"})
		return
	}

	// Respond with the created combo's ID
	h.log.Info("Combo Created Successfully!")
	c.JSON(http.StatusCreated, Response{Data: combo})
}

// @ID 			   get_all_combos
// @Router 		   /food/api/v1/getallcombos [GET]
// @Summary 	   Get All Combos
// @Description    Retrieve all Combo
// @Tags 		   combo
// @Accept 		   json
// @Produce 	   json
// @Param 		   search query string false "Search combos by name or description"
// @Param 		   page   query uint64 false "Page number"
// @Param 		   limit  query uint64 false "Limit number of results per page"
// @Success 	   200 {object} Response{data=string} "Success"
// @Response 	   400 {object} Response{data=string} "Bad Request"
// @Failure 	   500 {object} Response{data=string} "Server error"
func (h *Handler) GetAllCombos(c *gin.Context) {
	var req = &models.GetAllCombosRequest{}

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

	products, err := h.storage.Combo().GetAll(context.Background(), req)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error while getting all combos")
		c.JSON(http.StatusInternalServerError, "Error while getting all combos")
		return
	}

	h.log.Info("Combos retrieved successfully")
	c.JSON(http.StatusOK, Response{Data: products})
}

// @ID             get_combo
// @Router         /food/api/v1/getcombo/{id} [GET]
// @Summary        get_combo
// @Description    get_combo by its id
// @Tags           combo
// @Accept         json
// @Produces 	   json
// @Param          id path string true "Combo Id"
// @Success 	   200 {object}  Response{data=string} "Successfully retrieved combo"
// @Response 	   400 {object} Response{data=string} "Bad Request"
// @Failure 	   500 {object} Response{data=string} "Server error"
func (h *Handler) GetCombo(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.log.Error("missing combo id")
		c.JSON(http.StatusBadRequest, Response{Data: "you must fill the id"})
		return
	}

	combo, err := h.storage.Combo().GetCombo(c.Request.Context(), id)
	if err != nil {
		h.log.Error("error in Combo.GetByID: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Server Error!"})
		return
	}

	h.log.Info("Order retrieved successfully!")
	c.JSON(http.StatusOK, Response{Data: combo})
}

// @ID 			update_combo
// @Router 		/food/api/v1/updatecombo/{id} [PUT]
// @Summary 	Update Combo
// @Description Update an existing combo
// @Tags 		combo
// @Accept 		json
// @Produce 	json
// @Param 		id path string true "Combo ID"
// @Param 		Combo body models.ComboUpdateS true "UpdateOrderRequest"
// @Success 	200 {object} Response{data=string} "Successfully"
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) UpdateCombo(c *gin.Context) {
	var updateCombo *models.Combo

	if err := c.ShouldBindJSON(&updateCombo); err != nil {
		h.log.Error(err.Error() + " : " + "error Order Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	id := c.Param("id")

	resp, err := h.storage.Combo().Update(c.Request.Context(), id, updateCombo)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error Order Update")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("Order updated successfully!")
	c.JSON(http.StatusOK, Response{Data: resp})
}
