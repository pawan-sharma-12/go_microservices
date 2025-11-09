package order

import (
	"context"
	"log"
	"time"

	"github.com/pawan-sharma-12/go_microservices/order/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string)(*Client , error ){
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewOrderServiceClient(conn)
	return &Client{
		conn: conn,
		service: c,
	}, nil

}
func (c *Client) Close (){
	c.conn.Close()
}
func (c *Client) PostOrder(ctx context.Context, accountId string, products []OrderProduct)(*Order, error){
	protoProducts := []*pb.PostOrderRequest_OrderProduct{}
	for _, p := range products {
		protoProducts = append(protoProducts, &pb.PostOrderRequest_OrderProduct{
			ProductId: p.ID,
			Quantity:  p.Quantity,
		})

	}
	r, err := c.service.PostOrder(
		ctx,
		&pb.PostOrderRequest{
			AccountId: accountId,
			Products : protoProducts,
		},
	)
	if err != nil {
		log.Println("Order couldn't be posted")
		return nil , err
	}
	newOrder := r.Order
	newOrderCreatedAt := time.Time{}
	newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)
	return &Order {
		ID : newOrder.Id,
		CreatedAt: newOrderCreatedAt,
		TotalPrice: newOrder.TotalPrice,
		AccountID: newOrder.AccountId,
	}, nil
}
func (c *Client) GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error){
	r, err := c.service.GetOrderForAccount(ctx, &pb.GetOrderForAccountRequest{
		AccountId: accountId,
	})
	if err != nil {
		log.Println("Could not get orders with account id.. : ",accountId )
	}
	orders := []Order{}
	for _, orderProto := range r.Orders {
		newOrder := Order{
			ID:  orderProto.Id,
			TotalPrice: orderProto.TotalPrice,
			AccountID: orderProto.AccountId,
		}
		newOrder.CreatedAt = time.Time{}
		newOrder.CreatedAt.UnmarshalBinary(orderProto.CreatedAt)
		products := []OrderProduct{}
		for _, p := range orderProto.Products{
			products = append(products, OrderProduct{
				ID : p.Id,
				Quantity: p.Quantity,
				Name : p.Name,
				Description: p.Description,
				Price : p.Price,
			})
		}
		newOrder.Products  = products
		orders = append(orders, newOrder)
	}
	return orders, nil
}