package listing

type Product struct {
	Name     string `json:"name"`
	Sku      string `json:"sku"`
	Country  string `json:"country"`
	Quantity int    `json:"quantity"`
}
