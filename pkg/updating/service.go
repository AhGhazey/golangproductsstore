package updating

import "fmt"

type Service interface {
	ConsumeProduct(p Product) (bool, error)
	UpdateBulkRecords(products []ProductCSV) error
}

type Repository interface {
	ConsumeProduct(p Product) (bool, error)
	UpdateBulkRecords(products []ProductCSV) error
}

type service struct {
	productRepository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) ConsumeProduct(p Product) (bool, error) {

	// any validation can be done here

	product, err := s.productRepository.ConsumeProduct(p)
	if err != nil {
		return false, fmt.Errorf("failed to get product:%w", err)
	}

	return product, nil
}

func (s *service) UpdateBulkRecords(products []ProductCSV) error {

	// any validation can be done here
	err := s.productRepository.UpdateBulkRecords(products)
	if err != nil {
		return fmt.Errorf("failed to get product:%w", err)
	}

	return nil
}
