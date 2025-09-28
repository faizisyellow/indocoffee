package uploader

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
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

// CreateMultipartBody builds a multipart/form-data body with a single file field
// Returns (bodyBuffer, contentTypeHeader, error)
func CreateMultipartBody(file FileInput) (*bytes.Buffer, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the file field
	filePart, err := writer.CreateFormFile("file", file.Name)
	if err != nil {
		return nil, "", err
	}

	_, err = io.Copy(filePart, bytes.NewBuffer(file.Content))
	if err != nil {
		return nil, "", err
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return &buf, writer.FormDataContentType(), nil
}
