package request

import "mime/multipart"

type UploadFile struct {
	File *multipart.FileHeader `json:"file" form:"file" binding:"required"`
}

type UploadChunkFile struct {
	ChunkFile  *multipart.FileHeader `json:"chunk_file" form:"chunk_file" binding:"required"`
	ChunkIndex string                `json:"chunk_index" form:"chunk_index" binding:"required"`
	ChunkHash  string                `json:"chunk_hash" form:"chunk_hash" binding:"required"`
	ChunkSum   string                `json:"chunk_sum" form:"chunk_sum" binding:"required"`
	FileHash   string                `json:"file_hash" form:"file_hash" binding:"required"`
}
