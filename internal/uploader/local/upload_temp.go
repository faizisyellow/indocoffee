package local

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/faizisyellow/indocoffee/internal/uploader"
)

type TempUpload struct {
	LocSavePath string
}

func (t *TempUpload) UploadFile(ctx context.Context, file uploader.FileInput) (string, error) {

	saveDir := t.LocSavePath
	if saveDir == "" {
		saveDir = "."
	}

	if err := os.MkdirAll(saveDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(saveDir, file.Name)

	// Write the file
	if err := os.WriteFile(filePath, file.Content, 0o644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

func (t *TempUpload) DeleteFile(ctx context.Context, filename string) error {
	filePath := filepath.Join(t.LocSavePath, filename)

	if err := os.Remove(filePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// no-op if file doesn't exist
			return nil
		}
		return fmt.Errorf("failed to delete file %s: %w", filename, err)
	}

	return nil
}
