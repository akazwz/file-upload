package request

import "mime/multipart"

type UploadFile struct {
	File *multipart.FileHeader `json:"file" form:"file" binding:"required"`
}
