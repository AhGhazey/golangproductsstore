package adding

type Product struct {
	Sku      string `json:"sku"`
	Country  string `json:"country"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
