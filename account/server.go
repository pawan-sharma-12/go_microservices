//go:generate protoc --go_out=paths=source_relative:./pb --go-grpc_out=paths=source_relative:./pb account.proto


package account

import (
	"context"
	"net"
	"google.golang.org/grpc"
	"fmt"
	"google.golang.org/grpc/reflection"

	 pb "github.com/pawan-sharma-12/go_microservices/account/pb"

)

type grpcServer struct {
	service Service
	pb.UnimplementedAccountServiceServer

}
func ListenAndServeGRPC(service Service, port int) error {
	// Implementation for starting gRPC server goes here
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcSrv := grpc.NewServer()
	pb.RegisterAccountServiceServer(grpcSrv,  &grpcServer{
		service: service,
		UnimplementedAccountServiceServer : pb.UnimplementedAccountServiceServer{},
		
	} )
	reflection.Register(grpcSrv)
	return grpcSrv.Serve(lis)
}
//POST /accounts
func (s *grpcServer) PostAccount(ctx context.Context, req *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	account, err := s.service.PostAccount(ctx, req.Name, req.Email)
	if err != nil {
		return nil, err
	}
	return &pb.PostAccountResponse{
		Account : &pb.Account{
			Id:    account.ID,
			Name:  account.Name,
			Email: account.Email,
		}, 
	}, nil
}

//GET /accounts
func (s *grpcServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest)(*pb.GetAccountResponse, error){
	account, err := s.service.GetAccountByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountResponse{
		Account : &pb.Account{
			Id:    account.ID,
			Name:  account.Name,
			Email: account.Email,
		},

	}, nil
}

//GET /accounts
func (s * grpcServer) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest)(*pb.GetAccountsResponse, error){
	accounts, err := s.service.GetAccounts(ctx, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}
	var pbAccounts []*pb.Account
	for _, account := range accounts {
		pbAccounts = append(pbAccounts, &pb.Account{
			Id:    account.ID,
			Name:  account.Name,
			Email: account.Email,
		})
	}
	return &pb.GetAccountsResponse{
		Accounts: pbAccounts,
	}, nil
}

