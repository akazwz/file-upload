package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// UploadChunkFile
// received chunks and file md5, create a dir named md5, then put chunks in it
// 接受上传的文件分块和文件MD5和分块MD5,把上传的区块存进命名为文件MD5的文件夹内
func UploadChunkFile(c *gin.Context) {
	md5 := c.PostForm("md5")
	chunkMd5 := c.PostForm("chunk-md5")
	file, err := c.FormFile("chunk-file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "4000",
			"message": "form no file",
		})
		return
	}

	dst := fmt.Sprintf("public/file/%s/%s", md5, chunkMd5)
	exists, _ := PathExists(dst)
	if exists {
		c.JSON(http.StatusOK, gin.H{
			"code":    "2000",
			"message": "success",
		})
		return
	}
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "4000",
			"message": "file save failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    "2000",
		"message": "success",
	})
	return
}

// MergeChunk
// when upload all chunks ,just call this api to merge all chunks to complete file named md5
// 当上传完所有区块时调用, 把所有的分块文件合并为文件名为MD5的文件
func MergeChunk(c *gin.Context) {
	md5 := c.PostForm("md5")
	dir := fmt.Sprintf("public/file/%s", md5)
	exists, _ := PathExists(dir)
	if exists {
		c.JSON(http.StatusOK, gin.H{
			"code":    "2000",
			"message": "success",
			"url":     "",
		})
		return
	}

	files, _ := ioutil.ReadDir(dir)
	completeFile, _ := os.Create(dir + "/" + md5)

	defer func(completeFile *os.File) {
		err := completeFile.Close()
		if err != nil {
			panic(err)
		}
	}(completeFile)

	for _, file := range files {
		if file.Name() == ".BS_Store" {
			continue
		}
		bytes, _ := ioutil.ReadFile(dir + "/" + file.Name())
		_, _ = completeFile.Write(bytes)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "2000",
		"message": "success",
		"url":     "",
	})
	return
}

// ChunksState
// check chunks state
// 查询文件分块上传状态
func ChunksState(c *gin.Context) {
	md5 := c.Query("md5")
	if md5 == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "4000",
			"message": "no such md5",
		})
		return
	}
	dir := fmt.Sprintf("public/file/%s", md5)

	var chunkList []string

	exists, _ := PathExists(dir)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "4000",
			"message": "no such md5",
		})
		return
	}

	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		fileName := file.Name()
		chunkList = append(chunkList, fileName)
		fileBaseName := strings.Split(fileName, ".")[0]
		if fileBaseName == md5 {
			c.JSON(http.StatusOK, gin.H{
				"code":    "2000",
				"message": "success",
				"state":   1,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":       "2000",
		"message":    "success",
		"state":      0,
		"chunk-list": chunkList,
	})
	return
}

// PathExists
// check file or dir exists
// 检查文件或者文件夹是否存在
func PathExists(dst string) (bool, error) {
	_, err := os.Stat(dst)
	if err != nil {
		return false, err
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
