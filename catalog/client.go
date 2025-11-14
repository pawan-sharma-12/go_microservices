package catalog

import (
	"context"
	"log"

	"github.com/pawan-sharma-12/go_microservices/catalog/pb"
	"google.golang.org/grpc"
)
type Client struct {
	conn *grpc.ClientConn
	service pb.CatalogServiceClient
}
func NewClient(url string) (*Client, error){
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewCatalogServiceClient(conn)
	return &Client{
		conn: conn,
		service: c,
	}, nil
}
func (c *Client) Close() {
	c.conn.Close()
}
func (c *Client) PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error) {
	r , err := c.service.PostProduct(ctx, &pb.PostProductRequest{
		Name: name,
		Description: description,
		Price: price,
	})
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, nil
}
func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	r , err := c.service.GetProduct(ctx, &pb.GetProductRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, nil
}
func (c *Client) GetProducts(ctx context.Context, skip uint64, take uint64, ids []string, query string) ([]Product, error) {
	r , err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Skip: skip,
		Take: take,
		Ids: ids,
		Query: query,
	})
	if err != nil {
		log.Println("client.go GetProducts error : ", err)
		return nil, err
	}
	var products []Product
	for _, p := range r.Products {
		products = append(products, Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, nil
}