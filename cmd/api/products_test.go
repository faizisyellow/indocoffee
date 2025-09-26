package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/service/dto"
)

func TestProducts(t *testing.T) {
	t.Run("create new product", func(t *testing.T) {

		var (
			app     = setupTestApplication(t)
			handler = app.Mux()
			request = dto.CreateProductMetadataRequest{
				Roasted:  "light",
				Price:    17.5,
				Quantity: 10,
				Bean:     1,
				Form:     1,
			}
		)

		buf := &bytes.Buffer{}
		writer := multipart.NewWriter(buf)

		jsonData, err := json.Marshal(request)
		if err != nil {
			t.Fatal(err)
		}
		jsonPart, err := writer.CreateFormField("metadata")
		if err != nil {
			t.Fatal(err)
		}
		_, err = jsonPart.Write(jsonData)
		if err != nil {
			t.Fatal(err)
		}

		// Add the file field
		filePart, err := writer.CreateFormFile("file", "test.txt")
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(filePart, strings.NewReader("this is a test file"))
		if err != nil {
			t.Fatal(err)
		}

		// Close writer to set the terminating boundary
		writer.Close()

		req, err := http.NewRequest("POST", "/v1/products", buf)

		req.Header.Set("Content-Type", writer.FormDataContentType())

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != 201 {
			t.Errorf("should be success, got: %v", rr.Code)
		}
	})
}
