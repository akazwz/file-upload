package main

import (
	"file-upload/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
	}))
	r.Handle("POST", "/chunk-file", api.UploadChunkFile)
	r.Handle("POST", "/merge-chunks", api.MergeChunks)
	r.Handle("GET", "/chunks-state", api.ChunksState)
	_ = r.Run(":8888")
}
