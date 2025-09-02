package uploader

import (
	"context"
	"io"
)

type Uploader interface {
	UploadFile(ctx context.Context, key string, file io.ReadCloser, contentType string, errUpload chan<- error)
}
