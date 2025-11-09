package order

import (
	"context"
	"time"
	"github.com/segmentio/ksuid"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type Order struct {
	ID         string  `json:"id"`
	CreatedAt  time.Time   `json:"created_at"`
	AccountID  string  `json:"account_id"`
	TotalPrice float64 `json:"total_price"`
	Products   []OrderProduct `json:"products"`
}
type OrderProduct struct {
	ID 	  string `json:"id"`
	Name 	string `json:"name"`
	Description string `json:"description"`
	Price 	  float64 `json:"price"`
	Quantity 	uint64 `json:"quantity"`

}

type OrderService struct {
	repo Repository
}
func NewService(repo Repository) Service {
	return &OrderService{
		repo: repo,
	}
}
func (s *OrderService) PostOrder(ctx context.Context, accountID string, products []OrderProduct) (*Order, error) {
	var totalPrice float64
	for _, p := range products {
		totalPrice += p.Price * float64(p.Quantity)
	}
	order := Order{
		ID:         ksuid.New().String(),
		CreatedAt:  time.Now().UTC(),
		AccountID:  accountID,
		TotalPrice: totalPrice,
		Products:   products,
	}
	if err := s.repo.PutOrder(ctx, order); err != nil {
		return nil, err
	}
	return &order, nil
}


func (s *OrderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return s.repo.GetOrdersForAccount(ctx, accountID)
}