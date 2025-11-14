package order

import (
	"context"
	"log"

	"github.com/pawan-sharma-12/go_microservices/order/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:    conn,
		service: pb.NewOrderServiceClient(conn),
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

// PostOrder calls the gRPC PostOrder
func (c *Client) PostOrder(ctx context.Context, accountID string, products []OrderProduct) (*Order, error) {
	reqProducts := make([]*pb.PostOrderRequest_OrderProduct, len(products))
	for i, p := range products {
		reqProducts[i] = &pb.PostOrderRequest_OrderProduct{
			ProductId: p.ID,
			Quantity:  p.Quantity,
		}
	}

	resp, err := c.service.PostOrder(ctx, &pb.PostOrderRequest{
		AccountId: accountID,
		Products:  reqProducts,
	})
	if err != nil {
		log.Println("Error posting order:", err)
		return nil, err
	}

	orderProducts := make([]OrderProduct, len(resp.Order.Products))
	for i, p := range resp.Order.Products {
		orderProducts[i] = OrderProduct{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Quantity:    p.Quantity,
			Price:       p.Price,
		}
	}

	return &Order{
		ID:         resp.Order.Id,
		AccountID:  resp.Order.AccountId,
		TotalPrice: resp.Order.TotalPrice,
		CreatedAt:  resp.Order.GetCreatedAt().AsTime(),
		Products:   orderProducts,
	}, nil
}

// GetOrdersForAccount calls gRPC GetOrderForAccount
func (c *Client) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	resp, err := c.service.GetOrderForAccount(ctx, &pb.GetOrderForAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		log.Println("Error fetching orders:", err)
		return nil, err
	}

	orders := make([]Order, len(resp.Orders))
	for i, o := range resp.Orders {
		orders[i] = Order{
			ID:         o.Id,
			AccountID:  o.AccountId,
			TotalPrice: o.TotalPrice,
			CreatedAt:  o.GetCreatedAt().AsTime(),
			Products:   convertOrderProtoToOrderProducts(o.Products),
		}
	}

	return orders, nil
}

// Helper: convert response products to internal OrderProduct
func convertOrderProtoToOrderProducts(protoProducts []*pb.Order_OrderProduct) []OrderProduct {
	products := make([]OrderProduct, len(protoProducts))
	for i, p := range protoProducts {
		products[i] = OrderProduct{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Quantity:    p.Quantity,
			Price:       p.Price,
		}
	}
	return products
}
