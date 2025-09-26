package uploader

import (
	"context"
)

type FileInput struct {
	Name     string
	Size     int64
	MimeType string
	Content  []byte
}

type Uploader interface {
	UploadFile(ctx context.Context, file FileInput) (string, error)
	DeleteFile(ctx context.Context, filename string) error
}
