package postgres

import (
	"cmd/ims.server/pkg/adding"
	"cmd/ims.server/pkg/listing"
	"cmd/ims.server/pkg/updating"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type ProductStore interface {
	GetProductBySKU(sku string) (*[]listing.Product, error)
	ConsumeProduct(updating.Product) (bool, error)
	AddProduct(adding.Product) error
	UpdateBulkRecords(products []updating.ProductCSV) error
}

func NewInventoryStore(db *sqlx.DB) *InventoryStore {
	return &InventoryStore{
		DB: db,
	}
}

type InventoryStore struct {
	*sqlx.DB
}

func (istore *InventoryStore) GetProductBySku(sku string) (*[]listing.Product, error) {
	var products []listing.Product
	query := `SELECT sku, name, country,quantity FROM product WHERE sku= $1`
	err := istore.Select(&products, query, sku)
	if err != nil {
		return nil, fmt.Errorf("error getting product:%w", err)
	}
	return &products, nil
}

func (istore *InventoryStore) ConsumeProduct(p updating.Product) (bool, error) {
	tx := istore.MustBegin()
	var err error
	defer func() {
		if err != nil {
			log.Println("transaction rollback")
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	query := `SELECT id, sku, country, quantity, name FROM product WHERE sku= $1 AND country = $2 LIMIT 1`

	var product Product
	err = tx.Get(&product, query, p.Sku, p.Country)

	if err != nil {
		log.Println("unable to get product", err)
		return false, err
	}

	if product.Quantity < p.Quantity {
		return false, nil
	}

	newQuantity := product.Quantity - p.Quantity

	updateQuery := `UPDATE product SET quantity = $1 WHERE id = $2`

	_, err = tx.Exec(updateQuery, newQuantity, product.Id)
	if err != nil {
		log.Println("unable to update product", err)
	}

	return true, nil
}

func (istore *InventoryStore) AddProduct(p adding.Product) error {
	query := `INSERT INTO  product (sku, country, name, quantity)
				VALUES($1, $2, $3, $4) 
			 	ON CONFLICT ON CONSTRAINT sku_country_unique
				DO 
	   			UPDATE SET quantity = $5 + product.quantity;`

	_, err := istore.Exec(query, p.Sku, p.Country, p.Name, p.Quantity, p.Quantity)

	if err != nil {
		log.Println("unable to insert product", err)
		return err
	}

	return nil
}

func (istore *InventoryStore) UpdateBulkRecords(products []updating.ProductCSV) error {

	for _, product := range products {
		err := istore.updateCSVProductWithTx(product)
		if err != nil {
			continue
			//return fmt.Errorf("failed to bulk update products:%w product index:%d", err, index)
		}
	}
	return nil
}

func (istore *InventoryStore) updateCSVProduct(p updating.ProductCSV) error {

	query := `INSERT INTO  product (sku, country, name, quantity)
				VALUES($1, $2, $3, $4) 
			 	ON CONFLICT ON CONSTRAINT sku_country_unique
				DO 
	   			UPDATE SET quantity = EXCLUDED.quantity + product.quantity 
				WHERE EXCLUDED.quantity + product.quantity > 0;`

	_, err := istore.Exec(query, p.Sku, p.Country, p.Name, p.Quantity)

	if err != nil {
		log.Println("unable to insert product", err)
		return err
	}

	return nil
}

func (istore *InventoryStore) updateCSVProductWithTx(p updating.ProductCSV) error {

	tx := istore.MustBegin()
	var err error
	defer func() {
		if err != nil {
			log.Println("transaction rollback")
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	query := `SELECT id, sku, country, quantity, name FROM product WHERE sku= $1 AND country = $2 LIMIT 1`

	var product Product
	err = tx.Get(&product, query, p.Sku, p.Country)

	//item exist in the database
	if err == nil {
		newQuantity := product.Quantity + p.Quantity

		if newQuantity < 0 {
			err = fmt.Errorf("new quantity will be negative value not allowed")
			return err
		}

		updateQuery := `UPDATE product SET quantity = $1 WHERE id = $2`

		_, err = tx.Exec(updateQuery, newQuantity, product.Id)
		if err != nil {
			log.Println("unable to update product", err)
			return err
		}
		return nil
	}

	//item not exist
	if p.Quantity < 0 {
		err = fmt.Errorf("quantity is negative value not allowed for new items")
		return err
	}
	query = `INSERT INTO  product 
			(sku, country, name, quantity)
			VALUES($1, $2, $3, $4) 
			`
	_, err = tx.Exec(query, p.Sku, p.Country, p.Name, p.Quantity)
	if err != nil {
		log.Println("unable to insert product", err)
		return err
	}
	return nil
}
