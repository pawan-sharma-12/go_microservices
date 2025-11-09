package catalog
import (
 context "context"
 "github.com/segmentio/ksuid"
)
type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       float64 `json:"price"`
}
type catalogService struct {
	repo Repository
}
type Service interface	 {
	Close()
	PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	GetProductsByIds(	ctx context.Context, ids []string) ([]Product, error)
	SearchProducts (ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}	
func NewService(repo Repository) Service {
	return &catalogService{
		repo: repo,
	}
}
func (s *catalogService) Close() {
	s.repo.Close()
}
func (s *catalogService) PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error) {
	product := Product{
		ID:          ksuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
	}
	if err := s.repo.PutProduct(ctx, product); err != nil {
		return nil, err
	}
	return &product, nil
}
func (s *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetProductByID(ctx, id)
}
func (s *catalogService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repo.ListProducts(ctx, skip, take)
}
func (s *catalogService) GetProductsByIds(ctx context.Context, ids []string) ([]Product, error) {
	return s.repo.ListProductsWithIDs(ctx, ids)
}
func (s *catalogService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repo.SearchProducts(ctx, query, skip, take)
}


