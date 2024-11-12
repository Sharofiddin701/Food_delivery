package handler

import (
	"context"
	"net/http"
	"strconv"

	"food/api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @ID 			 create_branch
// @Router       /food/api/v1/createbranch [POST]
// @Summary      Create Branch
// @Description  Create Branch
// @Tags         branch
// @Accept       json
// @Category     json
// @Param        Branch body models.CreateBranch true "Branch"
// @Success      201 {object} models.Branch
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) CreateBranch(c *gin.Context) {
	branch := &models.Branch{}

	if err := c.ShouldBindJSON(&branch); err != nil {
		h.log.Error(err.Error() + " : " + "error Branch Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	resp, err := h.storage.Branch().Create(c.Request.Context(), branch)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error Branch Create")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("Branch created successfully!")
	c.JSON(http.StatusCreated, resp)
}

// @ID 			 update_branch
// @Router       /food/api/v1/updatebranch/{id} [PUT]
// @Summary      Update Branch
// @Description  Update Branch
// @Tags         branch
// @Accept       json
// @Category     json
// @Param        id path string true "Branch ID"
// @Param        Branch body models.UpdateBranch true "UpdateBranchRequest"
// @Success      200 {object} models.Branch
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) UpdateBranch(c *gin.Context) {
	var updateBranch models.UpdateBranch

	if err := c.ShouldBindJSON(&updateBranch); err != nil {
		h.log.Error(err.Error() + " : " + "error Branch Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, enter valid data!")
		return
	}

	id := c.Param("id")
	branch, err := h.storage.Branch().GetByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error Branch Not Found")
		c.JSON(http.StatusBadRequest, "Branch not found!")
		return
	}

	branch.Name = updateBranch.Name
	branch.Address = updateBranch.Address
	branch.Latitude = updateBranch.Latitude
	branch.Longitude = updateBranch.Longitude

	resp, err := h.storage.Branch().Update(c.Request.Context(), branch)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error Branch Update")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("Branch updated successfully!")
	c.JSON(http.StatusOK, resp)
}

// @ID 			 get_branch
// @Router		 /food/api/v1/getbranch/{id} [GET]
// @Summary		 Get Branch by ID
// @Description  Get a branch by its ID and return its info
// @Tags		 branch
// @Accept		 json
// @Produce		 json
// @Param		 id path string true "Branch ID"
// @Success		 200  {object}  models.Branch
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) GetBranchByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.log.Error("missing branch id")
		c.JSON(http.StatusBadRequest, "you must fill the id")
		return
	}

	branch, err := h.storage.Branch().GetByID(context.Background(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "Error while getting branch by id")
		c.JSON(http.StatusInternalServerError, "Server Error")
		return
	}

	h.log.Info("Branch was successfully gotten by ID")
	c.JSON(http.StatusOK, branch)
}

// @ID 			    getall_branch
// @Router 			/food/api/v1/getallbranches [GET]
// @Summary 		Get all branches
// @Description	Retrieves information about all branches
// @Tags 			branch
// @Accept 			json
// @Produce 		json
// @Param 			search query string false "Search branches"
// @Param 			page query uint64 false "Page number"
// @Param 			limit query uint64 false "Limit per page"
// @Success 		200 {object} models.GetAllBranchesResponse
// @Response     400 {object} Response{data=string} "Bad Request"
// @Failure      500 {object} Response{data=string} "Server error"
func (h *Handler) GetAllBranches(c *gin.Context) {
	var (
		req = &models.GetAllBranchesRequest{}
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

	branches, err := h.storage.Branch().GetAll(context.Background(), req)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while getting all branches")
		c.JSON(http.StatusInternalServerError, "error while getting all branches")
		return
	}

	h.log.Info("Branches were successfully retrieved")
	c.JSON(http.StatusOK, branches)
}

// @ID 			delete_branch
// @Router		/food/api/v1/deletebranch/{id} [DELETE]
// @Summary		Delete a branch by its ID
// @Description This API deletes a branch by its ID
// @Tags		branch
// @Accept		json
// @Produce		json
// @Param		id path string true "Branch ID"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *Handler) DeleteBranch(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.log.Error("missing branch id")
		c.JSON(http.StatusBadRequest, "fill the gap with id")
		return
	}

	err := uuid.Validate(id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while validating id")
		c.JSON(http.StatusBadRequest, "please enter a valid id")
		return
	}

	err = h.storage.Branch().Delete(context.Background(), id)
	if err != nil {
		h.log.Error(err.Error() + ":" + "error while deleting branch")
		c.JSON(http.StatusInternalServerError, "Server error!")
		return
	}

	h.log.Info("Branch deleted successfully!")
	c.JSON(http.StatusOK, id)
}
