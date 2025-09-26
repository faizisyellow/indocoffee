package main

import (
	"errors"
	"io"
	"net/http"

	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/faizisyellow/indocoffee/internal/uploader"
)

// @Summary		Add new product
// @Description	Create new coffee  product
// @Tags			Products
// @Accept			mpfd
// @Produce		json
//
// @Param			metadata	formData	string	true	"CreateVillaProp JSON string"	example({"roasted":"light","price":10.2,"quantity":,50,"bean":1,"form":1})
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
