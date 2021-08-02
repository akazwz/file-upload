package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// UploadChunkFile
// received chunks and file md5, create a dir named md5, then put chunks in it
// 接受上传的文件分块和文件MD5和分块MD5,把上传的区块存进命名为文件MD5的文件夹内
func UploadChunkFile(c *gin.Context) {
	fileMD5 := c.PostForm("file-md5")
	chunkMD5 := c.PostForm("chunk-md5")
	chunksCount := c.PostForm("chunks-count")
	chunksIndex := c.PostForm("chunk-index")
	chunks, _ := strconv.Atoi(chunksCount)

	// 判断文件是否已经存在,存在直接秒传
	completeFile := fmt.Sprintf("public/file/%v/%v", fileMD5, fileMD5)
	isExist, err := PathExists(completeFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "4000",
			"message": "check file existed failed",
		})
		return
	}

	if isExist {
		c.JSON(http.StatusOK, gin.H{
			"code":     "2000",
			"message":  "success",
			"progress": 100,
			"url":      "",
			"a-pass":   1,
		})
		return
	}

	// 接收文件
	chunkFile, err := c.FormFile("chunk-file")
	// form接收文件错误
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "4000",
			"message": "form no file",
		})
		return
	}

	// 文件夹路径
	dir := fmt.Sprintf("public/file/%s", fileMD5)
	// 完整的文件路径
	dst := fmt.Sprintf("public/file/%s/%s", fileMD5, chunksIndex+"-"+chunkMD5)

	// 判断文件夹是否存在,不存在创建文件夹
	pathExists, _ := PathExists(dir)
	if !pathExists {
		err = os.Mkdir(dir, os.ModePerm)
		// 创建文件夹错误
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "4000",
				"message": "create dir failed",
			})
			return
		}
	}

	// 判断分块文件是否已经存在,已经存在直接返回成功
	exists, _ := PathExists(dst)
	// 用于已经获取的文件数
	filesInfo, _ := ioutil.ReadDir(dir)
	countFiles := len(filesInfo)
	if exists {
		c.JSON(http.StatusOK, gin.H{
			"code":     "2000",
			"message":  "success",
			"progress": 100 * (float64(countFiles) / float64(chunks)),
			"count":    countFiles,
		})
		return
	}

	// 保存文件
	err = c.SaveUploadedFile(chunkFile, dst)
	// 保存文件错误
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "4000",
			"message": "file save failed",
		})
		return
	}

	// 按照文件名排序读取文件夹内的文件, 默认是按文件名排序
	files, _ := ioutil.ReadDir(dir)
	filesCount := len(files)
	// 最后一个分块上传成功,开始合并
	if filesCount == chunks {
		// 创建文件
		completeFile, _ := os.Create(dir + "/" + fileMD5)
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
			// 读取每个分块文件字节
			bytes, err := ioutil.ReadFile(dir + "/" + file.Name())
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "4000",
					"message": "read file failed",
				})
				return
			}
			// 分块文件写入创建的文件
			_, err = completeFile.Write(bytes)
			// 分块文件写入失败
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "4000",
					"message": "write file failed",
				})
				return
			}
		}

		// 最后一个分块上传并合并文件成功,进度100 返回文件url
		c.JSON(http.StatusOK, gin.H{
			"code":     "2000",
			"message":  "success",
			"progress": 100,
			"url":      "",
		})
		return
	}

	// 分块文件上传成功
	c.JSON(http.StatusOK, gin.H{
		"code":     "2000",
		"message":  "success",
		"progress": 100 * (float64(filesCount) / float64(chunks)),
		"count":    filesCount,
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
