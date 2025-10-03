package main

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/faizisyellow/indocoffee/internal/uploader"
	"github.com/go-chi/chi/v5"
)

// @Summary		Add new product
// @Description	Create new coffee  product
// @Tags			Products
// @Accept			mpfd
// @Produce		json
//
// @Param			metadata	formData	string	true	"Create product JSON string"	example({"roasted":"light","price":10.2,"quantity":50,"bean":1,"form":1})
//
// @Param			file		formData	file	true	"Image file"
// @Success		201			{object}	main.Envelope{data=string,error=nil}
// @Failure		400			{object}	main.Envelope{data=nil,error=string}
// @Failure		409			{object}	main.Envelope{data=nil,error=string}
// @Failure		500			{object}	main.Envelope{data=nil,error=string}
// @Router			/products [post]
func (app *Application) CreateProductsHandler(w http.ResponseWriter, r *http.Request) {

	// Limit request body to 3 MB
	r.Body = http.MaxBytesReader(w, r.Body, 3<<20+1024)

	if err := r.ParseMultipartForm(3 << 20); err != nil {
		ResponseClientError(w, r, errors.New("file too big"), http.StatusBadRequest)
		return
	}

	var request dto.CreateProductMetadataRequest
	if err := ReadJsonMultiPartForm(r, "metadata", &request); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(request); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	if handler.Size > 2<<20 {
		ResponseClientError(w, r, errors.New("file too big"), http.StatusBadRequest)
		return
	}

	limitedReader := io.LimitReader(file, 2<<20)
	fileBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusInternalServerError)
		return
	}

	uploaded := uploader.FileInput{
		Name:     handler.Filename,
		Size:     int64(len(fileBytes)),
		MimeType: http.DetectContentType(fileBytes),
		Content:  fileBytes,
	}

	if err := app.Services.ProductsService.Create(r.Context(), request, uploaded); err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrFileNotSupportedProducts:
			ResponseClientError(w, r, err, http.StatusBadRequest)
		case service.ErrFileTooBigProducts:
			ResponseClientError(w, r, err, http.StatusBadRequest)
		case service.ErrReferenceFailedProducts:
			ResponseClientError(w, r, err, http.StatusBadRequest)
		case service.ErrConflictProducts:
			ResponseClientError(w, r, err, http.StatusConflict)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, "success create new product", http.StatusCreated)
}

// @Summary		Get product
// @Description	Get coffee  product by id
// @Tags			Products
// @Accept			json
// @Produce		json
//
// @Param			id	path		int	true	"product id"
// @Success		200	{object}	main.Envelope{data=dto.GetProductResponse,error=nil}
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/products/{id} [get]
func (app *Application) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	product, err := app.Services.ProductsService.FindById(r.Context(), id)
	if err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrNotFoundProduct:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	response := dto.GetProductResponse{
		Id:       product.Id,
		Roasted:  product.Roasted,
		Price:    product.Price,
		Quantity: product.Quantity,
		Image:    product.Image,
		BeanId:   product.BeanId,
		FormId:   product.FormId,
	}
	response.Bean.Name = product.BeansModel.Name
	response.Form.Name = product.FormsModel.Name

	ResponseSuccess(w, r, response, http.StatusOK)
}

