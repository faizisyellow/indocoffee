package uploadthing

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/uploader"
	"github.com/joho/godotenv"
)

func TestUploadThing(t *testing.T) {
	t.Run("upload a file", func(t *testing.T) {
		upt, err := SetupUploadthing(t)
		if err != nil {
			t.Error(err)
			return
		}

		imageFile, err := os.ReadFile("../../service/file_test/lizzy.jpeg")
		if err != nil {
			t.Errorf("failed to read test file: %v", err)
			return
		}

		imageRequest := uploader.FileInput{
			Name:     "lizzymcalpine.jpg",
			Size:     int64(len(imageFile)),
			Content:  imageFile,
			MimeType: "image/jpeg",
		}

		url, err := upt.UploadFile(context.Background(), imageRequest)
		if err != nil {
			t.Error(err)
			return
		}

		if url == "" {
			t.Error("exptected got filename after upload but got empty")
		}
	})

	t.Run("upload a file that already exists", func(t *testing.T) {
		upt, err := SetupUploadthing(t)
		if err != nil {
			t.Error(err)
			return
		}

		imageFile, err := os.ReadFile("../../service/file_test/lizzy.jpeg")
		if err != nil {
			t.Errorf("failed to read test file: %v", err)
			return
		}

		imageRequest := uploader.FileInput{
			Name:     "lizzymcalpine.jpg",
			Size:     int64(len(imageFile)),
			Content:  imageFile,
			MimeType: "image/jpeg",
		}

		url, err := upt.UploadFile(context.Background(), imageRequest)
		if err != nil {
			t.Error(err)
			return
		}

		if url == "" {
			t.Error("exptected got filename after upload but got empty")
		}
	})

	t.Run("delete a file", func(t *testing.T) {
		upt, err := SetupUploadthing(t)
		if err != nil {
			t.Error(err)
			return
		}

		imageFile, err := os.ReadFile("../../service/file_test/lizzy.jpeg")
		if err != nil {
			t.Errorf("failed to read test file: %v", err)
			return
		}

		imageRequest := uploader.FileInput{
			Name:     "lizzymcalpine.jpg",
			Size:     int64(len(imageFile)),
			Content:  imageFile,
			MimeType: "image/jpeg",
		}

		url, err := upt.UploadFile(context.Background(), imageRequest)
		if err != nil {
			t.Error(err)
			return
		}

		if url == "" {
			t.Error("exptected got filename after upload but got empty")
		}

		if err := upt.DeleteFile(context.Background(), url); err != nil {
			t.Error(err)
			return
		}

	})

	t.Run("generate presigned an url", func(t *testing.T) {
		upt, err := SetupUploadthing(t)
		if err != nil {
			t.Error(err)
			return
		}

		req := struct {
			filename string
			filetype string
			filesize int
		}{
			filename: "lizzy.jpeg",
			filetype: "image/jpeg",
			filesize: 1468006,
		}

		response, err := upt.Register(req.filename, req.filetype, req.filesize)
		if err != nil {
			t.Error(err)
			return
		}

		if response.FileKey == "" {
			t.Error("should be got file key but empty")
			return
		}
	})

	t.Run("generate presigned an url and check the pool should be still working", func(t *testing.T) {
		upt, err := SetupUploadthing(t)
		if err != nil {
			t.Error(err)
			return
		}

		req := struct {
			filename string
			filetype string
			filesize int
		}{
			filename: "lizzy.jpeg",
			filetype: "image/jpeg",
			filesize: 1468006,
		}

		response, err := upt.Register(req.filename, req.filetype, req.filesize)
		if err != nil {
			t.Error(err)
			return
		}

		resPool, err := upt.PoolUpload(response.FileKey)
		if err != nil {
			t.Error(err)
			return
		}

		if resPool.Status != "still working" {
			t.Errorf("the upload image should be still working but got: %s", resPool.Status)
			return
		}
	})
}

func SetupUploadthing(t *testing.T) (*Uploadthing, error) {
	t.Helper()

	err := loadEnv()
	if err != nil {
		return nil, err
	}

	upt := New(
		os.Getenv("UPLOADTHING_API_KEY"),
		os.Getenv("UPLOADTHING_PRESIGNED_URL"),
		os.Getenv("UPLOADTHING_POOL_UPLOAD_URL"),
		"public-read",
		"imageUploader",
		os.Getenv("UPLOADTHING_UPLOAD_BY"),
		os.Getenv("UPLOADTHING_META_URL"),
		os.Getenv("UPLOADTHING_CALLBACK_URL"),
		os.Getenv("UPLOADTHING_DELETE_URL"),
		os.Getenv("UPLOADTHING_APP_ID"),
	)

	return upt, nil
}

func loadEnv() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return godotenv.Load(envPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return nil
}
