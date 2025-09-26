package service_test

import (
	"context"
	"os"
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
				return &products.InMemoryProducts{}, &local.TempUpload{}, func() {}
			},
		}.Test(t)
	})
}

type ProductsServiceTest struct {
	CreateDependencies func() (products.Products, uploader.Uploader, Cleanup)
}

func (p ProductsServiceTest) Test(t *testing.T) {
	t.Run("create new product", func(t *testing.T) {
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
			t.Fatalf("failed to read test file: %v", err)
		}

		imageRequest := uploader.FileInput{
			Name:    "lizzy.jpg",
			Size:    int64(len(imageFile)),
			Content: imageFile,
		}

		if err := sut.Create(ctx, request, imageRequest); err != nil {
			t.Errorf("expected not getting error but got %v", err.Error())
		}
	})

}
