package file

import (
	"fmt"
	"github.com/akazwz/file-upload/api/request"
	"github.com/akazwz/file-upload/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func UploadChunk(c *gin.Context) {
	chunkFileUp := request.UploadChunkFile{}

	err := c.ShouldBind(&chunkFileUp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}

	fileHeader := chunkFileUp.ChunkFile

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

	if chunkFileUp.ChunkHash != sha256Hash {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "文件不完整",
		})
		return
	}

	// 文件夹路径
	dir := fmt.Sprintf("public/file/%s", chunkFileUp.FileHash)
	// 完整的文件路径
	dst := fmt.Sprintf("public/file/%s/%s", chunkFileUp.FileHash, chunkFileUp.ChunkIndex+"-"+chunkFileUp.ChunkHash)

	// 判断文件夹是否存在,不存在创建文件夹
	pathExists, _ := utils.PathExists(dir)
	if !pathExists {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "创建文件夹失败",
			})
			return
		}
	}

	// 判断分块文件是否已经存在,已经存在直接返回成功
	exists, _ := PathExists(dst)
	if exists {
		c.JSON(http.StatusCreated, gin.H{
			"message": "已经有此分块文件",
		})
		return
	}

	// 保存文件
	err = c.SaveUploadedFile(chunkFileUp.ChunkFile, dst)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "文件保存失败",
		})
		return
	}

	// 按照文件名排序读取文件夹内的文件
	files, _ := ioutil.ReadDir(dir)

	chunkSum, _ := strconv.Atoi(chunkFileUp.ChunkSum)

	// 最后一个分块
	if len(files) == chunkSum {
		log.Println("finish")
		// 创建文件
		completeFile, _ := os.Create(dir + "/" + "complete_file")

		defer func(completeFile *os.File) {
			err := completeFile.Close()
			if err != nil {
				panic(err)
			}
		}(completeFile)

		// 写入文件
		for _, file := range files {
			if file.Name() == ".BS_Store" {
				continue
			}
			bytes, err := ioutil.ReadFile(dir + "/" + file.Name())
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "读取文件失败",
				})
			}
			_, err = completeFile.Write(bytes)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "写入文件失败",
				})
			}
		}

		c.JSON(http.StatusCreated, gin.H{
			"content_type": contentType,
			"filename":     filename,
			"size":         size,
			"chunk_index":  chunkFileUp.ChunkIndex,
			"chunk_sum":    chunkSum,
			"hash_sha256":  sha256Hash,
		})
	}
}
