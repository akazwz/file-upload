package file

import (
	"github.com/akazwz/file-upload/api/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadFile(c *gin.Context) {
	fileUp := request.UploadFile{}

	err := c.ShouldBind(&fileUp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	fileHeader := fileUp.File
	contentType := fileHeader.Header.Get("Content-Type")

	filename := fileHeader.Filename
	size := fileHeader.Size

	c.JSON(http.StatusCreated, gin.H{
		"content_type": contentType,
		"filename":     filename,
		"size":         size,
	})
}
