package main

import (
	"testing"

	"github.com/faizisyellow/indocoffee/internal/logger"
	"github.com/faizisyellow/indocoffee/internal/repository/products"
	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/uploader/local"
)

func setupTestApplication(t *testing.T) *Application {
	t.Helper()

	return &Application{
		Logger: logger.Logger,
		Services: *service.New(
			nil,
			nil,
			nil,
			nil,
			nil,
			&products.InMemoryProducts{},
			&local.TempUpload{},
			nil,
			nil,
			nil,
		),
	}
}
