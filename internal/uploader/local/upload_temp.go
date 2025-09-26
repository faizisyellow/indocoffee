package local

import (
	"context"

	"github.com/faizisyellow/indocoffee/internal/uploader"
)

type TempUpload struct {
}

func (t *TempUpload) UploadFile(ctx context.Context, file uploader.FileInput) (string, error) {

	return "", nil
}
