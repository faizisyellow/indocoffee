package products_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/db"
	"github.com/faizisyellow/indocoffee/internal/repository/products"
)

func TestProductsWithRealDB(t *testing.T) {
	if getEnvironment(t) != "development" {
		t.Skip("skipping test: only runs in development environment")
	}

	products.Contract{func() (products.Products, func()) {
		newDB, err := setupTestDB(t)
		if err != nil {
			t.Fatal(err)
		}
		return &products.ProductRepository{newDB}, func() {
			if err := newDB.Close(); err != nil {
				t.Fatal(err)
			}
		}
	}}.Test(t)
}

func setupTestDB(t *testing.T) (*sql.DB, error) {
	t.Helper()

	return db.New(
		os.Getenv("DB_TEST_ADDR"),
		5,
		5,
		"1m",
		"1m",
	)
}

func getEnvironment(t *testing.T) string {
	t.Helper()

	return os.Getenv("ENV")
}
