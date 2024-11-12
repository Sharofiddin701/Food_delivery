package handler

import (
	"context"
	"food/api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @ID 			create_banner
// @Router 		/food/api/v1/createbanner [POST]
// @Summary 	Create Banner
// @Description Create a new banner
// @Tags 		banner
// @Accept 		json
// @Produce 	json
// @Param 		Banner body models.CreateBanner true "Banner"
// @Success 	200 {object} models.Banner
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) CreateBanner(c *gin.Context) {
	var banner models.Banner

	if err := c.ShouldBindJSON(&banner); err != nil {
		h.log.Error(err.Error() + " : " + "error Banner Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	resp, err := h.storage.Banner().Create(c.Request.Context(), &banner)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error Banner Create")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("Banner created successfully!")
	c.JSON(http.StatusCreated, resp)
}

// @ID 			get_all_banners
// @Router 		/food/api/v1/getallbanners [GET]
// @Summary 	Get All Banners
// @Description Retrieve all banners
// @Tags 		banner
// @Accept 		json
// @Produce 	json
// @Param 		search query string false "Search banners by image_url"
// @Param 		page   query uint64 false "Page number"
// @Param 		limit  query uint64 false "Limit number of results per page"
// @Success 	200 {object} models.GetAllBannerResponse
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) GetAllBanners(c *gin.Context) {
	var req = &models.GetAllBannerRequest{}

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

	banners, err := h.storage.Banner().GetAll(context.Background(), req)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error while getting all banners")
		c.JSON(http.StatusInternalServerError, "Error while getting all banners")
		return
	}

	h.log.Info("Banners retrieved successfully")
	c.JSON(http.StatusOK, banners)
}

// @ID 			delete_banner
// @Router 		/food/api/v1/deletebanner/{id} [DELETE]
// @Summary 	Delete Banner by ID
// @Description Delete a banner by its ID
// @Tags 		banner
// @Accept 		json
// @Produce 	json
// @Param 		id path string true "Banner ID"
// @Success 	200 {object} Response{data=string} "Success Request"
// @Response 	400 {object} Response{data=string} "Bad Request"
// @Failure 	500 {object} Response{data=string} "Server error"
func (h *Handler) DeleteBanner(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.log.Error("missing banner id")
		c.JSON(http.StatusBadRequest, "fill the gap with id")
		return
	}

	err := uuid.Validate(id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while validating id")
		c.JSON(http.StatusBadRequest, "please enter a valid id")
		return
	}

	err = h.storage.Banner().Delete(context.Background(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while deleting banner")
		c.JSON(http.StatusBadRequest, "please input valid data")
		return
	}

	h.log.Info("Banner deleted successfully!")
	c.JSON(http.StatusOK, id)
}
