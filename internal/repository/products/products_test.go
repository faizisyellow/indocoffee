package products_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/db"
	"github.com/faizisyellow/indocoffee/internal/repository/products"
	"github.com/joho/godotenv"
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

	err := loadEnv()
	if err != nil {
		return nil, err
	}

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
	err := loadEnv()
	if err != nil {
		t.Fatal(err)
	}

	return os.Getenv("ENV")
}

func loadEnv() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return godotenv.Load(envPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return nil
}
