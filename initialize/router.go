package initialize

import (
	"net/http"

	"github.com/akazwz/file-upload/api/file"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
	}))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
	})

	//Teapot  418
	r.GET("teapot", func(c *gin.Context) {
		c.JSON(http.StatusTeapot, gin.H{
			"message": "I'm a teapot",
			"story": "This code was defined in 1998 " +
				"as one of the traditional IETF April Fools' jokes," +
				" in RFC 2324, Hyper Text Coffee Pot Control Protocol," +
				" and is not expected to be implemented by actual HTTP servers." +
				" However, known implementations do exist.",
		})
	})

	fileGroup := r.Group("/file")
	{
		fileGroup.POST("/", file.UploadFile)

		fileGroup.POST("/chunks", file.UploadChunkFile)
		fileGroup.POST("/chunks/merge", file.MergeChunks)
		fileGroup.GET("/chunks/state", file.ChunksState)
	}
	return r
}
