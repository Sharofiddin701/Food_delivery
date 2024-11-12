package handler

import (
	"food/pkg/helper"

	"github.com/gin-gonic/gin"
)

// UploadFiles godoc
// @ID upload_multiple_files
// @Router /food/api/v1/uploadfiles [post]
// @Summary Upload Multiple Files
// @Description Upload Multiple Files
// @Tags Upload File
// @Accept multipart/form-data
// @Produce json
// @Param file formData []file true "File to upload" 
// @Success 200 {object} Response{data=string} "Success Request"
// @Failure 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *Handler) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		h.log.Error(err.Error() + "  :  " + "File error")
		c.JSON(400, Response{Data: "File error"})
		return
	}

	resp, err := helper.UploadFiles(form)
	if err != nil {
		h.log.Error(err.Error() + "  :  " + "Upload error")
		c.JSON(500, Response{Data: "Upload error"})
		return
	}

	c.JSON(200, resp)
}

// DeleteFile godoc
// @ID delete_file
// @Router /food/api/v1/deletefile [delete]
// @Summary Delete File
// @Description Delete File
// @Tags Upload File
// @Accept json
// @Produce json
// @Param id query string true "ID of the file to delete"
// @Success 204 {string} string "Success Request"
// @Failure 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *Handler) DeleteFile(c *gin.Context) {
	err := helper.DeleteFile(c.Query("id"))
	if err != nil {
		h.log.Error(err.Error() + "  :  " + "Upload error")
		c.JSON(500, Response{Data: "Server error"})
		return
	}
	c.JSON(204, "success")
}
