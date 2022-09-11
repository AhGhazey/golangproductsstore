package postgres

import (
	"cmd/ims.server/pkg/adding"
	"cmd/ims.server/pkg/config"
	"cmd/ims.server/pkg/listing"
	postgres "cmd/ims.server/pkg/storage/postgres/productstore"
	"cmd/ims.server/pkg/updating"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository struct {
	postgres.InventoryStore
}

func NewRepository() (*Repository, error) {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load configurations")
		return nil, fmt.Errorf("Unable to load configurations %w", err)
	}

	db, err := sqlx.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Failed to open database")
		return nil, fmt.Errorf("error opening the database: %w", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("error connecting to the database")
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return &Repository{InventoryStore: *postgres.NewInventoryStore(db)}, nil
}

func (r *Repository) AddProduct(p adding.Product) error {
	err := r.InventoryStore.AddProduct(p)
	if err != nil {
		return fmt.Errorf("failed to add product:%w", err)
	}
	return nil
}

func (r *Repository) GetProductBySku(sku string) (*[]listing.Product, error) {
	products, err := r.InventoryStore.GetProductBySku(sku)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by sky:%w", err)
	}
	return products, nil
}

func (r *Repository) UpdateBulkRecords(products []updating.ProductCSV) error {
	err := r.InventoryStore.UpdateBulkRecords(products)
	if err != nil {
		return fmt.Errorf("failed to bulk update products:%w", err)
	}
	return nil
}
