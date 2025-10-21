package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/repository/products"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/faizisyellow/indocoffee/internal/uploader"
	"github.com/faizisyellow/indocoffee/internal/uploader/uploadthing"
)

type ProductsService struct {
	ProductsStore products.Products
	Uploader      uploader.Uploader
}

const FILE_SUPPORTED_MAIN = "png"
const FILE_SUPPORTED_SECOND = "jpeg"

var (
	ErrFileTooBigProducts       = errors.New("products: request upload image too big")
	ErrUploadImageProducts      = errors.New("products: error uploading image file")
	ErrFileNotSupportedProducts = errors.New("products: file not supported. Only png and jpeg files are supported")
	ErrInternalProducts         = errors.New("products: internal error")
	ErrConflictProducts         = errors.New("products: already exist")
	ErrReferenceFailedProducts  = errors.New("products: form or beans not found")
	ErrNotFoundProduct          = errors.New("products: product not found")
)

func (p *ProductsService) Create(ctx context.Context, metadatReq dto.CreateProductMetadataRequest, file uploader.FileInput) error {
	// // max upload file 2mb
	if file.Size > 2<<20 {
		return errorService.New(ErrFileTooBigProducts, ErrFileTooBigProducts)
	}

	if !strings.Contains(file.MimeType, FILE_SUPPORTED_MAIN) && !strings.Contains(file.MimeType, FILE_SUPPORTED_SECOND) {
		return errorService.New(ErrFileNotSupportedProducts, ErrFileNotSupportedProducts)
	}

	filename, err := p.Uploader.UploadFile(ctx, file)
	if err != nil {
		return errorService.New(ErrUploadImageProducts, err)
	}

	newProduct := models.Product{
		Roasted:  metadatReq.Roasted,
		Price:    metadatReq.Price,
		Quantity: metadatReq.Quantity,
		BeanId:   metadatReq.Bean,
		FormId:   metadatReq.Form,
		Image:    filename,
	}

	if err := p.ProductsStore.Insert(ctx, newProduct); err != nil {
		if err := p.Uploader.DeleteFile(ctx, filename); err != nil {
			return errorService.New(ErrInternalProducts, err)
		}

		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return errorService.New(ErrConflictProducts, err)
		}

		if strings.Contains(err.Error(), REFERENCES_CODE) {
			return errorService.New(ErrReferenceFailedProducts, err)
		}

		return errorService.New(err, err)
	}

	return nil

}

func (p *ProductsService) FindById(ctx context.Context, id int) (models.Product, error) {

	product, err := p.ProductsStore.GetById(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.Product{}, errorService.New(ErrNotFoundProduct, err)
		default:
			return models.Product{}, err
		}
	}

	return product, nil
}

func (p *ProductsService) FindProducts(ctx context.Context, r repository.PaginatedProductsQuery) ([]models.Product, error) {

	products, err := p.ProductsStore.GetAll(ctx, r)
	if err != nil {
		return nil, errorService.New(ErrInternalProducts, err)
	}

	return products, nil
}

func (p *ProductsService) Update(ctx context.Context, id int, req dto.UpdateProductMetadataRequest, file uploader.FileInput) error {

	product, err := p.FindById(ctx, id)
	if err != nil {
		return err
	}

	if req.Roasted != "" {
		product.Roasted = req.Roasted
	}

	if req.Quantity != 0 {
		product.Quantity = req.Quantity
	}

	if req.Price != 0 {
		product.Price = req.Price
	}

	if req.Form != 0 {
		product.FormId = req.Form
	}

	if req.Bean != 0 {
		product.BeanId = req.Bean
	}

	var existingImage = product.Image

	if len(file.Content) != 0 {
		if file.Size > 2<<20 {
			return errorService.New(ErrFileTooBigProducts, ErrFileTooBigProducts)
		}

		if !strings.Contains(file.MimeType, FILE_SUPPORTED_MAIN) && !strings.Contains(file.MimeType, FILE_SUPPORTED_SECOND) {
			return errorService.New(ErrFileNotSupportedProducts, ErrFileNotSupportedProducts)
		}

		filename, err := p.Uploader.UploadFile(ctx, file)
		if err != nil {
			return errorService.New(ErrUploadImageProducts, err)
		}

		product.Image = filename
	}

	if err := p.ProductsStore.Update(ctx, product); err != nil {
		if len(file.Content) != 0 {
			if err := p.Uploader.DeleteFile(ctx, uploadthing.GetFileKey(product.Image)); err != nil {
				log.Printf("error delete image in error update product: %v", err.Error())
			}
		}

		if strings.Contains(err.Error(), CONFLICT_CODE) {
			return errorService.New(ErrConflictProducts, err)
		}

		if strings.Contains(err.Error(), REFERENCES_CODE) {
			return errorService.New(ErrReferenceFailedProducts, err)
		}

		return errorService.New(ErrInternalProducts, err)
	}

	// if there's a file request, delete previous
	// image in file storage provider
	if len(file.Content) != 0 {
		err := p.Uploader.DeleteFile(ctx, uploadthing.GetFileKey(existingImage))
		if err != nil {
			return errorService.New(ErrInternalProducts, err)
		}
	}

	return nil
}

func (p *ProductsService) DecreaseQuantityProduct(ctx context.Context, tx *sql.Tx, prdId, quantity int) error {

	return p.ProductsStore.DecrementQuantity(ctx, tx, prdId, quantity)
}

func (p *ProductsService) IncreaseQuantityProduct(ctx context.Context, tx *sql.Tx, prdId, quantity int) error {

	return p.ProductsStore.IncrementQuantity(ctx, tx, prdId, quantity)
}

func (p *ProductsService) Destroy(ctx context.Context, id int) error {
	product, err := p.FindById(ctx, id)
	if err != nil {
		return err
	}

	err = p.ProductsStore.Delete(ctx, product.Id)
	if err != nil {
		return errorService.New(ErrInternalProducts, err)
	}

	err = p.Uploader.DeleteFile(ctx, uploadthing.GetFileKey(product.Image))
	if err != nil {
		return errorService.New(ErrInternalProducts, err)
	}

	return nil
}
