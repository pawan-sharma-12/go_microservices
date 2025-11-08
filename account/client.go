package account

import (
	"context"
	"google.golang.org/grpc"
	 pb "github.com/pawan-sharma-12/go_microservices/account/pb"

)

type Client struct {
	conn *grpc.ClientConn
	service pb.AccountServiceClient
}
func NewClient(url string) (*Client, error){
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewAccountServiceClient(conn)
	return &Client{
		conn: conn,
		service: c,
	}, nil
}

func  (c *Client) PostAccount(ctx context.Context, name string, email string) (*Account, error) {
	r , err := c.service.PostAccount(ctx, &pb.PostAccountRequest{
		Name: name,
		Email: email,
	})
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:    r.Account.Id,
		Name:  r.Account.Name,
		Email: r.Account.Email,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r , err := c.service.GetAccount(ctx, &pb.GetAccountRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:    r.Account.Id,
		Name:  r.Account.Name,
		Email: r.Account.Email,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	r , err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{
		Skip: skip,
		Take: take,
	})
	if err != nil {
		return nil, err
	}
	accounts := make([]Account, 0, len(r.Accounts))
	for _, a := range r.Accounts {
		accounts = append(accounts, Account{
			ID:    a.Id,
			Name:  a.Name,
			Email: a.Email,
		})
	}
	return accounts, nil
}
func (c *Client) Close() error {
	return c.conn.Close()
}