package service_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/repository/products"
	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	"github.com/faizisyellow/indocoffee/internal/uploader"
	"github.com/faizisyellow/indocoffee/internal/uploader/local"
	"github.com/faizisyellow/indocoffee/internal/uploader/uploadthing"
	"github.com/joho/godotenv"
)

func TestProductsService(t *testing.T) {
	t.Run("run in memory store and in local upload file", func(t *testing.T) {
		ProductsServiceTest{
			CreateDependencies: func() (products.Products, uploader.Uploader, Cleanup) {
				savePath := "file_test/temp/"

				return &products.InMemoryProducts{}, &local.TempUpload{
						LocSavePath: savePath,
					}, func() {
						entries, err := os.ReadDir(savePath)
						if err != nil {
							return
						}

						for _, e := range entries {
							if !e.IsDir() {
								_ = os.Remove(filepath.Join(savePath, e.Name()))
							}
						}
					}
			},
		}.Test(t)
	})

	t.Run("run in memory store and in uploadthing service to upload file", func(t *testing.T) {
		ProductsServiceTest{
			CreateDependencies: func() (products.Products, uploader.Uploader, Cleanup) {
				upld, err := SetupUploadthing(t)
				if err != nil {
					t.Fatal(err)
				}
				return &products.InMemoryProducts{}, upld, func() {

				}
			},
		}.Test(t)
	})
}

type ProductsServiceTest struct {
	CreateDependencies func() (products.Products, uploader.Uploader, Cleanup)
}

func (p ProductsServiceTest) Test(t *testing.T) {
	t.Run("create new product", func(t *testing.T) {

		t.Run("success create new product", func(t *testing.T) {
			var (
				ctx                                 = context.Background()
				productsStore, uploadFile, teardown = p.CreateDependencies()
				sut                                 = service.ProductsService{productsStore, uploadFile}
				request                             = dto.CreateProductMetadataRequest{
					Roasted:  "light",
					Price:    18.5,
					Quantity: 100,
					Bean:     1,
					Form:     1,
				}
			)
			t.Cleanup(teardown)

			imageFile, err := os.ReadFile("file_test/lizzy.jpeg")
			if err != nil {
				t.Errorf("failed to read test file: %v", err)
				return
			}

			imageRequest := uploader.FileInput{
				Name:     "lizzy.jpeg",
				Size:     int64(len(imageFile)),
				Content:  imageFile,
				MimeType: "image/jpeg",
			}

			if err := sut.Create(ctx, request, imageRequest); err != nil {
				t.Errorf("expected not getting error but got %v", err.Error())
				return
			}
		})

		t.Run("failed create new product because the file not suppported", func(t *testing.T) {
			var (
				ctx                                 = context.Background()
				productsStore, uploadFile, teardown = p.CreateDependencies()
				sut                                 = service.ProductsService{productsStore, uploadFile}
				request                             = dto.CreateProductMetadataRequest{
					Roasted:  "light",
					Price:    18.5,
					Quantity: 100,
					Bean:     1,
					Form:     1,
				}
			)
			t.Cleanup(teardown)

			textFile, err := os.ReadFile("file_test/greetings.txt")
			if err != nil {
				t.Errorf("failed to read test file: %v", err)
				return
			}

			fileRequest := uploader.FileInput{
				Name:     "greetings.txt",
				Size:     int64(len(textFile)),
				Content:  textFile,
				MimeType: "text/plain",
			}

			if err := sut.Create(ctx, request, fileRequest); err == nil {
				t.Error("expected getting error but got nil")
				return
			}
		})
	})

}

func SetupUploadthing(t *testing.T) (*uploadthing.Uploadthing, error) {
	t.Helper()

	err := loadEnv()
	if err != nil {
		return nil, err
	}

	upt := uploadthing.New(
		os.Getenv("UPLOADTHING_API_KEY"),
		os.Getenv("UPLOADTHING_PRESIGNED_URL"),
		os.Getenv("UPLOADTHING_POOL_UPLOAD_URL"),
		"public-read",
		"imageUploader",
		os.Getenv("UPLOADTHING_UPLOAD_BY"),
		os.Getenv("UPLOADTHING_META_URL"),
		os.Getenv("UPLOADTHING_CALLBACK_URL"),
		os.Getenv("UPLOADTHING_DELETE_URL"),
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
