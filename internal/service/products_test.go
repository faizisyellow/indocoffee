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
				Name:     "lizzy.jpg",
				Size:     int64(len(imageFile)),
				Content:  imageFile,
				MimeType: "jpeg",
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
				MimeType: "txt",
			}

			if err := sut.Create(ctx, request, fileRequest); err == nil {
				t.Error("expected getting error but got nil")
				return
			}
		})
	})

}
