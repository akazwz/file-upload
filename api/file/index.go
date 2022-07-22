package file

import (
	"fmt"
	"net/http"

	"github.com/akazwz/file-upload/api/request"
	"github.com/akazwz/file-upload/utils"
	"github.com/gin-gonic/gin"
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

	// 文件信息
	contentType := fileHeader.Header.Get("Content-Type")
	filename := fileHeader.Filename
	size := fileHeader.Size

	// 获取 sha256 hash
	sha256Hash, err := utils.HashFileByAlgo(fileHeader, "sha256")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "读取文件失败",
		})
		return
	}

	// 文件夹路径
	dst := fmt.Sprintf("public/file/%s", filename)
	// 保存文件
	err = c.SaveUploadedFile(fileHeader, dst)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "保存文件失败",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"content_type": contentType,
		"filename":     filename,
		"size":         size,
		"hash_sha256":  sha256Hash,
	})
}
