package postgres

type Product struct {
	Id       int64  `db:"id"`
	Name     string `db:"name"`
	Sku      string `db:"sku"`
	Country  string `db:"country"`
	Quantity int    `db:"quantity"`
}
