package main

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func WriteHttpJson(w http.ResponseWriter, data any, status int) error {

	switch v := data.(type) {
	case Envelope:
		if v.Error != nil {
			v.Error = v.Error.(error).Error()
			data = v
		}
	}

	if status == http.StatusNoContent {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadHttpJson(w http.ResponseWriter, r *http.Request, data any) error {

	if len(r.Header["Content-Type"]) == 0 {
		return fmt.Errorf("no header content-type found")
	}

	if r.Header["Content-Type"][0] != "application/json" {
		return fmt.Errorf("header request only accept json")
	}

	// limit of the body size for 1mb
	maxBytes := 1_048_578

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	return decode.Decode(data)

}

// ReadJsonMultiPartForm reads first value of the field
func ReadJsonMultiPartForm(r *http.Request, field string, data any) error {

	if len(r.MultipartForm.Value[field]) == 0 {
		return fmt.Errorf("no fields are found")
	}

	jsonField := r.MultipartForm.Value[field][0]

	decoder := json.NewDecoder(strings.NewReader(jsonField))
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func GetOptionalFormFile(r *http.Request, field string) (multipart.File, *multipart.FileHeader, error) {
	file, handler, err := r.FormFile(field)
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, nil, nil // gracefully no file
		}
		return nil, nil, err
	}
	return file, handler, nil
}
