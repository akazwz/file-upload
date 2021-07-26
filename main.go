package main

import (
	"file-upload/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Handle("POST", "/chunk-file", api.UploadChunkFile)
	r.Handle("POST", "/merge-chunk", api.MergeChunk)
	r.Handle("GET", "/chunks-state", api.ChunksState)
}
