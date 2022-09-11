package listing

import "fmt"

type Service interface {
	GetProductBySku(sku string) (*[]Product, error)
}

type Repository interface {
	GetProductBySku(sku string) (*[]Product, error)
}

type service struct {
	productRepository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetProductBySku(sku string) (*[]Product, error) {

	// any validation can be done here

	product, err := s.productRepository.GetProductBySku(sku)
	if err != nil {
		return nil, fmt.Errorf("failed to get product:%w", err)
	}

	return product, nil
}
