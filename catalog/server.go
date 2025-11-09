package catalog

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/pawan-sharma-12/go_microservices/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)
type grpcServer struct {
	service Service
	pb.UnimplementedCatalogServiceServer

}
func ListenAndServeGRPC(service Service, port int) error {
	// Implementation for starting gRPC server goes here
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcSrv := grpc.NewServer()
	pb.RegisterCatalogServiceServer(grpcSrv,  &grpcServer{
		service: service,
		UnimplementedCatalogServiceServer : pb.UnimplementedCatalogServiceServer{},
		
	} )
	reflection.Register(grpcSrv)
	return grpcSrv.Serve(lis)
}

func (s * grpcServer) PostProduct( ctx context.Context, req *pb.PostProductRequest)(*pb.PostProductResponse, error){
	product, err := s.service.PostProduct(ctx, req.Name, req.Description, req.Price)
	if err != nil {
		log.Println("Error posting product:", err)
		return nil, err
	}
	return &pb.PostProductResponse{
		Product : &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		},
	}, nil
}

func (s * grpcServer) GetProduct( ctx context.Context, req *pb.GetProductRequest)(*pb.GetProductResponse, error){
	product, err := s.service.GetProduct(ctx, req.Id)
	if err != nil {
		log.Println("Error getting product:", err)
		return nil, err
	}
	return &pb.GetProductResponse{
		Product : &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		},
	}, nil
}
func (s * grpcServer) GetProducts( ctx context.Context, r *pb.GetProductsRequest)(*pb.GetProductsResponse, error){	
	var res []Product 
	var err error
	if r.Query != ""{
		res, err = s.service.SearchProducts(ctx, r.Query, r.Skip, r.Take)
	}else if len(r.Ids) > 0{
		res, err = s.service.GetProductsByIds(ctx, r.Ids)
	}else{
		res, err = s.service.GetProducts(ctx, r.Skip, r.Take)
	}
	if err != nil {
		log.Println("Error getting products:", err)
		return nil, err
	}
	var pbProducts []*pb.Product
	for _, product := range res {
		pbProducts = append(pbProducts, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return &pb.GetProductsResponse{
		Products: pbProducts,
	}, nil
}
