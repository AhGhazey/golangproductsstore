package adding

type Service interface {
	AddProduct(...Product) error
}

type Repository interface {
	// AddProduct saves a given product to the repository.
	AddProduct(Product) error
}

type service struct {
	productRepository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) AddProduct(p ...Product) error {

	// any validation can be done here

	for _, product := range p {
		err := s.productRepository.AddProduct(product) // error handling omitted for simplicity
		if err != nil {
			return err
		}
	}
	return nil
}
