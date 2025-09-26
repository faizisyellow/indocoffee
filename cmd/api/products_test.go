package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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

		imageFile, err := os.ReadFile("../../testdata/file/lizzy.jpeg")
		if err != nil {
			t.Errorf("failed to read test file: %v", err)
			return
		}

		// Add the file field
		filePart, err := writer.CreateFormFile("file", "lizzy.jpeg")
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(filePart, bytes.NewBuffer(imageFile))
		if err != nil {
			t.Fatal(err)
		}

		// Close writer to set the terminating boundary
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/v1/products", buf)

		req.Header.Set("Content-Type", writer.FormDataContentType())

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != 201 {
			t.Errorf("should be success, got: %v", rr.Code)
		}
	})
}
