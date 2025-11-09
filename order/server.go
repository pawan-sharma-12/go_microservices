package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/pawan-sharma-12/go_microservices/account"
	"github.com/pawan-sharma-12/go_microservices/catalog"
	"github.com/pawan-sharma-12/go_microservices/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// grpcServer implements the gRPC service
type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
	pb.UnimplementedOrderServiceServer
}

// ListenGRPC starts the gRPC server
func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
		UnimplementedOrderServiceServer : pb.UnimplementedOrderServiceServer{},
	})
	reflection.Register(serv)
	log.Printf("🚀 gRPC Order service running on port %d", port)
	return serv.Serve(lis)
}

// PostOrder creates a new order
func (s *grpcServer) PostOrder(ctx context.Context, req *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, req.AccountId)
	if err != nil {
		log.Println("❌ Error fetching account:", err)
		return nil, err
	}

	productIDs := []string{}
	for _, p := range req.Products {
		productIDs = append(productIDs, p.ProductId)
	}

	orderedProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("❌ Error fetching products:", err)
		return nil, err
	}

	products := []OrderProduct{}
	for _, p := range orderedProducts {
		product := OrderProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    0,
		}

		for _, rp := range req.Products {
			if rp.ProductId == p.ID {
				product.Quantity = rp.Quantity
				break
			}
		}

		if product.Quantity > 0 {
			products = append(products, product)
		}
	}

	order, err := s.service.PostOrder(ctx, req.AccountId, products)
	if err != nil {
		log.Println("❌ Error posting order:", err)
		return nil, errors.New("failed to post order")
	}

	createdAt, _ := order.CreatedAt.MarshalBinary()
	orderProto := &pb.Order{
		Id:         order.ID,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  createdAt,
		Products:   []*pb.Order_OrderProduct{},
	}

	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}

	return &pb.PostOrderResponse{Order: orderProto}, nil
}

// GetOrdersForAccount returns all orders for an account
func (s *grpcServer) GetOrdersForAccount(ctx context.Context, req *pb.GetOrderForAccountRequest) (*pb.GetOrderForAccountResponse, error) {
	accountOrders, err := s.service.GetOrdersForAccount(ctx, req.AccountId)
	if err != nil {
		log.Println("❌ Could not find orders for account:", err)
		return nil, err
	}

	productIdMap := map[string]bool{}
	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIdMap[p.ID] = true
		}
	}

	productIDs := []string{}
	for id := range productIdMap {
		productIDs = append(productIDs, id)
	}

	products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("❌ Error getting product details:", err)
	}

	orders := []*pb.Order{}
	for _, o := range accountOrders {
		createdAt, _ := o.CreatedAt.MarshalBinary()
		op := &pb.Order{
			Id:         o.ID,
			AccountId:  o.AccountID,
			TotalPrice: o.TotalPrice,
			CreatedAt:  createdAt,
			Products:   []*pb.Order_OrderProduct{},
		}

		for _, product := range o.Products {
			for _, p := range products {
				if p.ID == product.ID {
					product.Name = p.Name
					product.Description = p.Description
					product.Price = p.Price
					break
				}
			}

			op.Products = append(op.Products, &pb.Order_OrderProduct{
				Id:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    product.Quantity,
			})
		}
		orders = append(orders, op)
	}

	return &pb.GetOrderForAccountResponse{Orders: orders}, nil
}
