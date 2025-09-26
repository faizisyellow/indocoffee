package service

import (
	"context"
	"errors"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository/products"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/faizisyellow/indocoffee/internal/uploader"
)

type ProductsService struct {
	ProductsStore products.Products
	Uploader      uploader.Uploader
}

var (
	ErrFileTooBigProducts  = errors.New("products: request upload image too big")
	ErrUploadImageProducts = errors.New("products: error uploading image file")
)

func (p *ProductsService) Create(ctx context.Context, metadatReq dto.CreateProductMetadataRequest, file uploader.FileInput) error {

	// max upload file 2mb
	if file.Size > 2<<20<<20 {
		return errorService.New(ErrFileTooBigProducts, ErrFileTooBigProducts)
	}

	filename, err := p.Uploader.UploadFile(ctx, file)
	if err != nil {
		return errorService.New(ErrUploadImageProducts, err)
	}

	newProduct := models.Product{
		Roasted:  metadatReq.Roasted,
		Price:    metadatReq.Price,
		Quantity: int(metadatReq.Quantity),
		BeanId:   int(metadatReq.Bean),
		FormId:   int(metadatReq.Bean),
		Image:    filename,
	}

	return p.ProductsStore.Insert(ctx, newProduct)
}
