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
	"google.golang.org/protobuf/types/known/timestamppb"
)

// grpcServer implements pb.OrderServiceServer
type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
	pb.UnimplementedOrderServiceServer
}

// ListenGRPC starts the gRPC server
func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {
	log.Printf("🔗 Connecting to Account service at: %s", accountURL)
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		log.Printf("❌ Failed to connect to Account service: %v", err)
		return err
	}
	log.Printf("🔗 Connecting to Catalog service at: %s", catalogURL)
	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		log.Printf("❌ Failed to connect to Catalog service: %v", err)
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}

	grpcSrv := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcSrv, &grpcServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	})
	reflection.Register(grpcSrv)

	log.Printf("🚀 gRPC Order service running on port %d", port)
	return grpcSrv.Serve(lis)
}

// PostOrder handles creating a new order
func (s *grpcServer) PostOrder(ctx context.Context, req *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	// Convert request products to internal OrderProduct
	products := convertRequestProtoToOrderProducts(req.Products)

	// Validate account exists
	_, err := s.accountClient.GetAccount(ctx, req.AccountId)
	if err != nil {
		log.Println("❌ Error fetching account:", err)
		return nil, err
	}

	// Fetch product details from catalog
	productIDs := []string{}
	for _, p := range products {
		productIDs = append(productIDs, p.ID)
	}

	catalogProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("❌ Error fetching products:", err)
		return nil, err
	}

	// Merge quantities with catalog details
	for i := range products {
		for _, cp := range catalogProducts {
			if products[i].ID == cp.ID {
				products[i].Name = cp.Name
				products[i].Description = cp.Description
				products[i].Price = cp.Price
				break
			}
		}
	}

	// Create order
	order, err := s.service.PostOrder(ctx, req.AccountId, products)
	if err != nil {
		log.Println("❌ Error posting order:", err)
		return nil, errors.New("failed to post order")
	}

	// Convert to protobuf response
	return &pb.PostOrderResponse{
		Order: &pb.Order{
			Id:         order.ID,
			AccountId:  order.AccountID,
			TotalPrice: order.TotalPrice,
			CreatedAt:  timestamppb.New(order.CreatedAt),
			Products:   convertOrderProductsToProto(order.Products),
		},
	}, nil
}

// GetOrderForAccount fetches all orders for an account
func (s *grpcServer) GetOrderForAccount(ctx context.Context, req *pb.GetOrderForAccountRequest) (*pb.GetOrderForAccountResponse, error) {
	orders, err := s.service.GetOrdersForAccount(ctx, req.AccountId)
	if err != nil {
		log.Println("❌ Error fetching orders:", err)
		return nil, err
	}

	var protoOrders []*pb.Order
	for _, o := range orders {
		// Fetch product details from catalog for each order
		productIDs := []string{}
		for _, p := range o.Products {
			productIDs = append(productIDs, p.ID)
		}

		catalogProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
		if err != nil {
			log.Println("❌ Error fetching product details:", err)
			return nil, err
		}

		// Merge quantities with catalog details
		for i := range o.Products {
			for _, cp := range catalogProducts {
				if o.Products[i].ID == cp.ID {
					o.Products[i].Name = cp.Name
					o.Products[i].Description = cp.Description
					o.Products[i].Price = cp.Price
					break
				}
			}
		}

		protoOrders = append(protoOrders, &pb.Order{
			Id:         o.ID,
			AccountId:  o.AccountID,
			TotalPrice: o.TotalPrice,
			CreatedAt:  timestamppb.New(o.CreatedAt),
			Products:   convertOrderProductsToProto(o.Products),
		})
	}

	return &pb.GetOrderForAccountResponse{Orders: protoOrders}, nil
}

// Helper: convert request products to internal OrderProduct
func convertRequestProtoToOrderProducts(protoProducts []*pb.PostOrderRequest_OrderProduct) []OrderProduct {
	products := make([]OrderProduct, len(protoProducts))
	for i, p := range protoProducts {
		products[i] = OrderProduct{
			ID:       p.ProductId,
			Quantity: p.Quantity,
		}
	}
	return products
}

// Helper: convert internal OrderProduct to protobuf Order_OrderProduct
func convertOrderProductsToProto(products []OrderProduct) []*pb.Order_OrderProduct {
	protoProducts := make([]*pb.Order_OrderProduct, len(products))
	for i, p := range products {
		protoProducts[i] = &pb.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Quantity:    p.Quantity,
			Price:       p.Price,
		}
	}
	return protoProducts
}
