package updating

type Product struct {
	Sku      string `json:"sku"`
	Country  string `json:"country"`
	Quantity int    `json:"quantity"`
}

type ProductCSV struct {
	Sku      string `csv:"sku"`
	Name     string `csv:"name"`
	Country  string `csv:"country"`
	Quantity int    `csv:"stock_change"`
}
