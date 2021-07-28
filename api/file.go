package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// UploadChunkFile
// received chunks and file md5, create a dir named md5, then put chunks in it
// 接受上传的文件分块和文件MD5和分块MD5,把上传的区块存进命名为文件MD5的文件夹内
func UploadChunkFile(c *gin.Context) {
	fileMD5 := c.PostForm("file-md5")
	chunkMD5 := c.PostForm("chunk-md5")
	chunksCount := c.PostForm("chunks-count")
	chunksIndex := c.PostForm("chunks-index")
	chunks, _ := strconv.Atoi(chunksCount)

	// 接收文件
	chunkFile, err := c.FormFile("chunk-file")
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
	if exists {
		c.JSON(http.StatusOK, gin.H{
			"code":    "2000",
			"message": "success",
		})
		return
	}

	// 保存文件
	err = c.SaveUploadedFile(chunkFile, dst)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "4000",
			"message": "file save failed",
		})
		return
	}

	// 按照文件名排序读取文件夹内的文件
	files, _ := ioutil.ReadDir(dir)
	// 最后一个分块
	if len(files) == chunks {
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
			bytes, err := ioutil.ReadFile(dir + "/" + file.Name())
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "4000",
					"message": "read file failed",
				})
			}
			_, err = completeFile.Write(bytes)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "4000",
					"message": "write file failed",
				})
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    "2000",
			"message": "success",
			"url":     "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     "2000",
		"message":  "success",
		"progress": len(files) / chunks,
	})

	return
}

// MergeChunks
// when upload all chunks ,just call this api to merge all chunks to complete file named md5
// 当上传完所有区块时调用, 把所有的分块文件合并为文件名为MD5的文件
func MergeChunks(c *gin.Context) {
	md5 := c.PostForm("file-md5")
	dst := fmt.Sprintf("public/file/%s/%s", md5, md5)
	dir := fmt.Sprintf("public/file/%s", md5)
	exists, _ := PathExists(dst)
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
		bytes, err := ioutil.ReadFile(dir + "/" + file.Name())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "4000",
				"message": "read file failed",
			})
		}
		_, err = completeFile.Write(bytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "4000",
				"message": "write file failed",
			})
		}
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