// @Summary		Get products
// @Description	Get all coffee  products
// @Tags			Products
// @Accept			json
// @Produce		json
// @Param			limit	query		string	false	"limit each page"
// @Param			offset	query		string	false	"skip rows"
// @Param			sort	query		string	false	"sort product by alphabet"
// @Param			roast	query		string	false	"roasted coffee"
// @Param			form	query		string	false	"what kind of form of the coffee (form id)"
// @Param			bean	query		string	false	"what kind of bean of the coffee (bean id)"
// @Success		200		{object}	main.Envelope{data=[]dto.GetProductsResponse,error=nil}
// @Success		400		{object}	main.Envelope{data=nil,error=string}
// @Failure		500		{object}	main.Envelope{data=nil,error=string}
// @Router			/products [get]
func (app *Application) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	queryValue := r.URL.Query()
	query := repository.QueryProducts{
		Limit:  queryValue.Get("limit"),
		Offset: queryValue.Get("offset"),
		Sort:   queryValue.Get("sort"),
		Roast:  queryValue.Get("roast"),
		Form:   queryValue.Get("form"),
		Bean:   queryValue.Get("bean"),
	}

	paginateProductQuery, err := repository.PaginatedProductsQuery{Limit: 8, Sort: "asc"}.Parse(query)
	if err != nil {
		ResponseClientError(w, r, errorService.New(errors.New("invalid query value"), err), http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(paginateProductQuery); err != nil {
		ResponseClientError(w, r, errorService.New(err, err), http.StatusBadRequest)
		return
	}

	products, err := app.Services.ProductsService.FindProducts(r.Context(), paginateProductQuery)
	if err != nil {
		ResponseServerError(w, r, err, http.StatusInternalServerError)
		return
	}

	var response []dto.GetProductsResponse
	for _, product := range products {
		res := dto.GetProductsResponse{
			Id:       product.Id,
			Roasted:  product.Roasted,
			Price:    product.Price,
			Quantity: product.Quantity,
			Image:    product.Image,
			BeanId:   product.BeanId,
			FormId:   product.FormId,
		}
		res.Bean.Name = product.BeansModel.Name
		res.Form.Name = product.FormsModel.Name

		response = append(response, res)
	}

	ResponseSuccess(w, r, response, http.StatusOK)
}

// @Summary		Edit product
// @Description	Update spesific product
// @Tags			Products
// @Accept			mpfd
// @Produce		json
//
// @Param			metadata	formData	string	true	"Update product JSON string"	example({"roasted":"medium","quantity":50})
//
// @Param			file		formData	file	false	"Image file"
// @Param			id			path		int		true	"Product id"
// @Success		200			{object}	main.Envelope{data=string,error=nil}
// @Failure		400			{object}	main.Envelope{data=nil,error=string}
// @Failure		404			{object}	main.Envelope{data=nil,error=string}
// @Failure		409			{object}	main.Envelope{data=nil,error=string}
// @Failure		500			{object}	main.Envelope{data=nil,error=string}
// @Router			/products/{id} [patch]
func (app *Application) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Limit request body to 3 MB
	r.Body = http.MaxBytesReader(w, r.Body, 3<<20+1024)

	if err := r.ParseMultipartForm(3 << 20); err != nil {
		ResponseClientError(w, r, errors.New("file too big"), http.StatusBadRequest)
		return
	}

	var request dto.UpdateProductMetadataRequest
	if err := ReadJsonMultiPartForm(r, "metadata", &request); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := Validate.Struct(request); err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	file, handler, err := GetOptionalFormFile(r, "file")
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	var uploadInput uploader.FileInput

	if file != nil && handler != nil {
		defer file.Close()

		if handler.Size > 2<<20 {
			ResponseClientError(w, r, errors.New("file too big"), http.StatusBadRequest)
			return
		}

		limitedReader := io.LimitReader(file, 2<<20)
		fileBytes, err := io.ReadAll(limitedReader)
		if err != nil {
			ResponseClientError(w, r, err, http.StatusInternalServerError)
			return
		}

		uploadInput.Name = handler.Filename
		uploadInput.Size = int64(len(fileBytes))
		uploadInput.MimeType = http.DetectContentType(fileBytes)
		uploadInput.Content = fileBytes
	}

	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := app.Services.ProductsService.Update(ctx, id, request, uploadInput); err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrNotFoundProduct:
			ResponseClientError(w, r, err, http.StatusNotFound)
		case service.ErrConflictProducts:
			ResponseClientError(w, r, err, http.StatusConflict)
		case service.ErrReferenceFailedProducts:
			ResponseClientError(w, r, err, http.StatusBadRequest)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, "success update product", http.StatusOK)
}

// @Summary		Delete product
// @Description	Delete product by id
// @Tags			Products
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"Product id"
// @Success		204
// @Failure		400	{object}	main.Envelope{data=nil,error=string}
// @Failure		404	{object}	main.Envelope{data=nil,error=string}
// @Failure		500	{object}	main.Envelope{data=nil,error=string}
// @Router			/products/{id} [delete]
func (app *Application) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ResponseClientError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := app.Services.ProductsService.Destroy(r.Context(), id); err != nil {
		errValue := errorService.GetError(err)
		switch errValue.E {
		case service.ErrNotFoundProduct:
			ResponseClientError(w, r, err, http.StatusNotFound)
		default:
			ResponseServerError(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	ResponseSuccess(w, r, nil, http.StatusNoContent)
}
